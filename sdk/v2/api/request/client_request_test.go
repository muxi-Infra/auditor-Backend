package request

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/muxi-Infra/auditor-Backend/sdk/v2/dto"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewAuditReq(t *testing.T) {
	// Define valid parameters
	hookUrl := "http://example.com/webhook"
	id := uint(1)
	contents := dto.NewContents(dto.WithTopicText("11", "11"))

	// Test creating a valid UploadReq
	req, err := NewUploadReq(hookUrl, id, contents)

	assert.NoError(t, err, "expected no error when creating valid UploadReq")
	assert.NotNil(t, req, "expected UploadReq to be non-nil")
	assert.Equal(t, hookUrl, *req.HookUrl, "expected HookUrl to match")
	assert.Equal(t, id, *req.Id, "expected Id to match")
	assert.Equal(t, contents, req.Content, "expected Content to match")
}

// Test NewUploadReq with invalid parameters
func TestNewAuditReq_InvalidParams(t *testing.T) {
	// Invalid parameters: empty hookUrl, zero id, nil contents
	req, err := NewUploadReq("", 0, nil)

	assert.Error(t, err, "expected error for invalid UploadReq parameters")
	assert.Nil(t, req, "expected UploadReq to be nil for invalid params")
}

// Test IsValid function
func TestIsValid(t *testing.T) {
	// Test for valid case
	hookUrl := "http://example.com/webhook"
	id := uint(1)
	contents := dto.NewContents(dto.WithTopicText("11", "11"))
	req := &UploadReq{
		HookUrl: &hookUrl,
		Id:      &id,
		Content: contents,
	}

	assert.True(t, req.IsValid(), "expected IsValid to return true for valid UploadReq")

	// Test for invalid cases
	tests := []struct {
		name     string
		req      *UploadReq
		expected bool
	}{
		{
			name:     "missing HookUrl",
			req:      &UploadReq{Id: &id, Content: contents},
			expected: false,
		},
		{
			name:     "empty HookUrl",
			req:      &UploadReq{HookUrl: tea.String(""), Id: &id, Content: contents},
			expected: false,
		},
		{
			name:     "missing Id",
			req:      &UploadReq{HookUrl: &hookUrl, Content: contents},
			expected: false,
		},
		{
			name:     "invalid Id",
			req:      &UploadReq{HookUrl: &hookUrl, Id: new(uint), Content: contents},
			expected: false,
		},
		{
			name:     "missing Content",
			req:      &UploadReq{HookUrl: &hookUrl, Id: &id},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.req.IsValid())
		})
	}
}

// Test optional fields using Options
func TestWithOptions(t *testing.T) {
	hookUrl := "http://example.com/webhook"
	id := uint(1)
	contents := dto.NewContents(dto.WithTopicText("11", "11"))

	// Test adding optional fields using options
	req, err := NewUploadReq(hookUrl, id, contents,
		WithUploadAuthor("John Doe"),
		WithUploadPublicTime(time.Now().Unix()), WithUploadTags([]string{"tag1", "tag2"}))

	assert.NoError(t, err, "expected no error when using options")
	assert.NotNil(t, req, "expected UploadReq to be non-nil")
	assert.Equal(t, "John Doe", *req.Author, "expected Author to be set")
	assert.NotNil(t, req.PublicTime, "expected PublicTime to be set")
	assert.Equal(t, []string{"tag1", "tag2"}, *req.Tags, "expected Tags to be set")
}

// Test WithUploadExtra option
func TestWithExtra(t *testing.T) {
	hookUrl := "http://example.com/webhook"
	id := uint(1)
	contents := dto.NewContents(dto.WithTopicText("11", "11"))

	extraData := map[string]interface{}{"key1": "value1", "key2": 2}
	req, err := NewUploadReq(hookUrl, id, contents, WithUploadExtra(extraData))

	assert.NoError(t, err, "expected no error when adding extra data")
	assert.NotNil(t, req, "expected UploadReq to be non-nil")
	assert.Equal(t, extraData, req.Extra, "expected Extra to match")
}

// Test WithUploadAuthor setting nil for empty author
func TestWithAuthorEmpty(t *testing.T) {
	hookUrl := "http://example.com/webhook"
	id := uint(1)
	contents := dto.NewContents(dto.WithTopicText("11", "11"))

	// Test empty author string, which should set Author to nil
	req, err := NewUploadReq(hookUrl, id, contents, WithUploadAuthor(""))

	assert.NoError(t, err, "expected no error")
	assert.Nil(t, req.Author, "expected Author to be nil when empty string passed")
}
