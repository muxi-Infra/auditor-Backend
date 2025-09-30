package core

type OpenAIClient struct {
	apiKey string
	url    string
}

func NewOpenAIClient(apiKey string, url string) *OpenAIClient {
	return &OpenAIClient{apiKey, url}
}

func (o *OpenAIClient) SendMessage(message string) error {}
