package request

import "github.com/muxi-Infra/auditor-Backend/api/response"

type AuditItem struct {
	ID        uint
	ProjectID uint
	Contents  response.Contents
}

type AuditByLLMReq struct {
	Data []AuditItem
}
