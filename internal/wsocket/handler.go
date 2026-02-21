package wsocket

import (
	"net/http"

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
	Hub        *Hub
	JWTService *service.JWTService
}

func NewWsHandler(hub *Hub, jwt *service.JWTService) *WebSocketHandler {
	return &WebSocketHandler{
		Hub:        hub,
		JWTService: jwt,
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

	go client.WritePump()
	go client.ReadPump()
}
