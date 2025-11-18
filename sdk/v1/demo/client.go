package main

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v1/client"
	"net/http"
	"time"
)

func main() {
	ac := "*****************"
	c := client.NewMuxiAuditClient(&http.Client{}, "http://0.0.0.0:8080/api/v1")
	var req = request.UploadReq{
		HookUrl:    "http://localhost:8081/audit",
		Id:         12,
		Author:     "chen",
		PublicTime: time.Now().Unix(),
		Tags:       make([]string, 0),
		Content: response.Contents{
			Topic: response.Topics{
				Title:    "test2",
				Content:  "test2",
				Pictures: nil,
			},
			LastComment: response.Comment{
				Content:  "11111",
				Pictures: nil,
			},
		},
	}
	c.UploadItem(ac, req)
	c.GetItem(ac, []uint{12})

}
