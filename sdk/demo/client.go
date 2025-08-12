package main

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/client"
	"net/http"
	"time"
)

func main() {
	ac := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NTIxNDE2OTYsInN1YiI6NH0.6bNgD_MF1zp8oNcpBPJnKaKU2i4-BdCRdqKNoBiU5Ys"
	c := client.NewMuxiAuditClient(&http.Client{}, "http://60.205.12.92:8080/api/v1")
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
