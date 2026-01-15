package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	pollIntervalSeconds = 15
	httpPort            = "8081"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	broadcast = make(chan Event)
)

type Event struct {
	Message string `json:"message"`
}

func broadcastEvent(event Event) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	message, _ := json.Marshal(event)
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("WebSocket error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// handleBroadcast listens for events and broadcasts them
func handleBroadcast() {
	for event := range broadcast {
		broadcastEvent(event)
	}
}

// handleWebSocket handles WebSocket connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	log.Printf("New WebSocket client connected")

	// Keep connection alive and handle disconnection
	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
		log.Printf("WebSocket client disconnected")
	}()

	// Read messages (just to detect disconnection)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// startHTTPServer starts the HTTP server
func startHTTPServer() error {
	r := mux.NewRouter()

	// API endpoints
	r.HandleFunc("/api/ws", handleWebSocket)

	// Serve static files from the frontend build directory
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend")))

	log.Printf("Starting HTTP server on port %s", httpPort)
	return http.ListenAndServe(":"+httpPort, r)
}

func StartServer() chan Event {

	// Start broadcast handler
	go handleBroadcast()

	// Start HTTP server in a goroutine
	go func() {
		if err := startHTTPServer(); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Return the channel to broadcast messages
	return broadcast
}
