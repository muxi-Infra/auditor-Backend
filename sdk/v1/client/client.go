package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/muxi-Infra/auditor-Backend/api/request"
	"github.com/muxi-Infra/auditor-Backend/api/response"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const UploadPath = "/remove/upload"
const UpdatePath = "/remove/update"
const DeletePath = "/remove/delete"
const GetPath = "/remove/get"

type MuxiAuditClient struct {
	c   *http.Client
	url string
	ctx context.Context
}

func NewMuxiAuditClient(c *http.Client, ul string) *MuxiAuditClient {
	return &MuxiAuditClient{
		c:   c,
		url: ul, //精确到版本即可
	}
}

// 供调用方上传Item
func (mc *MuxiAuditClient) httpServe(api_key string, data []byte, path string, method string) (response.Response, error) {
	rep, err := http.NewRequest(method, mc.url+path, bytes.NewBuffer(data))
	if err != nil {
		return response.Response{}, err
	}

	rep.Header.Set("Content-Type", "application/json")
	rep.Header.Set("Accept", "application/json")
	//注册时产生
	rep.Header.Set("api_key", api_key)
	re, err := mc.c.Do(rep)
	if err != nil {
		return response.Response{}, err
	}
	defer re.Body.Close()
	data, err = io.ReadAll(re.Body)
	if err != nil {
		return response.Response{}, err
	}
	fmt.Println(string(data))
	var resp response.Response
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return response.Response{}, err
	}
	return resp, nil
}

// 工具函数：将 []uint 转换为逗号分隔的字符串
func uintSliceToString(ids []uint) string {
	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = strconv.FormatUint(uint64(id), 10)
	}
	return strings.Join(strs, ",")
}

func (mc *MuxiAuditClient) UploadItem(apiKey string, req request.UploadReq) (response.Response, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return response.Response{}, err
	}
	resp, err := mc.httpServe(apiKey, data, UploadPath, http.MethodPost)
	if err != nil {
		return response.Response{}, err
	}
	return resp, nil
}

func (mc *MuxiAuditClient) UpdateItem(ac string, req request.UploadReq) (response.Response, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return response.Response{}, err
	}
	resp, err := mc.httpServe(ac, data, UpdatePath, http.MethodPut)
	if err != nil {
		return response.Response{}, err
	}
	return resp, nil
}

func (mc *MuxiAuditClient) DeleteItem(ac string, itemId uint) (response.Response, error) {
	path := fmt.Sprintf(DeletePath+"/%d", itemId)
	resp, err := mc.httpServe(ac, nil, path, http.MethodDelete)
	if err != nil {
		return response.Response{}, err
	}
	return resp, nil
}

func (mc *MuxiAuditClient) GetItem(ac string, ids []uint) (response.Response, error) {
	query := url.Values{}
	query.Set("ids", uintSliceToString(ids))
	resp, err := mc.httpServe(ac, nil, GetPath+"?"+query.Encode(), http.MethodGet)
	if err != nil {
		return response.Response{}, err
	}
	return resp, nil
}
