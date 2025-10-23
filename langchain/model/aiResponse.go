package model

import "github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/errorx"

type AIResponse struct {
	Data   any   `json:"data"`   //ai返回的结果
	Tokens int   `json:"tokens"` //方便后续打点记录token的消耗
	Error  error `json:"error"`  //ai返回的报错
}

func GetAIResponseData[T any](r *AIResponse) (T, error) {
	v, ok := r.Data.(T)
	if !ok {
		var zero T
		return zero, errorx.ErrDataInvalid
	}
	return v, nil
}
