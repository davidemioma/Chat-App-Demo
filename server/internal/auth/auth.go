package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var key = "chat_app_jwt_token"

func GetAuthToken (c *gin.Context) (string, error){
	// Retrieve the JWT token from the cookie
	cookie, err := c.Cookie(key)

	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("cookie not found: %v", err)
		}

        return "", fmt.Errorf("server cookie error: %v", err)
	}

	if len(cookie) == 0 {
		return "", fmt.Errorf("cookie is empty")
	}

	return cookie, nil
}

func SetAuthToken(c *gin.Context, token string) {
	c.SetCookie(key, token, 86400, "/", "", false, true)
}