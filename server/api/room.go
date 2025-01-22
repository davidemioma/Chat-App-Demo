package main

import (
	"encoding/json"
	"net/http"
	"server/utils"

	"github.com/google/uuid"
)

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