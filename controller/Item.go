// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"muxi_auditor/api/request"
	"muxi_auditor/api/response"
	"muxi_auditor/pkg/jwt"
	"muxi_auditor/repository/model"
	"muxi_auditor/service"
)

type ItemController struct {
	service ItemService
}
type ItemService interface {
	Select(ctx context.Context, req request.SelectReq) ([]model.Item, error)
	Audit(g context.Context, req request.AuditReq, id uint) (service.Data, model.Item, error)
	SearchHistory(g context.Context, id uint) ([]model.Item, error)
	Upload(g context.Context, req request.UploadReq, key string) error
	Hook(service.Data, model.Item) error
	RoleBack(item model.Item) error
}

func NewItemController(service *service.ItemService) *ItemController {
	return &ItemController{
		service: service,
	}
}

// Select 集成查询item
// @Summary 获取项目列表
// @Description 根据请求的条件获取项目和相关项目信息
// @Tags Item
// @Accept json
// @Produce json
// @Param selectReq body request.SelectReq true "查询条件"
// @Success 200 {object} response.Response{data=[]response.SelectResp} "成功返回项目列表"
// @Failure 400 {object} response.Response "查询失败"
// @Router /api/v1/item/select [post]
func (ic *ItemController) Select(c *gin.Context, req request.SelectReq) (response.Response, error) {

	it, err := ic.service.Select(c, req)
	if err != nil {
		return response.Response{
			Data: nil,
			Code: 400,
			Msg:  "搜索失败",
		}, err
	}
	var re []response.SelectResp
	var items []response.Item

	for _, item := range it {
		lastComment := response.Comment{}
		nextComment := response.Comment{}
		unixTimestamp := item.CreatedAt.UnixMilli()
		if len(item.Comments) > 0 {
			lastComment = response.Comment{
				Content:  item.Comments[0].Content,
				Pictures: item.Comments[0].Pictures,
			}
		}
		if len(item.Comments) > 1 {
			nextComment = response.Comment{
				Content:  item.Comments[1].Content,
				Pictures: item.Comments[1].Pictures,
			}
		}

		items = append(items, response.Item{
			ItemId:     item.ID,
			Author:     item.Author,
			Tags:       item.Tags,
			Status:     item.Status,
			PublicTime: unixTimestamp,
			Auditor:    item.Auditor,
			Content: response.Contents{
				Topic: response.Topics{
					Title:    item.Title,
					Content:  item.Content,
					Pictures: item.Pictures,
				},
				LastComment: lastComment,
				NextComment: nextComment,
			},
		})
	}
	re = append(re, response.SelectResp{
		Items: items,
	})

	return response.Response{
		Msg:  "success",
		Data: re,
		Code: 200,
	}, nil
}

// Audit 审核item
// @Summary 审核项目
// @Description 审核项目并更新审核状态
// @Tags Item
// @Accept json
// @Produce json
// @Param auditReq body request.AuditReq true "审核请求体"
// @Success 200 {object} response.Response "审核成功"
// @Failure 400 {object} response.Response "审核失败"
// @Security ApiKeyAuth
// @Router /api/v1/item/audit [post]
func (ic *ItemController) Audit(c *gin.Context, req request.AuditReq, cla jwt.UserClaims) (response.Response, error) {
	data, item, err := ic.service.Audit(c, req, cla.Uid)
	if err != nil {
		return response.Response{
			Msg:  "提交失败",
			Code: 400,
			Data: nil,
		}, err
	}

	go func() {
		err = ic.service.Hook(data, item)
		if err != nil {
			log.Println(err)
			err = ic.service.RoleBack(item)
			if err != nil {
				log.Println(err)
			}
		}
	}()
	return response.Response{
		Msg:  "success",
		Data: nil,
		Code: 200,
	}, nil

}

// SearchHistory 获取个人审核历史记录
// @Summary 获取历史记录
// @Description 获取用户的历史记录（审核历史）
// @Tags Item
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]response.Item} "成功返回历史记录"
// @Failure 400 {object} response.Response "获取历史记录失败"
// @Security ApiKeyAuth
// @Router /api/v1/item/searchHistory [get]
func (ic *ItemController) SearchHistory(g *gin.Context, cla jwt.UserClaims) (response.Response, error) {

	items, err := ic.service.SearchHistory(g, cla.Uid)
	if err != nil {
		return response.Response{
			Msg:  "获取历史记录失败",
			Code: 400,
			Data: nil,
		}, err
	}
	var it []response.Item
	for _, item := range items {
		lastComment := response.Comment{}
		nextComment := response.Comment{}
		unixTimestamp := item.CreatedAt.UnixMilli()
		if len(item.Comments) > 0 {
			lastComment = response.Comment{
				Content:  item.Comments[0].Content,
				Pictures: item.Comments[0].Pictures,
			}
		}
		if len(item.Comments) > 1 {
			nextComment = response.Comment{
				Content:  item.Comments[1].Content,
				Pictures: item.Comments[1].Pictures,
			}
		}

		it = append(it, response.Item{
			ItemId:     item.ID,
			Author:     item.Author,
			Tags:       item.Tags,
			Status:     item.Status,
			PublicTime: unixTimestamp,
			Auditor:    item.Auditor,
			Content: response.Contents{
				Topic: response.Topics{
					Title:    item.Title,
					Content:  item.Content,
					Pictures: item.Pictures,
				},
				LastComment: lastComment,
				NextComment: nextComment,
			},
		})
	}
	return response.Response{
		Msg:  "success",
		Data: it,
		Code: 200,
	}, nil

}

// Upload 上传item
// @Summary 上传项目
// @Description 上传新的项目或文件
// @Tags Item
// @Accept json
// @Produce json
// @Param uploadReq body request.UploadReq true "上传请求体"
// @Param api_key header string true "API 认证密钥(api_key)"
// @Success 200 {object} response.Response "上传成功"
// @Failure 400 {object} response.Response "上传失败"
// @Security ApiKeyAuth
// @Router /api/v1/item/upload [post]
func (ic *ItemController) Upload(g *gin.Context, req request.UploadReq, cla jwt.UserClaims) (response.Response, error) {
	key := g.GetHeader("api_key")
	err := ic.service.Upload(g, req, key)
	if err != nil {
		return response.Response{
			Msg:  "上传失败",
			Code: 400,
			Data: nil,
		}, err
	}
	return response.Response{
		Msg:  "success",
		Data: nil,
		Code: 200,
	}, nil
}
