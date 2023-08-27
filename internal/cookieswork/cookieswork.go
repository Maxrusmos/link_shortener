package cookieswork

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

var cookieName = "auth_cookie"

func GenerateCookie(id string) (*http.Cookie, error) {
	value, err := json.Marshal(id)
	if err != nil {
		return nil, err
	}
	signature := hmac.New(sha256.New, []byte("secret_key"))
	signature.Write(value)
	signedValue := signature.Sum(nil)
	cookie := &http.Cookie{
		Name:  "auth_cookie",
		Value: base64.StdEncoding.EncodeToString(value) + "|" + base64.StdEncoding.EncodeToString(signedValue),
		Path:  "/",
	}
	return cookie, nil
}
