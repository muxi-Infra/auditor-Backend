package mytoken
import (
"fmt"
"github.com/golang-jwt/jwt/v5"
"time"
)

var key = []byte("Muxishenhe")

func GenerateJwt(accesstoken string) (string, error) {
	claims := jwt.MapClaims{
		"Atoken": accesstoken,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)
}
func ParseToken(tokenString string) (*jwt.Token, error) {
	// 解析 token，并提供密钥来验证签名
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法是我们期望的签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// 返回密钥
		return key, nil
	})

	if err != nil {
		return nil, err
	}


	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Claims[Atoken]:", claims["Atoken"])
		fmt.Println("Claims[exp]:", claims["exp"])
	} else {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

