package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"github.com/zasper-io/bechmark/models"
	monitor "github.com/zasper-io/bechmark/montoring"
)

var (
	wg sync.WaitGroup
)

// Store WebSocket connections
var kernelConnections map[string]models.KernelWebSocketConnection

// Function to send a POST request asynchronously
func startKernelSession(
	url string, wsURL string, payload models.SessionPayload, messageChannel chan string,
	token string, xsrfToken string, numRequest int,
) int {
	payload.Name = fmt.Sprintf("Untitled-%d.ipynb", numRequest)
	payload.Path = payload.Name

	data, err := json.Marshal(payload)
	log.Println("payload name", payload.Name)
	if err != nil {
		log.Printf("Error marshalling payload: %v", err)
		return -1
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return -1
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("_xsrf", xsrfToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return -1
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return -1
	}

	log.Printf("Response Status: %s", resp.Status)
	log.Printf("Response Body: %s", string(responseBody))

	var response models.Response
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return -1
	}

	kernelID := response.Kernel.ID
	sessionID := response.ID
	wsURL = fmt.Sprintf(wsURL, kernelID, sessionID)

	wg.Add(1)
	go addKernelConnectionWorker(kernelID, sessionID, wsURL, messageChannel, token, xsrfToken)

	return resp.StatusCode
}

// Function to manage WebSocket connection for each kernel
func addKernelConnectionWorker(kernelID string, sessionID string, wsURL string, messageChannel chan string, token string, xsrfToken string) {
	defer wg.Done()

	headers := map[string][]string{
		"Authorization": {"Bearer " + token},
		"_xsrf":         {xsrfToken},
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}

	kernelConnections[kernelID] = models.KernelWebSocketConnection{
		Conn:     conn,
		KernelId: kernelID,
	}

	log.Printf("WebSocket connection established for kernel: %s", kernelID)
	go listenForMessages(conn)

	for {
		select {
		case <-messageChannel:
			sendKernelInfoRequest(conn, kernelID, sessionID)
			log.Println("Sending message to kernel", kernelID)
		}
	}
}

// Fire once Function to manage WebSocket connection for each kernel
func addKernelConnectionWorkerSlow(kernelID string, sessionID string, wsURL string, messageChannel chan string, token string, xsrfToken string) {
	defer wg.Done()

	headers := map[string][]string{
		"Authorization": {"Bearer " + token},
		"_xsrf":         {xsrfToken},
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}

	kernelConnections[kernelID] = models.KernelWebSocketConnection{
		Conn:     conn,
		KernelId: kernelID,
	}

	log.Printf("WebSocket connection established for kernel: %s", kernelID)
	go listenForMessages(conn)

	for {
		msg, ok := <-messageChannel
		if !ok {
			log.Printf("Message channel closed. Exiting worker for kernel: %s", kernelID)
			return
		}
		log.Printf("Sending message to kernel %s: %s", kernelID, msg)
		sendKernelInfoRequest(conn, kernelID, sessionID)
		time.Sleep(50 * time.Millisecond) // comment out in zasper
	}
}

// Function to send kernel_info_request to the WebSocket connection
func sendKernelInfoRequest(conn *websocket.Conn, kernelID string, sessionID string) {
	msg := models.Message{
		Channel: "shell",
		Content: models.Content{
			Silent:          false,
			StoreHistory:    true,
			UserExpressions: map[string]interface{}{},
			AllowStdin:      true,
			StopOnError:     true,
			Code:            "2+2",
		},
		Header: models.Header{
			Date:     time.Now().Format(time.RFC3339),
			MsgID:    uuid.New().String(),
			MsgType:  "execute_request",
			Session:  sessionID,
			Username: "prasunanand",
			Version:  "5.2",
		},
		Metadata: models.Metadata{
			DeletedCells: []interface{}{},
			RecordTiming: false,
			CellID:       uuid.New().String(),
			Trusted:      true,
		},
		ParentHeader: models.Header{},
	}

	err := conn.WriteJSON(msg)
	if err != nil {
		log.Printf("Failed to send kernel_info_request: %v", err)
	}
	log.Printf("Kernel Info Request sent for kernel: %s, %s", kernelID, sessionID)
}

// Function to listen for incoming messages from the WebSocket
func listenForMessages(conn *websocket.Conn) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		switch messageType {
		case websocket.TextMessage:
			var jsonMsg models.MessageReceived
			if err := json.Unmarshal(msg, &jsonMsg); err != nil {
				log.Printf("Received non-JSON message: %s", string(msg))
			} else {
				log.Printf("Received JSON message: %+v", jsonMsg)
			}
		case websocket.BinaryMessage:
			log.Printf("Received binary message of length: %d", len(msg))
		default:
			log.Printf("Received message of unknown type: %d", messageType)
		}
	}
}

// Function to run the benchmark and monitor resources
func benchmark(
	url string, wsURL string, payload models.SessionPayload, messageChannel chan string,
	token string, xsrfToken string, numRequests int, resultFile string,
) {
	for i := 0; i < numRequests; i++ {
		statusCode := startKernelSession(url, wsURL, payload, messageChannel, token, xsrfToken, i)
		if statusCode != 201 {
			log.Printf("Request failed with status: %d", statusCode)
		}
		time.Sleep(50 * time.Millisecond)
	}
	time.Sleep(5 * time.Second) // comment out in zasper (JupyterLab nudges the kernel)
}

// Benchmark against Zasper
func benchmarkZasper(payload models.SessionPayload, messageChannel chan string) {
	zasperURL := "http://localhost:8048/api/sessions"
	wsURL := "ws://localhost:8048/api/kernels/%s/channels?session_id=%s"

	go benchmark(zasperURL, wsURL, payload, messageChannel, "", "", 2, "zasper_benchmark_results.json")
}

// Benchmark against Jupyter
func benchmarkJupyter(payload models.SessionPayload, messageChannel chan string) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")
	xsrfToken := os.Getenv("XSRF_TOKEN")

	jupyterURL := "http://localhost:8888/api/sessions"
	wsURL := "ws://localhost:8888/api/kernels/%s/channels?session_id=%s"

	go benchmark(jupyterURL, wsURL, payload, messageChannel, token, xsrfToken, 2, "zasper_benchmark_results.json")
}

// Entry point
func main() {
	monitor.InitializeBenchmarkResults()
	kernelConnections = make(map[string]models.KernelWebSocketConnection)
	messageChannel := make(chan string)

	payload := models.SessionPayload{
		Type: "notebook",
		Kernel: models.Kernel{
			Name: "python3",
		},
	}

	// Parse CLI flag
	target := flag.String("target", "zasper", "Specify which backend to benchmark: 'jupyter' or 'zasper'")
	flag.Parse()

	go monitor.MonitorProcessByPID(42194)
	go monitor.WriteBenchmarkResultsPeriodically("benchmark_results_new.json", 15*time.Second)

	// Choose which backend to benchmark
	switch *target {
	case "jupyter":
		benchmarkJupyter(payload, messageChannel)
	case "zasper":
		benchmarkZasper(payload, messageChannel)
	default:
		log.Fatalf("Unknown target: %s. Use 'jupyter' or 'zasper'", *target)
	}
	go func() {
		time.Sleep(5 * time.Second)
		messageChannel <- "Hello from the channel!"
		close(messageChannel)
	}()

	time.Sleep(10 * time.Second)
	wg.Wait()
}
