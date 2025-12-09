package service

import (
	"github.com/muxi-Infra/auditor-Backend/langchain/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuditStatusToString(t *testing.T) {
	re := auditStatusToString(model.Pass)
	assert.Equal(t, re, "Pass")
}
