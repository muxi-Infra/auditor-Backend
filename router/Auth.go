package router

import (
	"github.com/gin-gonic/gin"
	"muxi_auditor/api/request"
	"muxi_auditor/pkg/ginx"
)

type OAuthController interface {
	Login(g *gin.Context, req request.LoginReq) (ginx.Response, error)
	//Logout(g *gin.Context) (ginx.Response, errors)
}

func RegisterOAuthRoutes(
	s *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	c OAuthController,
) {
	//认证服务
	authGroup := s.Group("/auth")
	authGroup.POST("/login", ginx.WrapReq(c.Login))
	//authGroup.POST("/logout", authMiddleware, ginx.Wrap(c.Logout))
}
