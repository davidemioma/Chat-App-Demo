package socket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Message struct {
	RoomID         string    `json:"roomId"`
	Content        string    `json:"content"`
	ClientID       string    `json:"clientId"`
    Username       string    `json:"username"`
}

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string            `json:"id"`
	RoomID   string            `json:"roomId"`
    Username string            `json:"username"`
}

type Room struct {
	ID      string              `json:"id"`
	Name    string              `json:"name"`
	Clients map[string] *Client `json:"clients"`
}

type Hub struct {
	Rooms      map[string] *Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func (h *Hub) Init() {
	h.Rooms = make(map[string]*Room)
	
	h.Register = make(chan *Client)

	h.Unregister = make(chan *Client)

	h.Broadcast = make(chan *Message)
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			// Check if room exists
			if _, exists := h.Rooms[cl.RoomID]; exists {
		        r := h.Rooms[cl.RoomID]

				// Check if client exists, if no add client to room.
				if _, exists := r.Clients[cl.ID]; !exists {
					r.Clients[cl.ID] = cl
				}
	        }

		case cl := <-h.Unregister:
			// Check if room exists
			if _, exists := h.Rooms[cl.RoomID]; exists {
				// Check if client is in the room
				if _, exists := h.Rooms[cl.RoomID].Clients[cl.ID]; exists {
					// Check if no client is in a room 
					if len(h.Rooms[cl.RoomID].Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  "user left the chat",
							RoomID:   cl.RoomID,
							Username: cl.Username,
						}
					}

					// Delete client
					delete(h.Rooms[cl.RoomID].Clients, cl.ID)

					// Close message channel of client
					close(cl.Message)
				}	
			}	

	    case msg := <-h.Broadcast:	
			// Check if room exists
			if _, exists := h.Rooms[msg.RoomID]; exists {
				// Log the broadcast message
				log.Printf("Broadcasting message to room %s: %s", msg.RoomID, msg.Content)

				// Send message to all clients
				for _, cl := range h.Rooms[msg.RoomID].Clients {
					select {
					case cl.Message <- msg:
						log.Printf("Message sent successfully")
					default:
						// If the client is not ready to receive, close the channel
						log.Printf("Closing message channel for client %s in room %s", cl.ID, msg.RoomID)

						close(cl.Message)

						delete(h.Rooms[msg.RoomID].Clients, cl.ID)
					}
				}
			} else {
				log.Printf("Room %s does not exist for broadcasting", msg.RoomID)
			}
		}
	}
}
