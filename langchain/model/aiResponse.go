package model

type AIResponse[T any] struct {
	Data   T     `json:"data"`   //ai返回的结果
	Tokens int   `json:"tokens"` //方便后续打点记录token的消耗
	Error  error `json:"error"`  //ai返回的报错
}

func GetAIResponseData[T any](r *AIResponse[T]) (T, error) {
	return r.Data, nil
}
