package cookieswork

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"time"
)

var secretKey = []byte("123")
var cookieName = "auth_cookie"

func GenerateCookie(userID string) (*http.Cookie, error) {
	cookieValue := []byte(userID)
	encryptedValue, err := encryptWithSecretKey(cookieValue, secretKey)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    base64.URLEncoding.EncodeToString(encryptedValue),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}

	return cookie, nil
}

func GetUserIDFromCookie(req *http.Request) (string, error) {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		return "", err
	}

	decryptedValue, err := decryptWithSecretKey([]byte(cookie.Value), secretKey)
	if err != nil {
		return "", err
	}

	return string(decryptedValue), nil
}

func decryptWithSecretKey(value, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("invalid key length")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(value) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := value[:aes.BlockSize]
	ciphertext := value[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	decrypted, err = unpadPKCS7(decrypted)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func unpadPKCS7(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("empty input")
	}

	padding := int(data[length-1])
	if padding > length || padding == 0 {
		return nil, errors.New("invalid padding")
	}

	for i := 1; i <= padding; i++ {
		if data[length-i] != byte(padding) {
			return nil, errors.New("invalid padding")
		}
	}

	return data[:length-padding], nil
}

func encryptWithSecretKey(value, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("invalid key length")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(value))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], value)

	return ciphertext, nil
}
