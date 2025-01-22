// Run go mod init <app name> to initialise app
// Run "echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc", "source ~/.zshrc" and "air" to start server.
// If you need a port, install "go get github.com/lpernett/godotenv", run "go mod vendor" and run "go mod tidy".
// To run a server, install "go get github.com/go-chi/chi" and "go get github.com/go-chi/cors", run "go mod vendor" and run "go mod tidy"

package main

import (
	"database/sql"
	"log"
	"os"
	"server/internal/database"

	"github.com/lpernett/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")

	// Port
	port := os.Getenv("PORT")

	if port == ""{
	    log.Fatal("PORT not found")
	}

	// JWT_Secret
	jwt_secret := os.Getenv("JWT_SECRET")

	if jwt_secret == ""{
	    log.Fatal("JWT secret not found")
	}

	// Postgres DB
	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == ""{
		log.Fatal("DATABASE_URL not found")
	}

	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	// Create a logger
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	cfg := config{
		db: db,
		dbQuery: database.New(db),
		jwtSecret: jwt_secret,
	}

	app := application{
		config: cfg,
		logger: logger,
	}

	// Intialise hub
	app.hub.Init()

	// Run hun channels
	go app.hub.Run()

	r:= app.mount()

	log.Printf("Server running on port %v", port)

	log.Fatal(app.run("0.0.0.0:" + port, r))
}