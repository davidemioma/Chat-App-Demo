package main

import (
	"encoding/json"
	"net/http"
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
	http.SetCookie(w, &http.Cookie{
		Name:     "chat_app_jwt_token",
		Value:    ss,
		Path:     "/api",
		HttpOnly: true,
		Secure:   false, // false for development, but true in production.
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(time.Hour * 24 / time.Second), // 24 hours
	})
	

	utils.RespondWithJSON(w, http.StatusOK, utils.JsonUser{
		ID: userExists.ID.String(),
		Email: userExists.Email,
		Username: userExists.Email,
		CreatedAt: userExists.CreatedAt,
		UpdatedAt: userExists.UpdatedAt,
	})
}



func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the JWT token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "chat_app_jwt_token",
		Value:    "",
		Path:     "/api",
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
	})

	utils.RespondWithJSON(w, http.StatusOK, "Logged out successfully")
}