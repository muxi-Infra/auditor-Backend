package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/muxi-Infra/auditor-Backend/api/request"
	"github.com/muxi-Infra/auditor-Backend/langchain/model"
)

func auditStatusToString(status int) string {
	switch status {
	case model.Pending:
		return "Pending"
	case model.Pass:
		return "Pass"
	case model.Reject:
		return "Reject"
	case model.PassBeforeHook:
		return "PassBeforeHook"
	case model.RejectBeforeHook:
		return "RejectBeforeHook"
	default:
		panic("unhandled default case")
	}
}

func auditStatusForHook(status int) string {
	switch status {
	case model.Pending:
		return "Pending"
	case model.Pass, model.PassBeforeHook:
		return "Pass"
	case model.Reject, model.RejectBeforeHook:
		return "Reject"
	default:
		panic("unhandled default case")
	}
}

func auditStatusToInt(status string) int {
	switch status {
	case "Pending":
		return model.Pending
	case "Pass":
		return model.Pass
	case "Reject":
		return model.Reject
	case "PassBeforeHook":
		return model.PassBeforeHook
	case "RejectBeforeHook":
		return model.RejectBeforeHook
	}
	return -1
}

func hookBack(t string, data request.HookPayload, authorization string) ([]byte, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal hook payload: %w", err)
	}

	var lastErr error
	var lastCode int
	client := &http.Client{}
	for i := 0; i < data.Try; i++ {
		reqs, err := http.NewRequest("POST", t, bytes.NewBuffer(jsonBytes))
		if err != nil {
			lastErr = err
			time.Sleep(time.Second)
			continue
		}
		reqs.Header.Set("Content-Type", "application/json")
		if authorization != "" {
			reqs.Header.Set("Authorization", authorization)
		}
		resp, err := client.Do(reqs)
		if err != nil {
			lastErr = err
			time.Sleep(time.Second)
			continue
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			break
		}
		// 确保只有这一条路径可以无错误返回
		if resp.StatusCode == http.StatusOK {
			return body, nil
		}
		lastCode = resp.StatusCode
	}

	return nil, fmt.Errorf("failed to call hook back: %w, lastStatusCode: %d", lastErr, lastCode)
}

func envInt(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}
