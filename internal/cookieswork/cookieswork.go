package cookieswork

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	authCookieName = "auth"
	authSecret     = "123"
)

func SetAuthCookie(w http.ResponseWriter, userID string) {
	cookie := &http.Cookie{
		Name:     authCookieName,
		Value:    userID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

func GetUserID(r *http.Request) string {
	cookie, err := r.Cookie(authCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func IsAuthenticated(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(authCookieName)
	if err != nil || cookie.Value == "" {
		userID := generateUniqueUserID()
		SetAuthCookie(w, userID)
	}
}

func generateUniqueUserID() string {
	return uuid.New().String()
}
