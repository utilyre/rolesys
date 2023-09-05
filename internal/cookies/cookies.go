package cookies

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

func Write(w http.ResponseWriter, cookie *http.Cookie) error {
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}

	http.SetCookie(w, cookie)
	return nil
}

func Read(r *http.Request, name string) (*http.Cookie, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return nil, err
	}

	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, ErrInvalidValue
	}

	cookie.Value = string(value)
	return cookie, nil
}

func WriteSigned(w http.ResponseWriter, cookie *http.Cookie, secret []byte) error {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(cookie.Value))
	signature := mac.Sum(nil)

	cookie.Value = string(signature) + cookie.Value
	return Write(w, cookie)
}

func ReadSigned(r *http.Request, name string, secret []byte) (*http.Cookie, error) {
	cookie, err := Read(r, name)
	if err != nil {
		return nil, err
	}

	if len(cookie.Value) < sha256.Size {
		return nil, ErrInvalidValue
	}

	signature := cookie.Value[:sha256.Size]
	cookie.Value = cookie.Value[sha256.Size:]

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(cookie.Value))
	if !hmac.Equal([]byte(signature), mac.Sum(nil)) {
		return nil, ErrInvalidValue
	}

	return cookie, nil
}

func WriteEncrypted(w http.ResponseWriter, cookie *http.Cookie, secret []byte) error {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return err
	}

	plaintext := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)
	value := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	cookie.Value = string(value)
	return Write(w, cookie)
}

func ReadEncrypted(r *http.Request, name string, secret []byte) (*http.Cookie, error) {
	cookie, err := Read(r, name)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(cookie.Value) < nonceSize {
		return nil, ErrInvalidValue
	}

	nonce := []byte(cookie.Value[:nonceSize])
	ciphertext := []byte(cookie.Value[nonceSize:])

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrInvalidValue
	}

	expectedName, value, ok := strings.Cut(string(plaintext), ":")
	if !ok {
		return nil, ErrInvalidValue
	}

	if expectedName != name {
		return nil, ErrInvalidValue
	}

	cookie.Value = value
	return cookie, nil
}
