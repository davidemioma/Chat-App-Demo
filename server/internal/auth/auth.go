package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var key = "chat_app_jwt_token"

func GetAuthToken (r *http.Request) (string, error){
	// Retrieve the JWT token from the cookie
	cookie, err := r.Cookie(key)

	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("cookie not found: %v", err)
		}

        return "", fmt.Errorf("server cookie error: %v", err)
	}

	token := cookie.Value

	if len(token) == 0 {
		return "", fmt.Errorf("cookie is empty")
	}

	return token, nil
}

func SetAuthToken(c *gin.Context, token string) {
	c.SetCookie(key, token, 86400, "/", "localhost", false, true)
}