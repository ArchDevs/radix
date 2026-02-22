package wsocket

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ArchDevs/radix/internal/message"
	"github.com/gorilla/websocket"
)

func (c *Client) ReadPump(msgSvc *message.MessageService) {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512 * 1024) // 512 kb max message size

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}

		c.handleMessage(message, msgSvc)
	}
}

func (c *Client) handleMessage(data []byte, msgSvc *message.MessageService) {
	var msg IncomingMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("failed to unmarshal message: %v, data: %s", err, string(data))
		return // invalid JSON
	}

	log.Printf("received message from %s to %s (content length: %d)", c.Address, msg.To, len(msg.Content))

	dbMsg, err := msgSvc.Create(context.Background(), c.Address, msg.To, msg.Content)
	if err != nil {
		log.Printf("failed to save message: %v", err)
		return
	}

	if recipient, ok := c.Hub.GetClient(msg.To); ok {
		out := OutgoingMessage{
			ID:        dbMsg.ID,
			From:      c.Address,
			Content:   msg.Content,
			Timestamp: msg.Timestamp,
		}
		payload, _ := json.Marshal(out)

		select {
		case recipient.Send <- payload:
			// Delivered successfully, mark in DB
			if err := msgSvc.MarkAsDelivered(context.Background(), dbMsg.ID); err != nil {
				log.Printf("failed to mark message as delivered: %v", err)
			}
		default:
			// Client buffer full, message stays undelivered
			log.Printf("recipient %s buffer full, message %s stays undelivered", msg.To, dbMsg.ID)
		}
	} else {
		// Recipient offline, message stays undelivered
		log.Printf("user %s offline, message %s saved for later delivery", msg.To, dbMsg.ID)
	}
}
