package service

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuditStatusToString(t *testing.T) {
	re := auditStatusToString(model.Pass)
	assert.Equal(t, re, "Pass")
}
