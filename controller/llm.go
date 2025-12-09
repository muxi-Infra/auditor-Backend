package controller

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"

	api_errors "github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/errors"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/jwt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/service"
)

const (
	MaxTaskNumber = 20
)

type LLMController struct {
	service LLMService
}

type LLMService interface {
	Audit(Data []request.AuditItem)
	Close()
}

func NewLLMController(service *service.LLMService) *LLMController {
	return &LLMController{
		service: service,
	}
}

// Audit
// @Summary ai审核条目
// @Description 根据请求将需要审核的条目加入ai审核队列
// @Tags LLM
// @Accept json
// @Produce json
// @Param auditReq body request.AuditByLLMReq true "审核请求"
// @Success 200 {object} response.Response{} "成功返回success"
// @Failure 400 {object} response.Response "审核失败"
// @Router /api/v1/llm/audit [post]
func (c *LLMController) Audit(ctx *gin.Context, req request.AuditByLLMReq, cla jwt.UserClaims) (response.Response, error) {
	if !cla.IfStaff() {
		return response.Response{Data: nil, Code: 40003, Msg: "no power"},
			api_errors.PERMISSION_DENIED_ERROR(errors.New("user is not our person "))
	}
	if len(req.Data) > MaxTaskNumber {
		return response.Response{Data: nil, Code: 40000, Msg: "task is too many"},
			api_errors.BAD_REQUEST_ERROR(errors.New("data is too many"))
	}
	c.service.Audit(req.Data)
	return response.Response{Msg: "success", Code: 0}, nil
}

func (c *LLMController) Close() {
	fmt.Println("LLMController Close")
	c.service.Close()
}
