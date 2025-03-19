package apikey

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SecretKey = []byte("Muxi-Team-Auditor-Backend")

func GenerateAPIKey(projectID uint) (string, error) {
	claims := jwt.MapClaims{
		"sub": projectID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(1000 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	apiKey, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return apiKey, nil
}
func ParseAPIKey(apiKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(apiKey, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}

}
