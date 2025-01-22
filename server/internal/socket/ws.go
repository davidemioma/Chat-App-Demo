package socket

import (
	"context"
	"fmt"
	"net/http"
	"server/utils"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	RoomID         string    `json:"roomId"`
	Content        string    `json:"content"`
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
	Mutex      sync.RWMutex
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

		case cl := <-h.Unregister:

	    case m := <-h.Broadcast:		
		}
	}
}

var ws = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// origin := r.Header.Get("Origin")

		// For production
		// return origin == "Your frontend url"

		return true
	},
}

func (h *Hub) CreateRoom(ctx context.Context, args utils.CreateRoomReq) error {
	h.Mutex.Lock()

	defer h.Mutex.Unlock()

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

func (h *Hub) JoinRoom(w http.ResponseWriter, r *http.Request, arg utils.JoinRoomReq) error {
	h.Mutex.Lock()

	defer h.Mutex.Unlock()
	
	conn, wsErr := ws.Upgrade(w, r, nil)

	if wsErr != nil {
		return fmt.Errorf("Web socket upgrade error:", wsErr)
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
		RoomID:   arg.RoomID,
		Username: arg.Username,
	}

	// Register a client through register channel

	// Broadcast the message

	// Write Message

	// Read Message
}