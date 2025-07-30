package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetKey(headers http.Header, scheme string) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("no authorization key found")
	}
	authslice := strings.SplitN(auth, " ", 2)
	if len(authslice) != 2 || authslice[0] != scheme {
		return "", errors.New("invalid authorization")
	}
	return authslice[1], nil
}
