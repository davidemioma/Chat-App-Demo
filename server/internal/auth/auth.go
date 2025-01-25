package auth

import (
	"fmt"
	"log"
	"net/http"
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

func SetAuthToken(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     key,
		Value:    token,
		Path:     "/",
		HttpOnly: true, // Prevents JavaScript access to the cookie
		Secure:   false, // Set to true in production (HTTPS)
		SameSite: http.SameSiteNoneMode,
		MaxAge:   86400, // 24 hours
	}

	http.SetCookie(w, &cookie)

	log.Println("Cookie set successfully")
}