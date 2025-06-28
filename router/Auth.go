package router

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/ginx"
	"github.com/gin-gonic/gin"
)

// OAuthController 登录登出接口
type OAuthController interface {
	Login(g *gin.Context, req request.LoginReq) (response.Response, error)

	Logout(g *gin.Context) (response.Response, error)
}

func RegisterOAuthRoutes(
	s *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	c OAuthController,
) {
	//认证服务
	authGroup := s.Group("/auth")
	authGroup.POST("/login", ginx.WrapReq(c.Login))
	authGroup.GET("/logout", authMiddleware, ginx.Wrap(c.Logout))

}
