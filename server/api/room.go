package main

import (
	"encoding/json"
	"net/http"
	"server/utils"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (app *application) getRoomHandler(w http.ResponseWriter, r *http.Request) {
	rooms := app.hub.GetRooms()

	utils.RespondWithJSON(w, http.StatusOK, rooms)
}

func (app *application) getClientsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the room Id from the URL params
    roomId := chi.URLParam(r, "roomId")

	if roomId == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Room ID required!")

		return
	}

	clients, err := app.hub.GetClients(roomId)

	if err != nil {
		app.logger.Printf("Couldn't get clients: %v", err)

		utils.RespondWithError(w, http.StatusNotFound, "Couldn't get clients! room not found.")

		return
	}

	utils.RespondWithJSON(w, http.StatusOK, clients)
}

func (app *application) createRoomHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	// Validating body
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		app.logger.Printf("Error parsing JSON: %v", err)
		
		utils.RespondWithError(w, http.StatusBadRequest, "Error parsing JSON")

		return
	}

	// Check if parameters is valid
	if params.Name == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Name is required!")

		return
	}

	// Create Room
	roomErr := app.hub.CreateRoom(r.Context(), utils.CreateRoomReq{
		ID: uuid.New().String(),
		Name: params.Name,
	})

	if roomErr != nil {
		app.logger.Printf("Couldn't create room: %v", roomErr)

		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create room")

		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, "New room created")
}

func (app *application) joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	// Get the room Id from the URL params
    roomId := chi.URLParam(r, "roomId")

	// Get parameters
	type parameters struct {
		UserID   string `json:"userId"`
        Username string `json:"username"`
	}

	// Validating body
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		app.logger.Printf("Error parsing JSON: %v", err)
		
		utils.RespondWithError(w, http.StatusBadRequest, "Error parsing JSON")

		return
	}

	// Check if parameters is valid
	if roomId == "" || params.UserID == "" || params.Username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Parameters!")

		return
	}

	// Join Room
	joinErr := app.hub.JoinRoom(w, r, utils.JoinRoomReq{
		RoomID: roomId,
		UserID: params.UserID,
		Username: params.Username,
	})

	if joinErr != nil {
		app.logger.Printf("Couldn't join room: %v", joinErr)

		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't join room")

		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "Successfully joined room")
}