package auth

import (
	"net/http"
)

func GetAPIKey(headers http.Header) (string, error) {
	key, err := GetKey(headers, "ApiKey")
	return key, err
}
