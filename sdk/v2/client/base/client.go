package base

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/api/errorx"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/config"
)

const UploadPath = "/remove/upload"
const UpdatePath = "/remove/update"
const DeletePath = "/remove/delete"
const GetPath = "/remove/get"

type Client struct {
	ApiKey string
	Region string
	client *http.Client
}

func (client *Client) Init(c config.Config) error {
	if c.ApiKey == "" || c.Region == "" {
		return errors.New("api key or endpoint is empty")
	}

	timeout := 5 * time.Second
	if c.ConnectTimeout > 0 {
		timeout = time.Duration(c.ConnectTimeout) * time.Millisecond
	}

	client.client = &http.Client{
		Timeout:   timeout,
		Transport: http.DefaultTransport,
	}
	client.Region = c.Region
	client.ApiKey = c.ApiKey
	return nil
}

// request 发送范式，所有真正的请求都调用这个方法。
func (client *Client) doRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("api_key", client.ApiKey)

	resp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}

	// 这里将所有非2**的请求都当做失败处理
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("http request failed: status=%d body=%s", resp.StatusCode, string(b))
	}

	return resp, nil
}

func (client *Client) UploadItem(ctx context.Context, req request.UploadReq) (response.Resp, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return response.Resp{}, errorx.MarshalErr(err)
	}

	resp, err := client.doRequest(ctx, http.MethodPost, client.Region+UploadPath, bytes.NewBuffer(data))
	if err != nil {
		return response.Resp{}, err
	}
	defer resp.Body.Close()
	var res response.Resp
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return response.Resp{}, errorx.MarshalErr(err)
	}

	return res, nil
}

func (client *Client) UpdateItem(ctx context.Context, req request.UpdateReq) (response.Resp, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return response.Resp{}, err
	}

	resp, err := client.doRequest(ctx, http.MethodPut, client.Region+UpdatePath, bytes.NewBuffer(data))
	if err != nil {
		return response.Resp{}, err
	}
	defer resp.Body.Close()

	var res response.Resp
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return response.Resp{}, err
	}

	return res, nil
}

func (client *Client) DeleteItem(ctx context.Context, id int) (response.Resp, error) {
	path := fmt.Sprintf(DeletePath+"/%d", id)
	resp, err := client.doRequest(ctx, http.MethodDelete, client.Region+path, nil)
	if err != nil {
		return response.Resp{}, err
	}
	defer resp.Body.Close()

	var res response.Resp
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return response.Resp{}, err
	}

	return res, nil
}

func (client *Client) GetItems(ctx context.Context, ids []int) (response.Resp, error) {
	query := url.Values{}
	query.Set("ids", uintSliceToString(ids))

	resp, err := client.doRequest(ctx, http.MethodGet, client.Region+GetPath+"?"+query.Encode(), nil)
	if err != nil {
		return response.Resp{}, err
	}

	defer resp.Body.Close()

	var res response.Resp
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return response.Resp{}, err
	}

	return res, nil
}
