package router

import (
	"github.com/gin-gonic/gin"
	"muxi_auditor/api/request"
	"muxi_auditor/api/response"
	"muxi_auditor/pkg/ginx"
	"muxi_auditor/pkg/jwt"
)

type OAuthController interface {
	Login(g *gin.Context, req request.LoginReq) (response.Response, error)

	Logout(g *gin.Context) (response.Response, error)
	UpdateMyInfo(g *gin.Context, req request.UpdateUserReq) (response.Response, error)
	GetQiToken(g *gin.Context) (response.Response, error)
	GetMyInfo(g *gin.Context, cla *jwt.UserClaims) (response.Response, error)
}

func RegisterOAuthRoutes(
	s *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	c OAuthController,
) {
	//认证服务
	authGroup := s.Group("/user")
	authGroup.POST("/login", ginx.WrapReq(c.Login))
	authGroup.GET("/logout", authMiddleware, ginx.Wrap(c.Logout))
	authGroup.POST("/updateMyInfo", authMiddleware, ginx.WrapReq(c.UpdateMyInfo))
	authGroup.GET("/GetQiToken", authMiddleware, ginx.Wrap(c.GetQiToken))
	authGroup.GET("/getMyInfo", authMiddleware, ginx.WrapClaims(c.GetMyInfo))
}
