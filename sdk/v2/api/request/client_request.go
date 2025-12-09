package request

import (
	"fmt"
	"github.com/muxi-Infra/auditor-Backend/sdk/v2/internal"
)

// UploadReq
// 这里重新定义了一遍req而并没用使用审核外部定义的结构体，
// 是为了不暴露内部实现让用户无感知.
type UploadReq struct {
	HookUrl    *string            `json:"hook_url,omitempty"`    // 必填
	Id         *uint              `json:"id,omitempty"`          // 必填，hook_id
	Author     *string            `json:"author,omitempty"`      // 可选
	PublicTime *int64             `json:"public_time,omitempty"` // 可选
	Tags       *[]string          `json:"tags,omitempty"`        // 可选
	Content    *internal.Contents `json:"content,omitempty"`     // 必填
	Extra      interface{}        `json:"extra,omitempty"`       // 可选扩展字段
}

type UploadOpt func(*UploadReq)
type UpdateOpt func(*UpdateReq)

func (req *UploadReq) IsValid() bool {
	if req.HookUrl == nil || req.Id == nil || req.Content == nil {
		return false
	}
	if *req.HookUrl == "" {
		return false
	}
	if *req.Id <= 0 {
		return false
	}
	return true
}

// NewUploadReq 你应当始终使用此函数来创建对象
func NewUploadReq(hookUrl string, id uint, contents *internal.Contents, opts ...UploadOpt) (*UploadReq, error) {
	req := &UploadReq{
		HookUrl: &hookUrl,
		Id:      &id,
		Content: contents,
	}
	if !req.IsValid() {
		return nil, fmt.Errorf("illegal UploadReq params: hookUrl and id must be non-empty, contents must not be nil, id must > 0")
	}

	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

func WithUploadAuthor(author string) UploadOpt {
	return func(req *UploadReq) {
		if author == "" {
			req.Author = nil
		} else {
			req.Author = &author
		}
	}
}

func WithUploadPublicTime(ts int64) UploadOpt {
	return func(req *UploadReq) {
		if ts <= 0 {
			req.PublicTime = nil
		} else {
			req.PublicTime = &ts
		}
	}
}

func WithUploadTags(tags []string) UploadOpt {
	return func(req *UploadReq) {
		if len(tags) == 0 {
			req.Tags = nil
		} else {
			req.Tags = &tags
		}
	}
}

func WithUploadExtra(extra interface{}) UploadOpt {
	return func(req *UploadReq) {
		req.Extra = extra
	}
}

type UpdateReq struct {
	HookUrl    *string            `json:"hook_url,omitempty"`    // 可选
	Id         *uint              `json:"id,omitempty"`          // 必填，hook_id
	Author     *string            `json:"author,omitempty"`      // 可选
	PublicTime *int64             `json:"public_time,omitempty"` // 可选
	Tags       *[]string          `json:"tags,omitempty"`        // 可选
	Content    *internal.Contents `json:"content,omitempty"`     // 可选
	Extra      interface{}        `json:"extra,omitempty"`       // 可选扩展字段
}

func NewUpdateReq(id uint, opts ...UpdateOpt) (*UpdateReq, error) {
	req := &UpdateReq{
		Id: &id,
	}
	if !req.IsValid() {
		return nil, fmt.Errorf("illegal UploadReq params:id must be non-empty, id must > 0")
	}

	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

func WithUpdateContent(content *internal.Contents) UpdateOpt {
	return func(req *UpdateReq) {
		if content == nil {
			return
		}
		req.Content = content
	}
}

func WithUpdateUrl(url string) UpdateOpt {
	return func(req *UpdateReq) {
		if url == "" {
			return
		}
		req.HookUrl = &url
	}
}

func WithUpdateAuthor(author string) UpdateOpt {
	return func(req *UpdateReq) {
		if author == "" {
			req.Author = nil
		} else {
			req.Author = &author
		}
	}
}

func WithUpdatePublicTime(ts int64) UpdateOpt {
	return func(req *UpdateReq) {
		if ts <= 0 {
			req.PublicTime = nil
		} else {
			req.PublicTime = &ts
		}
	}
}

func WithUpdateTags(tags []string) UpdateOpt {
	return func(req *UpdateReq) {
		if len(tags) == 0 {
			req.Tags = nil
		} else {
			req.Tags = &tags
		}
	}
}

func WithUpdateExtra(extra interface{}) UpdateOpt {
	return func(req *UpdateReq) {
		req.Extra = extra
	}
}

func (req *UpdateReq) IsValid() bool {
	if req.Id != nil && *req.Id > 0 {
		return true
	}
	return false
}

type DeleteReq struct {
	Id *uint `json:"id,omitempty"`
}

func NewDeleteReq(id uint) (*DeleteReq, error) {
	req := &DeleteReq{
		Id: &id,
	}

	if !req.IsValid() {
		return nil, fmt.Errorf("illegal DeleteReq params:id must be non-empty, id must > 0")
	}

	return req, nil
}

func (req *DeleteReq) IsValid() bool {
	if req.Id != nil && *req.Id > 0 {
		return true
	}
	return false
}

type GetItemsStatusReq struct {
	Ids *[]int `json:"ids,omitempty"` // 只传一个0代表获取该项目下的所有审核条目结果
}

func NewGetItemsStatusReq(ids []int) (*GetItemsStatusReq, error) {
	req := &GetItemsStatusReq{
		Ids: &ids,
	}

	if !req.IsValid() {
		return nil, fmt.Errorf("illegal Req params:ids must be non-empty, len must > 0")
	}

	return req, nil
}

func (req *GetItemsStatusReq) IsValid() bool {
	if req.Ids != nil && len(*req.Ids) > 0 {
		return true
	}
	return false
}
