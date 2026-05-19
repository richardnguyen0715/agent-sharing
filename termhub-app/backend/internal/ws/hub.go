// Package ws manages WebSocket client connections and message broadcasting.
package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 8192
)

// WSMessage is the JSON message format for WebSocket communication.
type WSMessage struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
	Cols uint16 `json:"cols,omitempty"`
	Rows uint16 `json:"rows,omitempty"`
}

// Client represents a connected WebSocket client.
type Client struct {
	ID         string
	TerminalID uuid.UUID
	UserID     string
	Conn       *websocket.Conn
	Send       chan []byte
	hub        *Hub
	done       chan struct{}
}

// Hub maintains the set of active WebSocket clients and broadcasts messages.
type Hub struct {
	mu         sync.RWMutex
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	logger     zerolog.Logger
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     log.With().Str("component", "ws-hub").Logger(),
	}
}

// Run starts the hub's main event loop. It should be called in a goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
			h.logger.Info().
				Str("client_id", client.ID).
				Str("terminal_id", client.TerminalID.String()).
				Str("user_id", client.UserID).
				Msg("client registered")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()
			h.logger.Info().
				Str("client_id", client.ID).
				Str("terminal_id", client.TerminalID.String()).
				Msg("client unregistered")
		}
	}
}

// Register adds a client to the hub.
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub.
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// GetClientsByTerminal returns all clients connected to a specific terminal.
func (h *Hub) GetClientsByTerminal(terminalID uuid.UUID) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var clients []*Client
	for _, client := range h.clients {
		if client.TerminalID == terminalID {
			clients = append(clients, client)
		}
	}
	return clients
}

// SendToTerminal broadcasts data to all clients connected to a specific terminal.
func (h *Hub) SendToTerminal(terminalID uuid.UUID, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		if client.TerminalID == terminalID {
			select {
			case client.Send <- data:
			default:
				h.logger.Warn().Str("client_id", client.ID).Msg("client send buffer full, dropping message")
			}
		}
	}
}

// ClientCount returns the number of connected clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// NewClient creates a new Client instance.
func NewClient(terminalID uuid.UUID, userID string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:         uuid.New().String(),
		TerminalID: terminalID,
		UserID:     userID,
		Conn:       conn,
		Send:       make(chan []byte, 256),
		hub:        hub,
		done:       make(chan struct{}),
	}
}

// WritePump pumps messages from the Send channel to the websocket connection.
// A goroutine running WritePump is started for each connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel.
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.done:
			return
		}
	}
}

// ReadPump pumps messages from the websocket connection.
// It parses JSON messages and writes stdin data to the PTY or handles resize.
func (c *Client) ReadPump(ptyWriter func([]byte) error, resizeFn func(cols, rows uint16) error) {
	defer func() {
		c.hub.Unregister(c)
		c.Conn.Close()
		close(c.done)
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Error().Err(err).Str("client_id", c.ID).Msg("unexpected websocket close")
			}
			return
		}

		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			// If not valid JSON, treat as raw stdin input.
			if writeErr := ptyWriter(message); writeErr != nil {
				log.Error().Err(writeErr).Str("client_id", c.ID).Msg("failed to write to PTY")
				return
			}
			continue
		}

		switch msg.Type {
		case "stdin":
			if writeErr := ptyWriter([]byte(msg.Data)); writeErr != nil {
				log.Error().Err(writeErr).Str("client_id", c.ID).Msg("failed to write to PTY")
				return
			}
		case "resize":
			if msg.Cols > 0 && msg.Rows > 0 {
				if resizeErr := resizeFn(msg.Cols, msg.Rows); resizeErr != nil {
					log.Error().Err(resizeErr).Str("client_id", c.ID).Msg("failed to resize PTY")
				}
			}
		default:
			log.Warn().Str("type", msg.Type).Str("client_id", c.ID).Msg("unknown websocket message type")
		}
	}
}
