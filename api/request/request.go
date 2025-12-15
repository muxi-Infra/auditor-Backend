package request

import (
	"github.com/muxi-Infra/auditor-Backend/api/response"
	"github.com/muxi-Infra/auditor-Backend/repository/model"
)

type LoginReq struct {
	Code string `json:"code"`
}
type RegisterReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type UpdateUserReq struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
type GetUserReq struct {
	Project_id uint `json:"project_id"`
}
type UpdateUserRoleReq struct {
	Role          int                   `json:"role"` //用户权限
	UserId        uint                  `json:"user_id"`
	ProjectPermit []model.ProjectPermit `json:"project_permit"` //允许的项目列表
}
type CreateProject struct {
	Name        string          `json:"name"`
	Logo        string          `json:"logo"`
	AuditRule   string          `json:"audit_rule"` //审核规则
	Users       []UserInProject `json:"users"`
	HookUrl     string          `json:"hook_url"`
	Description string          `json:"description"`
}
type GetProjectDetail struct {
	ProjectId uint `json:"project_id"`
}
type SelectReq struct {
	ProjectID uint               `json:"project_id"`
	RoundTime [][]int            `json:"round_time"` //日期
	Tags      []string           `json:"tags"`       //标签
	Statuses  []model.ItemStatus `json:"statuses"`
	Auditors  []uint             `json:"auditors"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
	Query     string             `json:"query"` //查询字段
}
type AuditReq struct {
	Reason string           `json:"reason"`
	Status model.ItemStatus `json:"status"` //0未审核，1通过，2未通过
	ItemId uint             `json:"item_id"`
}
type UploadReq struct {
	HookUrl    string            `json:"hook_url"`
	Id         uint              `json:"id"` //hook_id
	Author     string            `json:"author"`
	PublicTime int64             `json:"public_time"`
	Tags       []string          `json:"tags"`
	Content    response.Contents `json:"content"`
	Extra      interface{}       `json:"extra"`
}
type DeleteProject struct {
	ProjectId uint `json:"project_id"`
}
type UpdateProject struct {
	ProjectName string `json:"project_name"`
	Logo        string `json:"logo"`
	AuditRule   string `json:"audit_rule"`
	Description string `json:"description"`
}

type GetUsers struct {
	Query    string `json:"query"`
	Page     int    `json:"page"`
	PageSize int    `json:"size"`
}
type HookPayload struct {
	Event string      `json:"event"`
	Data  WebHookData `json:"data,omitempty"`
	Try   int         `json:"try"` // 重试次数，最大不超过五
}

type ReturnApiKey struct {
	ApiKey  string `json:"api_key"`
	Message string `json:"message"`
}
type ManyAuditReq struct {
	Reqs []AuditReq
}
type UserInProject struct {
	Userid      uint `json:"user_id"`
	ProjectRole int  `json:"project_role"`
}
type UserRole struct {
	Userid uint `json:"user_id"`
	Role   int  `json:"role"` //审核平台的权限，并非项目中的权限,0无权限，1普通用户，2管理者
}
type ChangeUserRolesReq struct {
	List []UserRole `json:"list"`
}
type AddUser struct {
	UserId      uint `json:"user_id"`
	ProjectRole int  `json:"project_role"`
}

type AddUsersReq struct {
	AddUsers []AddUser `json:"add_users"`
}

type DeleteUsers struct {
	Ids []uint `json:"ids"`
}

type WebHookData struct {
	Id     uint
	Status string
	Msg    string
}
