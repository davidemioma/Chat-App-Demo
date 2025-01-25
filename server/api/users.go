package main

import (
	"net/http"
	"server/internal/auth"
	"server/internal/database"
	"server/utils"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v4"

	"github.com/google/uuid"
)

func (app *application) registerHandler(c *gin.Context) {
	type parameters struct {
		Email       string  `json:"email"`
		Username    string  `json:"username"`
		Password    string  `json:"password"`
	}

	// Validating body
	var params parameters

	if err := c.ShouldBindJSON(&params); err != nil {
		app.logger.Printf("Error parsing JSON: %v", err)
		
		utils.RespondWithError(c, http.StatusBadRequest, "Error parsing JSON")

		return
	}

	// Check if parameters is valid
	if params.Email == "" || params.Username == "" || params.Password == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid Parameters!")

		return
	}

	// Check if user already exists
	userExists, existsErr := app.config.dbQuery.CheckUser(c.Request.Context(), params.Email)

	if existsErr == nil && userExists != (database.CheckUserRow{}) {
		utils.RespondWithError(c, http.StatusUnauthorized, "Email already exists!")
		
		return
	}

	// Hash Password
	hashedPassword, hashErr := utils.HashPassword(params.Password)

	if hashErr != nil{
		app.logger.Printf("Error hashing password: %v", hashErr)

		utils.RespondWithError(c, http.StatusInternalServerError, "Unable to encrypt password")
		
		return
	}

	// Create user
	dbErr := app.config.dbQuery.CreateUser(c.Request.Context(), database.CreateUserParams{
		ID: uuid.New(),
		Email: params.Email,
		Username: params.Username,
		Hashedpassword: hashedPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if dbErr != nil {
		app.logger.Printf("Couldn't create user: %v", dbErr)

		utils.RespondWithError(c, http.StatusInternalServerError, "Couldn't create user")

		return
	}

	utils.RespondWithJSON(c, http.StatusCreated, "New user created")
}

func (app *application) loginHandler(c *gin.Context) {
	type parameters struct {
		Email       string  `json:"email"`
		Password    string  `json:"password"`
	}

	// Validating body
	var params parameters

	if err := c.ShouldBindJSON(&params); err != nil {
		app.logger.Printf("Error parsing JSON: %v", err)
		
		utils.RespondWithError(c, http.StatusBadRequest, "Error parsing JSON")

		return
	}

	// Check if parameters is valid
	if params.Email == "" || params.Password == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid Parameters!")

		return
	}

	// Check if user exists
	userExists, existsErr := app.config.dbQuery.GetUserByEmail(c.Request.Context(), params.Email)

	if existsErr != nil || userExists == (database.User{}) {
		app.logger.Printf("User not found: %v", existsErr)

		utils.RespondWithError(c, http.StatusNotFound, "User not found!")
		
		return
	}

	// Check if password matched
	passErr := utils.CheckPassword(params.Password, userExists.Hashedpassword)

	if passErr != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "Password does not match!")
		
		return
	}

	// Set up JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.MyJWTClaims{
		ID: userExists.ID.String(),
		Email: userExists.Email,
		Username: userExists.Username,
		CreatedAt: userExists.CreatedAt,
		UpdatedAt: userExists.UpdatedAt,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: userExists.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 24 hrs
		},
	})

	ss, signErr := token.SignedString([]byte(app.config.jwtSecret))

	if signErr != nil {
		app.logger.Printf("Unable to sign JWT token: %v", signErr)

		utils.RespondWithError(c, http.StatusInternalServerError, "Something went wrong! try again.")
		
		return
	}

	// Store token in http cookies
	auth.SetAuthToken(c, ss)
	
	utils.RespondWithJSON(c, http.StatusOK, "Login Successful")
}

func (app *application) logoutHandler(c *gin.Context) {
	// Clear the JWT token cookie
	auth.SetAuthToken(c, "")

	utils.RespondWithJSON(c, http.StatusOK, "Logged out successfully")
}

func (app *application) getCurrentUser(c *gin.Context, user utils.JsonUser) {
	utils.RespondWithJSON(c, http.StatusOK, user)
}