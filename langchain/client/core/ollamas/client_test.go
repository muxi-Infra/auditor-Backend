package ollamas

import (
	"fmt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"testing"
)

func TestSendMessage(t *testing.T) {
	c := NewOllamaClient("http://localhost:11434/api/generate", "deepseek-r1:1.5b")
	re, err := c.SendMessage(c.Transform("不违法", response.Contents{
		Topic: response.Topics{
			Title:   "你好",
			Content: "大家都好",
		},
		LastComment: response.Comment{},
		NextComment: response.Comment{},
	}))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(re)
}
