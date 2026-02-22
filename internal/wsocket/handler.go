package wsocket

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ArchDevs/radix/internal/message"
	"github.com/ArchDevs/radix/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins for development
	},
}

type WebSocketHandler struct {
	Hub            *Hub
	JWTService     *service.JWTService
	MessageService *message.MessageService
}

func NewWsHandler(hub *Hub, jwt *service.JWTService, msgSvc *message.MessageService) *WebSocketHandler {
	return &WebSocketHandler{
		Hub:            hub,
		JWTService:     jwt,
		MessageService: msgSvc,
	}
}

func (h *WebSocketHandler) Handle(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	address, err := h.JWTService.Parse(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		Hub:     h.Hub,
		Conn:    conn,
		Address: address,
		Send:    make(chan []byte, 256),
	}

	h.Hub.Register(client)

	// Send undelivered messages
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		msgs, err := h.MessageService.GetUndeliveredMessages(ctx, address)
		if err != nil {
			return
		}

		for _, m := range msgs {
			out := OutgoingMessage{
				ID:        m.ID,
				From:      m.Sender,
				Content:   m.Content,
				Timestamp: m.CreatedAt.Unix(),
			}
			payload, _ := json.Marshal(out)
			client.Send <- payload
			h.MessageService.MarkAsDelivered(ctx, m.ID)
		}
	}()

	go client.WritePump()
	go client.ReadPump(h.MessageService)
}
