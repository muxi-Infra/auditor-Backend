package request

import "github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"

type AuditItem struct {
	ID        uint
	ProjectID uint
	Contents  response.Contents
}

type AuditByLLMReq struct {
	Data []AuditItem
}
