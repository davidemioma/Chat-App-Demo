package socket

import (
	"log"

	"github.com/gorilla/websocket"
)

func (c *Client) writeMessage() {
	defer func ()  {
		c.Conn.Close()
	}()

	for {
		// Get message from channel
		message, ok := <-c.Message

		if !ok {
			return
		}

		err := c.Conn.WriteJSON(message)

		if err != nil{
			log.Printf("Error writing message for client %v: %v", c.Username, err)

			c.Conn.Close()
		}
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func ()  {
		hub.Unregister <- c

		c.Conn.Close()
	}()

	for {
		// Read message from channel
		_, m, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Websocket error: %v", err)
			}

			break
		}

		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			ClientID: c.ID,
			Username: c.Username,
		}

		hub.Broadcast <- msg
	}
}