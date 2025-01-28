package socket

import (
	"context"
	"fmt"
	"net/http"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var ws = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")

		// For production
		return origin == "http://localhost:3000"
	},
}

func (h *Hub) GetRooms() []utils.GetRoomRes {
	var rooms []utils.GetRoomRes

	for _, room := range h.Rooms {
		rooms = append(rooms, utils.GetRoomRes{
			ID: room.ID,
			Name: room.Name,
		})
	}

	return rooms
}

func (h *Hub) GetClients(roomId string) ([]utils.GetClientRes, error) {
	var clients []utils.GetClientRes

	// Check if room exists
	if _, exists := h.Rooms[roomId]; !exists {
		return []utils.GetClientRes{}, fmt.Errorf("Room not found")
	}

	// Get clients
	for _, c := range h.Rooms[roomId].Clients {
		clients = append(clients, utils.GetClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	return clients, nil
}

func (h *Hub) CreateRoom(ctx context.Context, args utils.CreateRoomReq) error {
	if _, exists := h.Rooms[args.ID]; exists {
		return fmt.Errorf("room with ID %s already exists", args.ID)
	}

	h.Rooms[args.ID] = &Room{
		ID: args.ID,
		Name: args.Name,
		Clients: make(map[string]*Client),
	}

	return nil
}

func (h *Hub) JoinRoom(c *gin.Context, arg utils.JoinRoomReq) error {
	conn, wsErr := ws.Upgrade(c.Writer, c.Request, nil)

	if wsErr != nil {
		return fmt.Errorf("web socket upgrade error: %v", wsErr)
    }

	cl := &Client{
		Conn: conn,
		ID: arg.UserID,
		RoomID: arg.RoomID,
		Username: arg.Username,
		Message: make(chan *Message, 10),
	}

	msg := &Message{
		Content:  "A new user has joined the room",
		ClientID: arg.UserID,
		RoomID:   arg.RoomID,
		Username: arg.Username,
	}

	// Register a client through register channel
	h.Register <- cl

	// Broadcast the message
	h.Broadcast <- msg

	// Write Message
	go cl.writeMessage()

	// Read Message
	cl.readMessage(h)

	return nil
}