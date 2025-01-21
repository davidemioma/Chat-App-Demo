package main

import (
	"encoding/json"
	"net/http"
	"server/internal/database"
	"server/utils"
	"time"

	"github.com/google/uuid"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email       string  `json:"email"`
		Username    string  `json:"username"`
		Password    string  `json:"password"`
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
	if params.Email == "" || params.Username == "" || params.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Paramenters!")

		return
	}

	// Check if user already exists
	userExists, existsErr := app.config.dbQuery.CheckUser(r.Context(), params.Email)

	if existsErr == nil && userExists != (database.CheckUserRow{}) {
		utils.RespondWithJSON(w, http.StatusOK, "Email already exists!")
		
		return
	}

	// Hash Password
	hashedPassword, hashErr := utils.HashPassword(params.Password)

	if hashErr != nil{
		app.logger.Printf("Error hashing password: %v", hashErr)

		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to encrypt password")
		
		return
	}

	// Create user
	dbErr := app.config.dbQuery.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		Email: params.Email,
		Username: params.Username,
		Hashedpassword: hashedPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if dbErr != nil {
		app.logger.Printf("Couldn't create user: %v", dbErr)

		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create user")

		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, "New user created")
}