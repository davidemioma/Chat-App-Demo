package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespondWithJSON(c *gin.Context, code int, payload interface{}) {
    _, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})

		return
	}

	c.JSON(code, payload)
}

func RespondWithError(c *gin.Context, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}
	
	RespondWithJSON(c, code, errorResponse{
		Error: msg,
	})
}
