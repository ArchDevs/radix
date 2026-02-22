package message

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *MessageService
}

func NewHandler(svc *MessageService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetHistory(c *gin.Context) {
	user := c.GetString("address")
	other := c.Query("with")

	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	messages, err := h.svc.GetMessageHistory(c.Request.Context(), user, other, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get message history"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *Handler) GetUndelivered(c *gin.Context) {
	user := c.GetString("address")

	messages, err := h.svc.GetUndeliveredMessages(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get undelivered messages"})
		return
	}

	// Mark delivered
	for _, m := range messages {
		if err := h.svc.MarkAsDelivered(c.Request.Context(), m.ID); err != nil {
			// Log error but continue
			continue
		}
	}

	c.JSON(http.StatusOK, messages)
}
