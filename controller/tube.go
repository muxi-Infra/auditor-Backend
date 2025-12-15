package controller

import (
	"context"
	"github.com/gin-gonic/gin"

	"github.com/muxi-Infra/auditor-Backend/api/response"
	"github.com/muxi-Infra/auditor-Backend/config"
	"github.com/muxi-Infra/auditor-Backend/service"
)

type TubeController struct {
	service TubeService
	qi      *config.QiNiuYunConfig
}
type TubeService interface {
	GetQiToken(ctx context.Context) (string, error)
}

func NewTuberController(service *service.TubeService) *TubeController {
	return &TubeController{
		service: service,
	}
}

// GetQiToken 获取七牛云上传Token
// @Summary 获取七牛云上传 Token
// @Description 获取用于上传文件的七牛云 Token
// @Tags Tube
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "返回上传 Token"
// @Failure 400 {object} response.Response "获取图床token失败"
// @Security ApiKeyAuth
// @Router /api/v1/tube/GetQiToken [get]
func (c *TubeController) GetQiToken(ctx *gin.Context) (response.Response, error) {
	token, err := c.service.GetQiToken(ctx)
	if err != nil {
		return response.Response{
			Code: 400,
			Data: nil,
			Msg:  "获取图床token失败",
		}, err
	}
	return response.Response{
		Code: 200,
		Data: token,
		Msg:  "获取成功",
	}, nil
}
