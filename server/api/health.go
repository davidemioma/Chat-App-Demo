package main

import (
	"net/http"
	"server/utils"

	"github.com/gin-gonic/gin"
)


func handlerReadiness(c *gin.Context) { 
	utils.RespondWithJSON(c, http.StatusOK, map[string]string{"status": "ok"}) 
}

func handlerErr(c *gin.Context) { 
	utils.RespondWithError(c, http.StatusInternalServerError, "Internal Server Error") 
}