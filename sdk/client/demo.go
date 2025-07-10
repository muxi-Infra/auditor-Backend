package client

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"net/http"
	"time"
)

func main() {
	ac := "xxxxxx"
	c := NewMuxiAuditClient(&http.Client{}, "http://localhost:8080/api/v1")
	var req = request.UploadReq{
		HookUrl:    "http://localhost:8081",
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
	//req.Content.Topic.Title = "test2"
	//c.UpdateItem(ac, sc, req)
	//c.DeleteItem(ac, sc, 12)
	c.GetItem(ac, []uint{12})

}
