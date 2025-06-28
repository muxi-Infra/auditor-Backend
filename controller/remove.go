package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type RemoveController struct {
	service RemoveItemService
}
type RemoveItemService interface {
	CheckPower(c context.Context, ac, temp, signature string) (bool, uint, error)
	Upload(c context.Context, req request.UploadReq, projectId uint) (uint, error)
	Update(c context.Context, req request.UploadReq, projectId uint) (uint, error)
	Delete(c context.Context, itemId uint, projectId uint) error
	Get(c context.Context, itemIds []uint, projectId uint) ([]model.Item, error)
}

// NewRemoveController remove方面的控制器，处理http请求与结果处理，数据库方面逻辑交给service
func NewRemoveController(service RemoveItemService) *RemoveController {
	return &RemoveController{
		service: service,
	}
}
func (c *RemoveController) Upload(g *gin.Context, req request.UploadReq) (response.Response, error) {
	//ac := g.GetHeader("AccessKey")
	//temp := g.GetHeader("Timestamp")
	//signature := g.GetHeader("Signature")
	//if ac == "" || temp == "" || signature == "" {
	//	var re = response.Response{
	//		Code: http.StatusBadRequest,
	//		Msg:  "header参数缺失",
	//		Data: nil,
	//	}
	//
	//	return re, errors.New("http header parameters required")
	//}
	////鉴权
	//ok, id, err := c.service.CheckPower(g, ac, temp, signature)
	//if !ok || err != nil {
	//	var re = response.Response{
	//		Code: http.StatusBadRequest,
	//		Msg:  fmt.Errorf("鉴权失败%w", err).Error(),
	//	}
	//	return re, err
	//}
	id, err := c.CheckPower(g)
	if err != nil {
		return response.Response{
			Code: http.StatusBadRequest,
			Msg:  fmt.Errorf("power check error:%w", err).Error(),
		}, err
	}
	//上传item的逻辑
	itemId, err := c.service.Upload(g, req, id)
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
func (c *RemoveController) Update(g *gin.Context, req request.UploadReq) (response.Response, error) {
	id, err := c.CheckPower(g)
	if err != nil {
		return response.Response{Code: http.StatusBadRequest, Msg: fmt.Errorf("power check err:%w", err).Error()}, err
	}
	itemId, err := c.service.Update(g, req, id)
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
func (c *RemoveController) Delete(g *gin.Context) (response.Response, error) {
	id, err := c.CheckPower(g)
	if err != nil {
		return response.Response{Code: http.StatusBadRequest, Msg: fmt.Errorf("power check err:%w", err).Error()}, err
	}
	//获取参数,其实是hook_id
	data := g.Param("Itemid")
	itemId, err := strconv.ParseUint(data, 10, 64)
	if err != nil {
		return response.Response{Code: http.StatusBadRequest, Msg: fmt.Errorf("item id err:%w", err).Error()}, err
	}

	//删除逻辑
	err = c.service.Delete(g, uint(itemId), id)
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
func (c *RemoveController) Get(g *gin.Context) (response.Response, error) {
	idsStr := g.Query("ids")          // 获取字符串 "1,2,3"
	ids := strings.Split(idsStr, ",") // 手动切割成字符串数组
	id, err := c.CheckPower(g)
	if err != nil {
		return response.Response{Code: http.StatusBadRequest, Msg: fmt.Errorf("power check err:%w", err).Error()}, err
	}
	ItemIds, err := stringSliceToUintSlice(ids)
	if err != nil {
		return response.Response{}, err
	}
	items, err := c.service.Get(g, ItemIds, id)
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
func (c *RemoveController) CheckPower(g *gin.Context) (uint, error) {
	ac := g.GetHeader("AccessKey")
	temp := g.GetHeader("Timestamp")
	signature := g.GetHeader("Signature")
	if ac == "" || temp == "" || signature == "" {

		return 0, errors.New("http header parameters required")
	}
	//鉴权,返回project_id
	ok, id, err := c.service.CheckPower(g, ac, temp, signature)
	if !ok || err != nil {
		return 0, err
	}
	return id, nil
}
