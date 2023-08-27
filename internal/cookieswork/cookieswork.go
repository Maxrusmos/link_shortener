package cookieswork

import (
	"net/http"
	"time"
)

var secretKey = []byte("123")
var cookieName = "auth_cookie"

func GenerateCookie(userID string) (*http.Cookie, error) {
	cookieValue := []byte(userID)
	cookie := &http.Cookie{
		Name:     "auth_cookie",
		Value:    string(cookieValue),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	return cookie, nil
}

func GetUserIDFromCookie(req *http.Request) (string, error) {
	cookie, err := req.Cookie("auth_cookie")
	if err != nil {
		return "", err
	}
	return string(cookie.Value), nil
}
