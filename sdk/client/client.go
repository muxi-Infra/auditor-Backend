package client

import (
	"net/http"
)

type MuxiAuditClient struct {
	c *http.Client
}

func NewMuxiAuditClient(c *http.Client) *MuxiAuditClient {
	return &MuxiAuditClient{
		c: c,
	}
}

func (mc *MuxiAuditClient) UploadItem(path string, data []byte) ([]byte, error) {
	return nil, nil
}
func (mc *MuxiAuditClient) UpdateItem(path string) ([]byte, error) {
	return nil, nil
}
func (mc *MuxiAuditClient) DeleteItem(path string) error        {}
func (mc *MuxiAuditClient) GetItem(path string) ([]byte, error) {}
