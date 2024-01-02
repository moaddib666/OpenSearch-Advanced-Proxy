package websockets

import (
	"OpenSearchAdvancedProxy/internal/core/ports"
	"context"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type ServerClient struct {
	Conn           *websocket.Conn
	Server         *WebSocketServer
	EventProcessor ports.WebsocketServerEventProcessor
	writeMu        *sync.Mutex
}

// NewServerClient creates a new client
func NewServerClient(conn *websocket.Conn, server *WebSocketServer, eventProcessor ports.WebsocketServerEventProcessor) *ServerClient {
	return &ServerClient{
		Conn:           conn,
		Server:         server,
		EventProcessor: eventProcessor,
		writeMu:        &sync.Mutex{},
	}
}

func (c *ServerClient) SendMessage(message []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return c.Conn.WriteMessage(websocket.TextMessage, message)
}

// Listen for the close event in a separate goroutine
func (c *ServerClient) Listen() {
	defer func() {
		log.Debugf("Closing connection: %s", c.Conn.RemoteAddr().String())
		c.Server.Unregister <- c
		_ = c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		// log.Debugf("Received message: %s", string(message))
		if err != nil {
			break
		}

		// Pass the incoming message to the event processor
		err = c.EventProcessor.OnIncomingMessage(c, message)
		if err != nil {
			log.Errorf("Error processing incoming message: %s", err.Error())
			break
		}
	}
}

type WebSocketServer struct {
	Clients        map[*ServerClient]bool
	Register       chan *ServerClient
	Unregister     chan *ServerClient
	Mutex          sync.Mutex
	Upgrader       websocket.Upgrader
	EventProcessor ports.WebsocketServerEventProcessor
	httpServer     *http.ServeMux
	bindAddress    string
}

func NewWebSocketServer(bindAddress string, processor ports.WebsocketServerEventProcessor) *WebSocketServer {
	return &WebSocketServer{
		Clients:        make(map[*ServerClient]bool),
		Register:       make(chan *ServerClient),
		Unregister:     make(chan *ServerClient),
		Upgrader:       websocket.Upgrader{},
		EventProcessor: processor,
		bindAddress:    bindAddress,
		httpServer:     http.NewServeMux(),
	}
}

func (server *WebSocketServer) start(ctx context.Context) {
	server.httpServer.HandleFunc("/", server.HandleWebSocket)

	// Start the HTTP server in a goroutine
	go func() {
		log.Println("WebSocket server starting on", server.bindAddress)
		if err := http.ListenAndServe(server.bindAddress, server.httpServer); err != nil {
			log.Fatalf("Error starting WebSocket server: %s", err)
		}
	}()

	// Listen for context cancellation
	<-ctx.Done()
	log.Println("Shutting down WebSocket server")
}

func (server *WebSocketServer) Run(ctx context.Context) {
	go server.start(ctx)
	for {
		select {
		case client := <-server.Register:
			server.Mutex.Lock()
			server.Clients[client] = true
			server.Mutex.Unlock()
			go client.Listen() // Start listening for close event
		case client := <-server.Unregister:
			server.Mutex.Lock()
			if _, ok := server.Clients[client]; ok {
				delete(server.Clients, client)
			}
			server.Mutex.Unlock()
		case <-ctx.Done():
			return // Stop the server when context is canceled
		}
	}
}

func (server *WebSocketServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Handling WebSocket connection: %s", r.RemoteAddr)
	conn, err := server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Error upgrading to WebSocket:", err)
		return
	}
	client := NewServerClient(conn, server, server.EventProcessor)
	server.Register <- client
}

func (server *WebSocketServer) Broadcast(message []byte) int {
	for client := range server.Clients {
		err := client.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Warnf("Error broadcasting message: %s", err.Error())
		}
	}
	return len(server.Clients)
}
