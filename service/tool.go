package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
)

const (
	Pending = iota
	Pass
	Reject
	PassBeforeHook
	RejectBeforeHook
)

func auditStatusToString(status int) string {
	switch status {
	case Pending:
		return "Pending"
	case Pass:
		return "Pass"
	case Reject:
		return "Reject"
	case PassBeforeHook:
		return "PassBeforeHook"
	case RejectBeforeHook:
		return "RejectBeforeHook"
	default:
		panic("unhandled default case")
	}
}

func auditStatusForHook(status int) string {
	switch status {
	case Pending:
		return "Pending"
	case Pass, PassBeforeHook:
		return "Pass"
	case Reject, RejectBeforeHook:
		return "Reject"
	default:
		panic("unhandled default case")
	}
}

func auditStatusToInt(status string) int {
	switch status {
	case "Pending":
		return Pending
	case "Pass":
		return Pass
	case "Reject":
		return Reject
	case "PassBeforeHook":
		return PassBeforeHook
	case "RejectBeforeHook":
		return RejectBeforeHook
	}
	return -1
}

func hookBack(t string, data request.HookPayload, authorization string) ([]byte, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal hook payload: %w", err)
	}
	var lasterr error
	for i := 0; i < data.Try; i++ {
		reqs, err := http.NewRequest("POST", t, bytes.NewBuffer(jsonBytes))
		if err != nil {
			lasterr = err
			time.Sleep(time.Second)
			continue
		}
		reqs.Header.Set("Content-Type", "application/json")
		if authorization != "" {
			reqs.Header.Set("Authorization", authorization)
		}
		client := &http.Client{}
		resp, err := client.Do(reqs)
		if err != nil {
			lasterr = err
			time.Sleep(time.Second)
			continue
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lasterr = readErr
			break
		}
		if resp.StatusCode == http.StatusOK {
			return body, nil
		}
	}

	return nil, lasterr
}
