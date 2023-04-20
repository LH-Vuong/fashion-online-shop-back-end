package grant

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateAPIKey generates an API key and secret using cryptographic random numbers, encoded in base64.
func GenerateAPIKey() (apiKey, apiSecret string, err error) {
	key := make([]byte, 16)
	_, err = rand.Read(key)
	if err != nil {
		return "", "", err
	}

	secret := make([]byte, 32)
	_, err = rand.Read(secret)
	if err != nil {
		return "", "", err
	}
	apiKey = base64.StdEncoding.EncodeToString(key)
	apiSecret = base64.StdEncoding.EncodeToString(secret)
	return apiKey, apiSecret, nil
}
