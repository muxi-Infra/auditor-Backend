package controller

import (
	"context"

	"github.com/gin-gonic/gin"
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
	Select(ctx context.Context, req request.SelectReq) ([]model.Project, error)
	Audit(g context.Context, req request.AuditReq, id uint) error
	SearchHistory(g context.Context, id uint) ([]model.Item, error)
	Upload(g context.Context, req request.UploadReq) error
}

func NewItemController(service *service.ItemService) *ItemController {
	return &ItemController{
		service: service,
	}
}
func (ic *ItemController) Select(c *gin.Context, cla jwt.UserClaims, req request.SelectReq) (response.Response, error) {
	projects, err := ic.service.Select(c, req)
	if err != nil {
		return response.Response{
			Data: nil,
			Code: 400,
			Msg:  "搜索失败",
		}, err
	}
	var re []response.SelectResp
	var items []response.Item
	for _, project := range projects {
		for _, item := range project.Items {
			lastComment := response.Comment{}
			nextComment := response.Comment{}
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
				PublicTime: item.CreatedAt,
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
			Items:     items,
			ProjectId: project.ID,
		})
	}
	return response.Response{
		Msg:  "success",
		Data: re,
		Code: 200,
	}, nil
}
func (ic *ItemController) Audit(c *gin.Context, cla jwt.UserClaims, req request.AuditReq) (response.Response, error) {
	err := ic.service.Audit(c, req, cla.Uid)
	if err != nil {
		return response.Response{
			Msg:  "提交失败",
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
			PublicTime: item.CreatedAt,
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
func (ic *ItemController) Upload(g *gin.Context, cla jwt.UserClaims, req request.UploadReq) (response.Response, error) {
	err := ic.service.Upload(g, req)
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
