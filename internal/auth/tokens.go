package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	randnum := make([]byte, 32)
	_, err := rand.Read(randnum)
	if err != nil {
		return "", err
	}
	randstring := hex.EncodeToString(randnum)
	return randstring, nil
}
