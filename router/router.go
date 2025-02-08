package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"muxi_auditor/controller"
	"muxi_auditor/middleware"
)

func NewRouter(
	OAuth *controller.AuthController,
	User *controller.UserController,
	AuthMiddleware *middleware.AuthMiddleware,
	corsMiddleware *middleware.CorsMiddleware,
	loggerMiddleware *middleware.LoggerMiddleware,
	Project *controller.ProjectController,

) *gin.Engine {

	r := gin.New()
	//使用gin的Panic捕获中间件
	r.Use(gin.Recovery())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// 添加 CORS 中间件,跨域中间件
	r.Use(corsMiddleware.MiddlewareFunc())
	// 添加日志与打点中间件
	r.Use(loggerMiddleware.MiddlewareFunc())
	//暴露给前端的api前缀
	g := r.Group("/api/v1")
	//注册router
	RegisterOAuthRoutes(g, AuthMiddleware.MiddlewareFunc(), OAuth)
	UserRoutes(g, AuthMiddleware.MiddlewareFunc(), User)
	RegisterProjectRoutes(g, AuthMiddleware.MiddlewareFunc(), Project)
	return r
}
