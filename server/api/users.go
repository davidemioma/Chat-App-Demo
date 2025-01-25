package main

import (
	"encoding/json"
	"net/http"
	"server/internal/auth"
	"server/internal/database"
	"server/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/google/uuid"
)

func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
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
		utils.RespondWithError(w, http.StatusUnauthorized, "Email already exists!")
		
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

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email       string  `json:"email"`
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
	if params.Email == "" || params.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Paramenters!")

		return
	}

	// Check if user exists
	userExists, existsErr := app.config.dbQuery.GetUserByEmail(r.Context(), params.Email)

	if existsErr != nil || userExists == (database.User{}) {
		utils.RespondWithError(w, http.StatusNotFound, "User not found!")
		
		return
	}

	// Check if password matched
	passErr := utils.CheckPassword(params.Password, userExists.Hashedpassword)

	if passErr != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Password does not match!")
		
		return
	}

	// Set up JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.MyJWTClaims{
		ID: userExists.ID.String(),
		Email: userExists.Email,
		Username: userExists.Email,
		CreatedAt: userExists.CreatedAt,
		UpdatedAt: userExists.UpdatedAt,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: userExists.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 24 hrs
		},
	})

	ss, signErr := token.SignedString([]byte(app.config.jwtSecret))

	if signErr != nil {
		app.logger.Printf("Unable to signed JWT token: %v", signErr)

		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong! try again.")
		
		return
	}

	// Store token in http cookies
	auth.SetAuthToken(w, ss)
	
	utils.RespondWithJSON(w, http.StatusOK, utils.JsonUser{
		ID: userExists.ID.String(),
		Email: userExists.Email,
		Username: userExists.Username,
		CreatedAt: userExists.CreatedAt,
		UpdatedAt: userExists.UpdatedAt,
	})
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the JWT token cookie
	auth.SetAuthToken(w, "")

	utils.RespondWithJSON(w, http.StatusOK, "Logged out successfully")
}

func (app *application) getCurrentUser(w http.ResponseWriter, r *http.Request, user utils.JsonUser) {
	utils.RespondWithJSON(w, http.StatusOK, user)
}