package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/client"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/internal"
)

func main() {
	c, err := client.NewClient(client.Config{
		ApiKey:         "YOUR_API_KEY",
		ConnectTimeout: 3000,
		Region:         "http://localhost:8080/api/v1",
	})
	if err != nil {
		panic(err)
	}
	// upload 的使用实例：
	sendToAudit(c, 1, "hao", "test", "这是测试内容", []string{"http://lib.cqhasy.top/0-1758728125.jpeg"})

	// getItemStatus 的使用实例：
	getAuditStatus(c, []int{1})

	// updateItem 的使用实例：
	ur, _ := request.NewUpdateReq(1, request.WithUpdateAuthor("chen"),
		request.WithUpdateContent(internal.NewContents(internal.WithTopicText("update_test",
			"这是更新后的测试内容^^"),
			internal.WithTopicPictures([]string{"http://lib.cqhasy.top/up0-1758728125.jpeg"}),
			internal.WithLastCommentText("comment"))),
	)
	updateItem(c, ur)

	// deleteItem 的使用实例：
	//dr, _ := request.NewDeleteReq(1)
	//deleteItem(c, dr)
}

func sendToAudit(c *client.Client, id uint, author, title, content string, pics []string) {
	con := internal.NewContents(
		internal.WithTopicText(title, content),
		internal.WithTopicPictures(pics))

	req, err := request.NewUploadReq("http://localhost:8081/api/v1/webhook ", id, con, request.WithUploadAuthor(author))
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := c.UploadItem(context.Background(), req)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.Basic.Errorx != nil {
		log.Println(resp.Basic.Errorx)
		return
	}

	fmt.Println(resp)
}

func getAuditStatus(c *client.Client, ids []int) {
	req, err := request.NewGetItemsStatusReq(ids)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := c.GetItems(context.Background(), req)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.Basic.Errorx != nil {
		log.Println(resp.Basic.Errorx)
		return
	}

	fmt.Println(resp)
}

func updateItem(c *client.Client, req *request.UpdateReq) {
	resp, err := c.UpdateItem(context.Background(), req)

	if err != nil {
		log.Println(err)
		return
	}

	if resp.Basic.Errorx != nil {
		log.Println(resp.Basic.Errorx)
		return
	}

	fmt.Println(resp)
}

func deleteItem(c *client.Client, req *request.DeleteReq) {
	resp, err := c.DeleteItem(context.Background(), req)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.Basic.Errorx != nil {
		log.Println(resp.Basic.Errorx)
		return
	}

	fmt.Println(resp)
}
