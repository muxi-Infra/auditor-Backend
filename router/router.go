package router

import (
	"github.com/gin-gonic/gin"
	"muxi_auditor/controller"
	"muxi_auditor/middleware"
)

func NewRouter(
	OAuth *controller.AuthController,
	AuthMiddleware *middleware.AuthMiddleware,
	corsMiddleware *middleware.CorsMiddleware,
) *gin.Engine {

	r := gin.New()
	//使用gin的Panic捕获中间件
	r.Use(gin.Recovery())

	// 添加 CORS 中间件,跨域中间件
	r.Use(corsMiddleware.Fc)

	g := r.Group("/api/v1")
	RegisterOAuthRoutes(g, AuthMiddleware.Fc, OAuth)

	return r
}
