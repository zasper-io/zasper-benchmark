package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/zasper-io/bechmark/models"
	"github.com/zasper-io/bechmark/monitoring"
)

var (
	wg sync.WaitGroup
)

// Store WebSocket connections
var kernelConnections map[string]models.KernelWebSocketConnection

// Function to send a POST request asynchronously
func startKernelSession(
	url string, wsURL string,
	ctx context.Context, triggerChan chan string,
	token string, xsrfToken string, numRequest int,
) int {

	payload := models.SessionPayload{
		Type: "notebook",
		Kernel: models.Kernel{
			Name: "python3",
		},
	}
	payload.Name = fmt.Sprintf("Untitled-%d.ipynb", numRequest)
	payload.Path = payload.Name

	data, err := json.Marshal(payload)
	log.Debug().Msgf("payload name %s", payload.Name)
	if err != nil {
		log.Debug().Msgf("Error marshalling payload: %v", err)
		return -1
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Debug().Msgf("Error creating request: %v", err)
		return -1
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("_xsrf", xsrfToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Debug().Msgf("Request failed: %v", err)
		return -1
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Debug().Msgf("Error reading response body: %v", err)
		return -1
	}

	log.Debug().Msgf("Response Status: %s", resp.Status)
	log.Debug().Msgf("Response Body: %s", string(responseBody))

	var response models.Response
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		log.Debug().Msgf("Error unmarshaling response: %v", err)
		return -1
	}

	kernelID := response.Kernel.ID
	sessionID := response.ID
	wsURL = fmt.Sprintf(wsURL, kernelID, sessionID)

	wg.Add(1)
	log.Debug().Msgf("created session  for %s , %s", kernelID, sessionID)
	go addKernelConnectionWorker(kernelID, sessionID, wsURL, ctx, triggerChan, token, xsrfToken)

	return resp.StatusCode
}

// Fire once Function to manage WebSocket connection for each kernel
func addKernelConnectionWorker(kernelID string, sessionID string, wsURL string,
	ctx context.Context, triggerChan chan string,
	token string, xsrfToken string,
) {
	defer wg.Done()

	headers := map[string][]string{
		"Authorization": {"Bearer " + token},
		"_xsrf":         {xsrfToken},
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
	if err != nil {
		log.Fatal().Msgf("Failed to connect to WebSocket: %v", err)
	}

	kernelConnections[kernelID] = models.KernelWebSocketConnection{
		Conn:     conn,
		KernelId: kernelID,
	}

	log.Debug().Msgf("WebSocket connection established for kernel: %s", kernelID)
	go listenForMessages(conn)

	if err := godotenv.Load(); err != nil {
		log.Debug().Msg("Error loading .env file")
	}

	delayStr := os.Getenv("DELAY")
	delay, err := strconv.Atoi(delayStr)
	if err != nil {
		log.Fatal().Msgf("Invalid DELAY value: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("Context cancelled for kernel: %s", kernelID)
			return
		case <-triggerChan:
			log.Debug().Msgf("started sending request %s %s", kernelID, sessionID)

			start := time.Now()
			ticker := time.NewTicker(time.Duration(delay) * time.Millisecond) // send request every 100ms
			defer ticker.Stop()
			for t := range ticker.C {
				if time.Since(start) > 60*time.Second {
					log.Debug().Msg("Stop sending requests")
					return
				}
				log.Debug().Msgf("triggering request at %v", t)
				sendKernelExecuteRequest(conn, kernelID, sessionID)
			}

		}
	}
}

// Function to send kernel_info_request to the WebSocket connection
func sendKernelExecuteRequest(conn *websocket.Conn, kernelID string, sessionID string) {
	atomic.AddInt64(&monitoring.MessagesSentCount, 1)
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
		log.Debug().Msgf("Failed to send kernel_info_request: %v", err)
	}
	log.Debug().Msgf("Kernel Execute Request sent for kernel: %s, %s", kernelID, sessionID)
}

// Function to listen for incoming messages from the WebSocket
func listenForMessages(conn *websocket.Conn) {
	for {
		messageType, msg, err := conn.ReadMessage()
		atomic.AddInt64(&monitoring.MessagesReceivedCount, 1)
		if err != nil {
			log.Debug().Msgf("Error reading message: %v", err)
			return
		}

		switch messageType {
		case websocket.TextMessage:
			var jsonMsg models.MessageReceived
			if err := json.Unmarshal(msg, &jsonMsg); err != nil {
				log.Debug().Msgf("Received non-JSON message: %s", string(msg))
			} else {
				log.Debug().Msgf("Received JSON message: %+v", jsonMsg)
			}
		case websocket.BinaryMessage:
			log.Debug().Msgf("Received binary message of length: %d", len(msg))
		default:
			log.Debug().Msgf("Received message of unknown type: %d", messageType)
		}
	}
}

