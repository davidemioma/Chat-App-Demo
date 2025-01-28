package main

import (
	"database/sql"
	"log"
	"net/http"
	"server/internal/database"
	"server/internal/socket"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type config struct {
	db *sql.DB
	dbQuery *database.Queries
	jwtSecret string
}

type application struct {
	config  config
	logger  *log.Logger
	hub     socket.Hub
}

var r *gin.Engine

// Handle Routes
func (app *application) mount() http.Handler {
	r = gin.Default()

	// Cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With", "X-Token"},
		ExposeHeaders:    []string{"Content-Length", "Link", "Set-Cookie"},
		AllowCredentials: true,  // Change this to true to allow cookies
		MaxAge:           time.Hour * 12,
	}))

	// Routes
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/health", handlerReadiness)

		apiGroup.GET("/err", handlerErr)

		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/sign-up", app.registerHandler)

			authGroup.POST("/sign-in", app.loginHandler)

			authGroup.GET("/sign-out", app.logoutHandler)

			authGroup.GET("/user", app.middlewareAuth(app.getCurrentUser))
		}

		roomsGroup := apiGroup.Group("/rooms")
		{
			roomsGroup.GET("/", app.getRoomHandler)

			roomsGroup.POST("/create", app.createRoomHandler)

			roomsGroup.GET("/:roomId/join", app.joinRoomHandler)

			roomsGroup.GET("/:roomId/clients", app.getClientsHandler)
		}
	}

	return r
}

// Run Server
func (app *application) run(port string, handler http.Handler) error {
	srv := &http.Server{
		Addr:         port,
		Handler:      handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	err := srv.ListenAndServe()

	if err != nil {
		return err
	}

	return nil
}