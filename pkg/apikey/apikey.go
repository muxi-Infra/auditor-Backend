package apikey

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SecretKey = []byte("Muxi-Team-Auditor-Backend")

func GenerateAPIKey(projectID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": projectID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	apiKey, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return apiKey, nil
}
func ParseAPIKey(apiKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(apiKey, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return SecretKey, nil
	})

	return token, err
}
