package client

import (
	"context"
	"errors"

	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/api/errorx"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/client/base"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/config"
)

type Config = config.Config

// Client 后期可基于此拓展其他功能
type Client struct {
	base.Client
}

func NewClient(config Config) (*Client, error) {
	c := new(Client)
	err := c.Client.Init(config)
	if err != nil {
		return nil, err
	}
	return c, nil
}

/*
*
	接口设计理念: 错误职责明确区分，开箱即用，参数先验，异常信息保留返回。

	1. 永远先行检查用户请求数据是否正确，否则可能造成审核服务端脏数据堆积和ai服务token的浪费。

	2. 要有明确的错误边界，能让用户感知到是最基础性的错误还是使用sdk时发生的错误还是服务端那边产生的错误，这里
	目前的处理是对于用户传参解析之类的基础错误直接以err形式返回，让用户在错误检查时就可以知道是自己使用不当，至
	于服务端的bad 情求错误和后续造成的数据格式转换错误统一放在response的errorx字段，让用户明白自己并非使用问
	题引发的错误而是服务端处的问题。

	3. 自定义了SDKErr用于帮助最后详细错误信息的展示和在真正逻辑处区分服务端调用错误和sdk处的错误。
*
*/

// UploadItem 用于向已创建项目上传审核内容
func (c *Client) UploadItem(ctx context.Context, req *request.UploadReq) (response.UploadItemResp, error) {
	if ok := req.IsValid(); !ok {
		return response.UploadItemResp{}, errors.New("invalid request")
	}

	row, err := c.Client.UploadItem(ctx, *req)
	if err != nil {
		if errors.As(err, &errorx.DefaultErr) {
			return response.UploadItemResp{}, err
		}

		return response.UploadItemResp{
			Basic: response.Basic{
				Code:   errorx.RequestErrCode,
				Msg:    "请求失败",
				Errorx: errorx.New(row.Code, errorx.RequestErrCode, row.Msg, err),
			},
			ItemID: 0,
		}, nil
	}

	da, err := extractInt(row)
	if err != nil {
		return response.UploadItemResp{
			Basic: response.Basic{
				Code:   errorx.SeverDataIllegalCode,
				Msg:    "响应数据data格式错误",
				Errorx: err,
			},

			ItemID: 0,
		}, nil
	}

	return response.UploadItemResp{
		Basic: response.Basic{
			Code:   errorx.SuccessCode,
			Msg:    "success",
			Errorx: nil,
		},
		ItemID: da,
	}, nil
}

// UpdateItem 更新已创建的审核内容
func (c *Client) UpdateItem(ctx context.Context, req *request.UpdateReq) (response.UpdateItemResp, error) {
	if ok := req.IsValid(); !ok {
		return response.UpdateItemResp{}, errors.New("invalid request")
	}

	row, err := c.Client.UpdateItem(ctx, *req)
	if err != nil {
		if errors.As(err, &errorx.DefaultErr) {
			return response.UpdateItemResp{}, err
		}

		return response.UpdateItemResp{
			Basic: response.Basic{
				Code:   errorx.RequestErrCode,
				Msg:    "请求失败",
				Errorx: errorx.New(row.Code, errorx.RequestErrCode, row.Msg, err),
			},
			ItemID: 0,
		}, nil
	}

	da, err := extractInt(row)
	if err != nil {
		return response.UpdateItemResp{
			Basic: response.Basic{
				Code:   errorx.SeverDataIllegalCode,
				Msg:    "响应数据data格式错误",
				Errorx: err,
			},
			ItemID: 0,
		}, nil
	}

	return response.UpdateItemResp{
		Basic: response.Basic{
			Code:   errorx.SuccessCode,
			Msg:    "success",
			Errorx: nil,
		},
		ItemID: da,
	}, nil
}

// DeleteItem 删除已创建的审核内容。
func (c *Client) DeleteItem(ctx context.Context, req *request.DeleteReq) (response.DeleteItemResp, error) {
	if ok := req.IsValid(); !ok {
		return response.DeleteItemResp{}, errors.New("invalid request")
	}

	row, err := c.Client.DeleteItem(ctx, int(*req.Id))
	if err != nil {
		if errors.As(err, &errorx.DefaultErr) {
			return response.DeleteItemResp{}, err
		}

		return response.DeleteItemResp{
			Basic: response.Basic{
				Code:   errorx.RequestErrCode,
				Msg:    "请求失败",
				Errorx: errorx.New(row.Code, errorx.RequestErrCode, row.Msg, err),
			},
			ItemID: 0,
		}, nil
	}

	da, err := extractInt(row)
	if err != nil {
		return response.DeleteItemResp{
			Basic: response.Basic{
				Code:   errorx.SeverDataIllegalCode,
				Msg:    "响应数据data格式错误",
				Errorx: err,
			},
			ItemID: 0,
		}, nil
	}

	return response.DeleteItemResp{
		Basic: response.Basic{
			Code:   errorx.SuccessCode,
			Msg:    "success",
			Errorx: nil,
		},
		ItemID: da,
	}, nil
}

func (c *Client) GetItems(ctx context.Context, req *request.GetItemsStatusReq) (response.GetItemsResp, error) {
	if ok := req.IsValid(); !ok {
		return response.GetItemsResp{}, errors.New("invalid request")
	}

	row, err := c.Client.GetItems(ctx, *req.Ids)
	if err != nil {
		if errors.As(err, &errorx.DefaultErr) {
			return response.GetItemsResp{}, err
		}

		return response.GetItemsResp{
			Basic: response.Basic{
				Code:   errorx.RequestErrCode,
				Msg:    "请求失败",
				Errorx: errorx.New(row.Code, errorx.RequestErrCode, row.Msg, err),
			},
			Items: nil,
		}, nil
	}

	da, err := extractItemsStatus(row)
	if err != nil {
		return response.GetItemsResp{
			Basic: response.Basic{
				Code:   errorx.SeverDataIllegalCode,
				Msg:    "响应数据data格式错误",
				Errorx: err,
			},
			Items: nil,
		}, nil
	}

	re := make([]response.ItemsStatus, 0, len(da.Items))
	for _, v := range da.Items {
		re = append(re, response.ItemsStatus{
			Id:     v.HookId,
			Status: TransformStatusToString(v.Status),
		})
	}

	return response.GetItemsResp{
		Basic: response.Basic{
			Code:   errorx.SuccessCode,
			Msg:    "success",
			Errorx: nil,
		},
		Items: re,
	}, nil
}