// Function to collect the performace data
func measurePerformance(
	url string, wsURL string, ctx context.Context, triggerChan chan string, numRequests int,
) {
	fmt.Println("Creating kernel sessions ⏳")
	if err := godotenv.Load(); err != nil {
		log.Info().Msg("Error loading .env file")
	}

	token := os.Getenv("TOKEN")
	xsrfToken := os.Getenv("XSRF_TOKEN")
	for i := 0; i < numRequests; i++ {
		statusCode := startKernelSession(url, wsURL, ctx, triggerChan, token, xsrfToken, i)
		if statusCode != 201 {
			log.Debug().Msgf("Request failed with status: %d", statusCode)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// Benchmark against Zasper
func measureZasperPerformance(numKernels int, ctx context.Context, triggerChan chan string) {
	zasperURL := "http://localhost:8048/api/sessions"
	wsURL := "ws://localhost:8048/api/kernels/%s/channels?session_id=%s"

	go measurePerformance(zasperURL, wsURL, ctx, triggerChan, numKernels)
}

// Benchmark against Jupyter
func measureJupyterPerformance(numKernels int, ctx context.Context, triggerChan chan string) {

	jupyterURL := "http://localhost:8888/api/sessions"
	wsURL := "ws://localhost:8888/api/kernels/%s/channels?session_id=%s"

	go measurePerformance(jupyterURL, wsURL, ctx, triggerChan, numKernels)
}

// Entry point
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	monitoring.InitializeBenchmarkResults()
	kernelConnections = make(map[string]models.KernelWebSocketConnection)
	triggerChan := make(chan string)

	if err := godotenv.Load(); err != nil {
		log.Debug().Msg("Error loading .env file")
	}

	target := os.Getenv("TARGET")
	numKernelsStr := os.Getenv("NUM_KERNELS")
	numKernels, err := strconv.Atoi(numKernelsStr)
	if err != nil {
		log.Fatal().Msgf("Invalid NUM_KERNELS value: %v", err)
	}

	pidStr := os.Getenv("PID")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		log.Fatal().Msgf("Invalid PID value: %v", err)
	}

	delayStr := os.Getenv("DELAY")
	delay, err := strconv.Atoi(delayStr)
	if err != nil {
		log.Fatal().Msgf("Invalid DELAY value: %v", err)
	}

	resultFile := fmt.Sprintf("data/%dms/benchmark_results_%s_%dkernels.json", delay, target, numKernels)
	fmt.Println("====================================================================")
	fmt.Println("*******            Measuring performance                     *******")
	fmt.Println("====================================================================")
	fmt.Println("Target:", target)
	fmt.Println("PID:", pid)
	fmt.Println("Number of kernels:", numKernels)
	fmt.Println("Output file:", resultFile)
	fmt.Println("====================================================================")
	go monitoring.MonitorProcessByPID(int32(pid))
	go monitoring.WriteBenchmarkResultsPeriodically(resultFile, 2*time.Second)

	ctx, cancel := context.WithCancel(context.Background())

	// Choose which backend to benchmark
	switch target {
	case "jupyter":
		measureJupyterPerformance(numKernels, ctx, triggerChan)
	case "zasper":
		measureZasperPerformance(numKernels, ctx, triggerChan)
	default:
		log.Fatal().Msgf("Unknown target: %s. Use 'jupyter' or 'zasper'", target)
	}
	time.Sleep(10 * time.Second)
	fmt.Println("Sessions created:  ✅ ")
	fmt.Println("Start sending requests: ⏳")
	for i := 0; i < numKernels; i++ {
		triggerChan <- "start"
	}

	time.Sleep(4 * time.Second)
	wg.Wait()
	fmt.Println("Kernel messages sent:  ✅ ")
	fmt.Println("====================================================================")
	fmt.Println("*******                   Summary                            *******")
	fmt.Println("====================================================================")
	fmt.Println("Messages sent:", monitoring.TotalMessagesSentCount)
	fmt.Println("Messages received:", monitoring.TotalMessagesReceievedCount)
	fmt.Println("====================================================================")

	cancel()

}
