package apikey

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SecretKey = []byte("Muxi-Team-Auditor-Backend")

func GenerateAPIKey(projectID uint) (string, error) {
	claims := jwt.MapClaims{
		"sub": projectID,
		"iat": time.Now().Unix(),
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

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
func GenerateKeyPair() (accessKey, secretKey string) {
	accessKey = "cli_" + randomString(16)
	secretKey = randomString(32) // 只返回给调用方一次
	return
}
func SignRequest(secret, timestamp string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	return hex.EncodeToString(mac.Sum(nil))
}
