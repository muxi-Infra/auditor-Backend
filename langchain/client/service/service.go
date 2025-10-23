package service

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/strategy"
)

type ServiceInterface interface {
	RegisterService(r interface{})
	RegisterStrategy(s strategy.AIAuditStrategy)
}
