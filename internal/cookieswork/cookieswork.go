package cookieswork

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	authCookieName = "authgjdgsd1ssssasass2asassds"
	authSecret     = "123"
)

func SetAuthCookie(w http.ResponseWriter, r *http.Request, userID string) {
	cookie := &http.Cookie{
		Name:     authCookieName,
		Value:    userID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
	r.AddCookie(cookie)
}

func GetUserID(r *http.Request) string {
	cookie, err := r.Cookie(authCookieName)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return cookie.Value
}

func IsAuthenticated(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(authCookieName)
	if err != nil || cookie.Value == "" {
		userID := generateUniqueUserID()
		SetAuthCookie(w, r, userID)
	}
}

func generateUniqueUserID() string {
	return uuid.New().String()
}

func generateCookieName() string {
	return "cookie_" + uuid.New().String()
}
