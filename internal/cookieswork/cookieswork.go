package cookieswork

import (
	"net/http"
	"time"
)

const (
	authCookieName = "auth_"
	// Здесь лучше использовать более безопасный метод для генерации секретного ключа.
	authSecret = "123"
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
