package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"

	"github.com/muxi-Infra/auditor-Backend/api/request"
	"github.com/muxi-Infra/auditor-Backend/api/response"
	"github.com/muxi-Infra/auditor-Backend/repository/model"
	"github.com/muxi-Infra/auditor-Backend/service"
)

type RemoveController struct {
	service RemoveItemService
}
type RemoveItemService interface {
	CheckPower(c context.Context, apikey string) (bool, uint, error)
	Upload(c context.Context, req request.UploadReq, projectId uint) (uint, error)
	Update(c context.Context, req request.RemoveUpdateReq, projectId uint) (uint, error)
	Delete(c context.Context, itemId uint, projectId uint) error
	Get(c context.Context, itemIds []uint, projectId uint) (*model.RemoveItemsStatus, error)
}

// NewRemoveController remove方面的控制器，处理http请求与结果处理，数据库方面逻辑交给service
func NewRemoveController(service *service.RemoveService) *RemoveController {
	return &RemoveController{
		service: service,
	}
}

// Upload sdk上传项目
// @Summary sdk上传项目，无需对接
// @Description 通过使用提供的sdk供调用方快速上传项目
// @Tags Remove
// @Accept json
// @Produce json
// @Param api_key header string true "访问凭证 ApiKey"
// @Param UploadReq body request.UploadReq true "上传请求体"
// @Success 200 {object} response.Response "成功返回项目id"
// @Failure 400 {object} response.Response "获取项目列表失败"
// @Router /api/v1/remove/upload [post]
func (c *RemoveController) Upload(ctx *gin.Context, req request.UploadReq) (response.Response, error) {
	id, err := c.CheckPower(ctx)
	if err != nil {
		return response.Response{
			Code: http.StatusBadRequest,
			Msg:  fmt.Errorf("power check error:%w", err).Error(),
		}, err
	}

	//上传item的逻辑
	itemId, err := c.service.Upload(ctx, req, id)
	if err != nil {
		var re = response.Response{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Errorf("%d:%w", itemId, err).Error(),
		}
		return re, err
	}
	re := response.Response{
		Code: http.StatusOK,
		Data: itemId,
	}
	return re, nil

}

// Update sdk更新项目
// @Summary sdk更新项目，无需对接
// @Description 通过使用提供的sdk供调用方快速更改项目信息
// @Tags Remove
// @Accept json
// @Produce json
// @Param api_key header string true "访问凭证 ApiKey"
// @Param UploadReq body request.UploadReq true "上传请求体"
// @Success 200 {object} response.Response "成功返回项目id"
// @Failure 400 {object} response.Response "修改失败"
// @Router /api/v1/remove/update [put]
func (c *RemoveController) Update(ctx *gin.Context, req request.RemoveUpdateReq) (response.Response, error) {
	id, err := c.CheckPower(ctx)
	if err != nil {
		return response.Response{Code: http.StatusBadRequest, Msg: fmt.Errorf("power check err:%w", err).Error()}, err
	}
	itemId, err := c.service.Update(ctx, req, id)
	if err != nil {
		var re = response.Response{
			Code: http.StatusInternalServerError,
			Data: fmt.Errorf("%d:%w", itemId, err).Error(),
		}
		return re, err
	}
	re := response.Response{
		Code: http.StatusOK,
		Data: itemId,
	}
	return re, nil
}

// Delete sdk删除项目
// @Summary sdk删除项目，无需对接
// @Description 通过使用提供的sdk供调用方快速删除项目
// @Tags Remove
// @Accept json
// @Produce json
// @Param api_key header string true "访问凭证 ApiKey"
// @Param Itemid path uint true "要删除的项目ID（Itemid）"
// @Success 200 {object} response.Response "成功返回删除的项目id"
// @Failure 400 {object} response.Response "删除项目失败"
// @Router /api/v1/remove/delete/{Itemid} [delete]
func (c *RemoveController) Delete(ctx *gin.Context) (response.Response, error) {
	id, err := c.CheckPower(ctx)
	if err != nil {
		return response.Response{Code: http.StatusBadRequest, Msg: fmt.Errorf("power check err:%w", err).Error()}, err
	}
	//获取参数,其实是hook_id
	data := ctx.Param("Itemid")
	itemId, err := strconv.ParseUint(data, 10, 64)
	if err != nil {
		return response.Response{Code: http.StatusBadRequest, Msg: fmt.Errorf("item id err:%w", err).Error()}, err
	}

	//删除逻辑
	err = c.service.Delete(ctx, uint(itemId), id)
	if err != nil {
		var re = response.Response{
			Code: http.StatusInternalServerError,
			Data: fmt.Errorf("%d:%w", itemId, err).Error(),
		}
		return re, err
	}
	re := response.Response{
		Code: http.StatusOK,
		Data: itemId,
	}
	return re, nil
}

// Get sdk获取项目信息
// @Summary sdk获取项目信息，无需对接
// @Description 通过使用提供的sdk供调用方快速获取项目信息，如果只传一个0表示获取全部
// @Tags Remove
// @Accept json
// @Produce json
// @Param api_key header string true "访问凭证 ApiKey"
// @Param ids query string true "项目ID列表，多个ID用英文逗号分隔，如: 1,2,3"
// @Success 200 {object} response.Response "成功返回项目信息"
// @Failure 400 {object} response.Response "获取项目失败"
// @Router /api/v1/remove/get [get]
func (c *RemoveController) Get(ctx *gin.Context) (response.Response, error) {
	idsStr := ctx.Query("ids")        // 获取字符串 "1,2,3"
	ids := strings.Split(idsStr, ",") // 手动切割成字符串数组
	id, err := c.CheckPower(ctx)
	if err != nil {
		return response.Response{Code: http.StatusBadRequest, Msg: fmt.Errorf("power check err:%w", err).Error()}, err
	}
	ItemIds, err := stringSliceToUintSlice(ids)
	if err != nil {
		return response.Response{}, err
	}
	items, err := c.service.Get(ctx, ItemIds, id)
	if err != nil {
		var re = response.Response{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
			Data: nil,
		}
		return re, err
	}
	re := response.Response{
		Code: http.StatusOK,
		Data: items,
	}
	return re, nil
}

func stringSliceToUintSlice(strs []string) ([]uint, error) {
	var result []uint
	for _, s := range strs {
		num, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, uint(num))
	}
	return result, nil
}

// CheckPower 鉴权控制器
func (c *RemoveController) CheckPower(ctx *gin.Context) (uint, error) {
	ac := ctx.GetHeader("api_key")
	if ac == "" {
		return 0, errors.New("http header parameters required")
	}
	//鉴权,返回project_id
	ok, id, err := c.service.CheckPower(ctx, ac)
	if !ok || err != nil {
		return 0, err
	}
	return id, nil
}
