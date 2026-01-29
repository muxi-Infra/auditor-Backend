// @securityDefinitions.keyget ApiKeyAuth
// @in header
// @name Authorization

package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"strconv"

	"github.com/muxi-Infra/auditor-Backend/api/request"
	"github.com/muxi-Infra/auditor-Backend/api/response"
	"github.com/muxi-Infra/auditor-Backend/pkg/jwt"
	"github.com/muxi-Infra/auditor-Backend/repository/model"
	"github.com/muxi-Infra/auditor-Backend/service"
)

type ItemController struct {
	service ItemService
}
type ItemService interface {
	Select(ctx context.Context, req request.SelectReq) ([]model.Item, error)
	Audit(g context.Context, req request.AuditReq, id uint) (request.WebHookData, model.Item, error)
	SearchHistory(g context.Context, id uint) ([]model.Item, error)
	Upload(g context.Context, req request.UploadReq, key string) (uint, error)
	Hook(request.WebHookData, model.Item) error
	RoleBack(item model.Item)
	GetDetail(ctx context.Context, id uint) (model.Item, error)
	AuditMany(g context.Context, items []request.AuditReq, uid uint) []request.WebHookData
}

func NewItemController(service *service.ItemService) *ItemController {
	return &ItemController{
		service: service,
	}
}

// Select 集成查询item
// @Summary 获取条目列表
// @Description 根据请求的条件获取项目和相关项目信息
// @Tags Item
// @Accept json
// @Produce json
// @Param selectReq body request.SelectReq true "查询条件"
// @Success 200 {object} response.Response{data=[]response.Item} "成功返回项目列表"
// @Failure 400 {object} response.Response "查询失败"
// @Router /api/v1/item/select [post]
func (ic *ItemController) Select(c *gin.Context, req request.SelectReq) (response.Response, error) {
	if req.ProjectID == 0 {
		return response.Response{
			Data: nil,
			Msg:  "需要project_id",
			Code: 400,
		}, nil
	}
	it, err := ic.service.Select(c, req)
	if err != nil {
		return response.Response{
			Data: nil,
			Code: 400,
			Msg:  "搜索失败",
		}, err
	}
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
			Id:         item.ID,
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
		Data: items,
		Code: 200,
	}, nil
}

// Audit 审核item
// @Summary 审核条目
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
		er := ic.service.Hook(data, item)

		if er != nil {
			ic.service.RoleBack(item)
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
			Id:         item.ID,
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
// @Summary 上传条目,这个似乎不用接
// @Description 上传新的项目或文件
// @Tags Item
// @Accept json
// @Produce json
// @Param uploadReq body request.UploadReq true "上传请求体"
// @Param api_key header string true "API 认证密钥(api_key)"
// @Success 200 {object} response.Response "上传成功"
// @Failure 400 {object} response.Response "上传失败"
// @Security ApiKeyAuth
// @Router /api/v1/item/upload [put]
func (ic *ItemController) Upload(g *gin.Context, req request.UploadReq) (response.Response, error) {
	key := g.GetHeader("api_key")
	id, err := ic.service.Upload(g, req, key)
	if err != nil {
		return response.Response{
			Msg:  "上传失败",
			Code: 400,
			Data: nil,
		}, err
	}
	return response.Response{
		Msg:  "success",
		Data: id,
		Code: 200,
	}, nil
}

// Detail 获取单个条目信息
// @Summary 通过id获取条目具体信息
// @Description 通过id获取条目具体信息
// @Tags Item
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=response.Item} "成功返回条目"
// @Failure 400 {object} response.Response "获取条目失败"
// @Security ApiKeyAuth
// @Router /api/v1/item/{item_id}/detail [get]
func (ic *ItemController) Detail(ctx *gin.Context) (response.Response, error) {
	ItemID := ctx.Param("item_id")
	if ItemID == "" {
		return response.Response{
			Code: 400,
			Msg:  "需要item_id",
		}, nil
	}
	itemId, err := strconv.ParseUint(ItemID, 10, 64)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  "获取item_id失败",
		}, err
	}
	id := uint(itemId)
	item, err := ic.service.GetDetail(ctx, id)
	if err != nil {
		return response.Response{
			Code: 400,
			Msg:  "获取条目失败",
		}, err
	}
	if len(item.Comments) == 0 {
		re := response.Item{
			Id:         item.ID,
			Author:     item.Author,
			Tags:       item.Tags,
			Status:     item.Status,
			PublicTime: item.CreatedAt.UnixMilli(),
			Content: response.Contents{
				Topic: response.Topics{
					Title:    item.Title,
					Content:  item.Content,
					Pictures: item.Pictures,
				},
			},
		}
		return response.Response{
			Code: 200,
			Msg:  "获取条目成功",
			Data: re,
		}, nil
	} else if len(item.Comments) == 1 {
		re := response.Item{
			Id:         item.ID,
			Author:     item.Author,
			Tags:       item.Tags,
			Status:     item.Status,
			PublicTime: item.CreatedAt.UnixMilli(),
			Content: response.Contents{
				Topic: response.Topics{
					Title:    item.Title,
					Content:  item.Content,
					Pictures: item.Pictures,
				},
				LastComment: response.Comment{
					Content:  item.Comments[0].Content,
					Pictures: item.Comments[0].Pictures,
				},
			},
		}
		return response.Response{
			Code: 200,
			Msg:  "获取条目成功",
			Data: re,
		}, nil
	} else {
		re := response.Item{
			Id:         item.ID,
			Author:     item.Author,
			Tags:       item.Tags,
			Status:     item.Status,
			PublicTime: item.CreatedAt.UnixMilli(),
			Content: response.Contents{
				Topic: response.Topics{
					Title:    item.Title,
					Content:  item.Content,
					Pictures: item.Pictures,
				},
				LastComment: response.Comment{
					Content:  item.Comments[0].Content,
					Pictures: item.Comments[0].Pictures,
				},
				NextComment: response.Comment{
					Content:  item.Comments[1].Content,
					Pictures: item.Comments[1].Pictures,
				},
			},
		}
		return response.Response{
			Code: 200,
			Msg:  "获取条目成功",
			Data: re,
		}, nil
	}

}

// AuditAll 批量审核item
// @Summary 批量审核条目，可接受拒绝交杂
// @Description 批量审核项目并更新审核状态,不要超过10个
// @Tags Item
// @Accept json
// @Produce json
// @Param auditReq body request.ManyAuditReq true "审核请求体"
// @Success 200 {object} response.Response "批量审核成功"
// @Failure 400 {object} response.Response "批量审核失败"
// @Failure 400 {object} response.Response "too many items"
// @Security ApiKeyAuth
// @Router /api/v1/item/auditMany [post]
func (ic *ItemController) AuditAll(ctx *gin.Context, req request.ManyAuditReq, cla jwt.UserClaims) (response.Response, error) {
	if cla.UserRule == 0 {
		return response.Response{
			Msg:  "no power",
			Code: 403,
			Data: nil,
		}, nil
	}
	if len(req.Reqs) > 10 {
		return response.Response{
			Msg:  "too many items",
			Code: 400,
			Data: nil,
		}, nil
	}
	datas := ic.service.AuditMany(ctx, req.Reqs, cla.Uid)
	if len(datas) == 0 {

		return response.Response{
			Msg:  "批量审核成功",
			Code: 200,
			Data: nil,
		}, nil
	} else {
		return response.Response{
			Msg:  "批量审核失败案例",
			Code: 400,
			Data: datas,
		}, nil
	}
}
