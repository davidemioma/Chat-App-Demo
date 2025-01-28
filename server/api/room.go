package main

import (
	"net/http"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (app *application) getRoomHandler(c *gin.Context) {
	rooms := app.hub.GetRooms()

	utils.RespondWithJSON(c, http.StatusOK, rooms)
}

func (app *application) getClientsHandler(c *gin.Context) {
	// Get the room Id from the URL params
	roomId := c.Param("roomId")

	if roomId == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Room ID required!")
		return
	}

	clients, err := app.hub.GetClients(roomId)

	if err != nil {
		app.logger.Printf("Couldn't get clients: %v", err)

		utils.RespondWithError(c, http.StatusNotFound, "Couldn't get clients! room not found.")

		return
	}

	utils.RespondWithJSON(c, http.StatusOK, clients)
}

func (app *application) createRoomHandler(c *gin.Context) {
	type parameters struct {
		Name string `json:"name"`
	}

	// Validating body
	var params parameters
	if err := c.ShouldBindJSON(&params); err != nil {
		app.logger.Printf("Error parsing JSON: %v", err)

		utils.RespondWithError(c, http.StatusBadRequest, "Error parsing JSON")

		return
	}

	// Check if parameters is valid
	if params.Name == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Name is required!")

		return
	}

	// Create Room
	roomId := uuid.New().String()

	roomErr := app.hub.CreateRoom(c.Request.Context(), utils.CreateRoomReq{
		ID: roomId,
		Name: params.Name,
	})

	if roomErr != nil {
		app.logger.Printf("Couldn't create room: %v", roomErr)

		utils.RespondWithError(c, http.StatusInternalServerError, "Couldn't create room")

		return
	}

	utils.RespondWithJSON(c, http.StatusCreated, utils.CreateRoomReq{
		ID: roomId,
		Name: params.Name,
	})
}

func (app *application) joinRoomHandler(c *gin.Context) {
	// Get the room Id from the URL params
	roomId := c.Param("roomId")

	// Get userId and username query params
	userId := c.Query("userId")

	username := c.Query("username") 

	// Check if parameters is valid
	if roomId == "" || userId == "" || username == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid Parameters!")
		
		return
	}

	// Join Room
	joinErr := app.hub.JoinRoom(c, utils.JoinRoomReq{
		RoomID: roomId,
		UserID: userId,
		Username: username,
	})

	if joinErr != nil {
		app.logger.Printf("Couldn't join room: %v", joinErr)

		utils.RespondWithError(c, http.StatusInternalServerError, "Couldn't join room")

		return
	}

	utils.RespondWithJSON(c, http.StatusOK, "Successfully joined room")
}