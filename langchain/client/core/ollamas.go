package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/model"
)

type OllamaClient struct {
	url string
}

func NewOllamaClient(url string) *OllamaClient {
	return &OllamaClient{url: url}
}

func (c *OllamaClient) SendMessage(prompt string) (model.AIResponse, error) {
	body := fmt.Sprintf(`{"prompt": "%s"}`, prompt)
	resp, err := http.Post(c.url+"/api/generate", "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return model.AIResponse{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.AIResponse{}, err
	}
	var response model.AIResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return model.AIResponse{}, err
	}
	return response, nil
}
