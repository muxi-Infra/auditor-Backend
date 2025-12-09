package ollamas

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/muxi-Infra/auditor-Backend/api/response"
	"github.com/muxi-Infra/auditor-Backend/langchain/model"
	"github.com/muxi-Infra/auditor-Backend/langchain/prompt"
	"github.com/muxi-Infra/auditor-Backend/pkg/logger"
)

type OllamaClient struct {
	url   string
	model string
}

func NewOllamaClient(url string, model string) *OllamaClient {
	return &OllamaClient{url: url, model: model}
}

func (c *OllamaClient) SendMessage(content string, pic []string) (model.AuditResult, error) {
	body := map[string]interface{}{
		"model":  c.model,
		"prompt": content, // 这里是prompt
		"format": "json",
		"stream": true,
	}
	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(c.url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return model.AuditResult{}, err
	}
	defer resp.Body.Close()

	var result strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			fmt.Println("JSON parse error:", err)
			continue
		}
		if part, ok := obj["response"].(string); ok {
			result.WriteString(part)
		}
		if done, ok := obj["done"].(bool); ok && done {
			break
		}
	}
	var response model.AuditResult
	err = json.Unmarshal([]byte(result.String()), &response)
	if err = scanner.Err(); err != nil {
		return model.AuditResult{}, err
	}
	return response, nil
}

func (c *OllamaClient) Transform(role string, contents response.Contents) (string, []string) {
	return prompt.BuildPrompt(role, contents), nil
}

func (c *OllamaClient) WrapLogger(logger logger.Logger) { return }
