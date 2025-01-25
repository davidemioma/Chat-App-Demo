package main

import (
	"fmt"
	"net/http"
	"server/internal/auth"
	"server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthHandler func (*gin.Context, utils.JsonUser)

func (app *application) middlewareAuth(handler AuthHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := auth.GetAuthToken(c.Request)

		if err != nil {
			app.logger.Printf("Couldn't get token: %v", err)

			utils.RespondWithError(c, http.StatusNotFound, "Couldn't get token")

			return
		}

		// Verify Token
		token, verifyErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(app.config.jwtSecret), nil 
		})

		if verifyErr != nil || !token.Valid {
			app.logger.Printf("Couldn't validate token: %v", verifyErr)

			utils.RespondWithError(c, http.StatusUnauthorized, "Couldn't validate token")

			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			app.logger.Printf("Invalid token claims")

			utils.RespondWithError(c, http.StatusUnauthorized, "Invalid token claims")
			
			return
		}

		// Get User Info
		user :=  utils.JsonUser{
			ID: claims["ID"].(string)  ,
			Email: claims["Email"].(string) ,
			Username: claims["Username"].(string),
			CreatedAt: claims["CreatedAt"].(time.Time),
            UpdatedAt: claims["UpdatedAt"].(time.Time),
		}

		handler(c, user)
	}
}