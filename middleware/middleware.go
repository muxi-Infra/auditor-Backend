package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"muxi_auditor/pkg/jwt"
	"net/http"
	"time"
)

type CorsMiddleware struct {
	Fc gin.HandlerFunc
}

type AuthMiddleware struct {
	Fc gin.HandlerFunc
}

func NewCorsMiddleware() *CorsMiddleware {
	return &CorsMiddleware{Fc: cors.New(cors.Config{
		// 允许的请求头
		AllowHeaders: []string{"Content-ContentType", "Authorization", "Origin"},
		// 是否允许携带凭证（如 Cookies）
		AllowCredentials: true,
		// 解决跨域问题,这个地方允许所有请求跨域了,之后要改成允许前端的请求,比如localhost
		AllowOriginFunc: func(origin string) bool {
			//检测请求来源是否以localhost开头
			//if strings.HasPrefix(origin, "localhost") {
			//	return true
			//}else {
			//	return false
			//}
			return true
		},

		// 预检请求的缓存时间
		MaxAge: 12 * time.Hour,
	})}
}

func NewAuthMiddleware(jwtHandler *jwt.RedisJWTHandler) *AuthMiddleware {

	return &AuthMiddleware{Fc: func(ctx *gin.Context) {
		// 从请求中提取并解析 Token
		userClaims, err := jwtHandler.ParseToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"errors": "Unauthorized: " + err.Error(),
			})
			ctx.Abort() // 中断请求
			return
		}

		// 将解析后的用户信息存入上下文，供后续逻辑使用
		ctx.Set("user", userClaims)

		// 继续处理请求
		ctx.Next()
	}}
}
