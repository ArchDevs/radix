package wsocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func (c *Client) ReadPump() {
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

		c.handleMessage(message)
	}
}

func (c *Client) handleMessage(data []byte) {
	var msg IncomingMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return // invalid JSON
	}

	// Send to recipient
	if recipient, ok := c.Hub.GetClient(msg.To); ok {
		out := OutgoingMessage{
			From:    c.Address,
			Content: msg.Content,
		}
		payload, _ := json.Marshal(out)
		recipient.Send <- payload
	} else {
		// TODO Save to database
		log.Printf("user %s offline, message saved", msg.To)
	}
}
