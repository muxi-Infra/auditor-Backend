package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

// JWT 实现了处理 JWT 的接口，并使用 Redis 进行支持
type JWT struct {
	signingMethod jwt.SigningMethod // JWT 签名方法
	rcExpiration  time.Duration     // 刷新令牌的过期时间，防止缓存过大
	jwtKey        []byte            // 用于签署 JWT 的密钥
}

// NewRedisJWTHandler 创建并返回一个新的 JWT 实例
func NewJWT(expiration time.Duration, secretKey string) *JWT {

	return &JWT{
		signingMethod: jwt.SigningMethodHS256, //签名的加密方式
		rcExpiration:  expiration,
		jwtKey:        []byte(secretKey),
	}
}

// ExtractToken 从请求中提取并返回解析完成的结构体
func (j *JWT) ParseToken(tokenStr string) (UserClaims, error) {
	//解析token
	uc := UserClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &uc, func(*jwt.Token) (interface{}, error) {
		// 可以根据具体情况给出不同的key
		return j.jwtKey, nil
	})

	if j.jwtKey == nil {
		return UserClaims{}, errors.New("jwtKey 为空，无法解析 token")
	}
	if err != nil {
		return UserClaims{}, err
	}

	//检查有效性
	if token == nil || !token.Valid {
		return UserClaims{}, errors.New("token无效")
	}

	return uc, nil
}

// SetJWTToken 生成并设置用户的 JWT
func (j *JWT) SetJWTToken(uid uint, email string, userRole int) (string, error) {

	uc := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.rcExpiration)),
			ID:        uuid.New().String(),
		},
		Uid:      uid,
		Email:    email,
		UserRule: userRole,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenStr, err := token.SignedString(j.jwtKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// UserClaims 定义了 JWT 中用户相关的声明
type UserClaims struct {
	jwt.RegisteredClaims
	Uid      uint   // 用户 ID
	Email    string //用户名称
	UserRule int
}
