package router

import (
	"github.com/gin-gonic/gin"

	"github.com/muxi-Infra/auditor-Backend/api/request"
	"github.com/muxi-Infra/auditor-Backend/api/response"
	"github.com/muxi-Infra/auditor-Backend/pkg/ginx"
)

type RemoveController interface {
	Upload(g *gin.Context, req request.UploadReq) (response.Response, error)
	Update(g *gin.Context, req request.RemoveUpdateReq) (response.Response, error)
	Get(g *gin.Context) (response.Response, error)
	Delete(g *gin.Context) (response.Response, error)
}

// RemoveRoutes 其他应用上传或修改item的接口
func RemoveRoutes(
	s *gin.RouterGroup,

	c RemoveController,
) {
	removeGroup := s.Group("/remove")

	removeGroup.POST("/upload", ginx.WrapReq(c.Upload))
	removeGroup.PUT("/update", ginx.WrapReq(c.Update))
	removeGroup.GET("/get", ginx.Wrap(c.Get))
	removeGroup.DELETE("/delete/:Itemid", ginx.Wrap(c.Delete))
}
