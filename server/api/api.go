package main

import (
	"database/sql"
	"log"
	"net/http"
	"server/internal/database"
	"server/internal/socket"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
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

// Handle Routes
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	
	r.Use(middleware.RealIP)

	r.Use(middleware.Logger)

	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	// Cors
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,  // Change this to true to allow cookies
		MaxAge:           300,
	}))

	// Routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", handlerReadiness)

		r.Get("/err", handlerErr)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", app.registerHandler)

			r.Post("/sign-in", app.loginHandler)

			r.Get("/sign-out", app.logoutHandler)

			r.Get("/user", app.middlewareAuth(app.getCurrentUser))
		})

		r.Route("/rooms", func(r chi.Router) {
			r.Get("/", app.getRoomHandler)

			r.Post("/create", app.createRoomHandler)

			r.Post("/{roomId}/join", app.joinRoomHandler)

			r.Get("/{roomId}/clients", app.getClientsHandler)
		})
	})

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