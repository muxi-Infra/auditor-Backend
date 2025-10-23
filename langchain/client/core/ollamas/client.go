package ollamas

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/client/service"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/model"
)

type OllamaClient struct {
	url     string
	model   string
	service service.ServiceInterface
}

func NewOllamaClient(url string, model string) *OllamaClient {
	return &OllamaClient{url: url, model: model}
}

func (c *OllamaClient) SendMessage(prompt string) (model.AuditResult, error) {
	body := map[string]interface{}{
		"model":  c.model,
		"prompt": prompt,
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
