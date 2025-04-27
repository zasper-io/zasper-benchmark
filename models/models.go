package models

import (
	"github.com/gorilla/websocket"
)

// SessionPayload and other related structs
type SessionPayload struct {
	Path   string `json:"path"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Kernel Kernel `json:"kernel"`
}

type Kernel struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	LastActivity   string `json:"last_activity"`
	ExecutionState string `json:"execution_state"`
	Connections    int    `json:"connections"`
}

type Notebook struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type Response struct {
	ID       string   `json:"id"`
	Path     string   `json:"path"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Kernel   Kernel   `json:"kernel"`
	Notebook Notebook `json:"notebook"`
}

type BenchmarkData struct {
	Timestamp                 string  `json:"timestamp"`
	CPUUsage                  float64 `json:"cpu_usage"`
	MemoryUsageMB             float64 `json:"memory_usage_mb"`
	MessagesSentCount         int64   `json:"messages_sent_count"`
	MessagesReceivedCount     int64   `json:"messages_received_count"`
	MessageSentThroughput     float64 `json:"message_sent_throughput"`
	MessageReceivedThroughput float64 `json:"message_received_throughput"`
}

type KernelWebSocketConnection struct {
	Conn                 *websocket.Conn
	Send                 chan []byte
	KernelId             string
	IOPubWindowMsgCount  int
	IOPubWindowByteCount int
	IOPubMsgsExceeded    int
	IOPubDataExceeded    int
}

// Store WebSocket connections
var kernelConnections map[string]KernelWebSocketConnection

// Channel to signal all WebSocket connections to send a kernel_info_request
var sendKernelInfoSignal chan struct{}

type MessageHeader struct {
	MsgID           string `json:"msg_id"`
	MsgType         string `json:"msg_type"`
	Username        string `json:"username"`
	Session         string `json:"session"`
	Date            string `json:"date"`
	ProtocolVersion string `json:"version"`
}

type MessageReceived struct {
	Header       MessageHeader `json:"header"`
	ParentHeader MessageHeader `json:"parent_header"`
	MsgId        string        `json:"msg_id"`
	MsgType      string        `json:"msg_type"`
	Content      interface{}   `json:"content"`
	Buffers      []byte        `json:"buffers"`
	Metadata     interface{}   `json:"metadata"`
	Tracker      int           `json:"tracker"`
	Error        error         `json:"error"`
	Channel      string        `json:"channel"`
}

type Content struct {
	Silent          bool                   `json:"silent"`
	StoreHistory    bool                   `json:"store_history"`
	UserExpressions map[string]interface{} `json:"user_expressions"`
	AllowStdin      bool                   `json:"allow_stdin"`
	StopOnError     bool                   `json:"stop_on_error"`
	Code            string                 `json:"code"`
}

type Header struct {
	Date     string `json:"date"`
	MsgID    string `json:"msg_id"`
	MsgType  string `json:"msg_type"`
	Session  string `json:"session"`
	Username string `json:"username"`
	Version  string `json:"version"`
}

type Metadata struct {
	DeletedCells []interface{} `json:"deleted_cells"`
	RecordTiming bool          `json:"record_timing"`
	CellID       string        `json:"cell_id"`
	Trusted      bool          `json:"trusted"`
}

type Message struct {
	Channel      string   `json:"channel"`
	Content      Content  `json:"content"`
	Header       Header   `json:"header"`
	Metadata     Metadata `json:"metadata"`
	ParentHeader Header   `json:"parent_header"`
}
