package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/apikey"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const UploadPath = "/remove/upload"
const UpdatePath = "/remove/update"
const DeletePath = "/remove/delete"
const GetPath = "/remove/get"

type MuxiAuditClient struct {
	c   *http.Client
	url string
}

func NewMuxiAuditClient(c *http.Client, ul string) *MuxiAuditClient {
	return &MuxiAuditClient{
		c:   c,
		url: ul, //精确到版本即可
	}
}

// 供调用方上传Item
func (mc *MuxiAuditClient) httpServe(ac string, se string, data []byte, path string, method string) (response.Response, error) {
	rep, err := http.NewRequest(method, mc.url+path, bytes.NewBuffer(data))
	if err != nil {
		return response.Response{}, err
	}
	t := time.Now().Format(time.RFC3339)

	signature := apikey.SignRequest(se, t)
	rep.Header.Set("Content-Type", "application/json")
	rep.Header.Set("Accept", "application/json")
	//注册时产生
	rep.Header.Set("AccessKey", ac)
	//在服务端会进行相似度检验
	rep.Header.Set("Signature", signature)
	rep.Header.Set("Timestamp", t)
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
func (mc *MuxiAuditClient) UploadItem(ac string, se string, req request.UploadReq) (response.Response, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return response.Response{}, err
	}
	resp, err := mc.httpServe(ac, se, data, UploadPath, http.MethodPost)
	if err != nil {
		return response.Response{}, err
	}
	return resp, nil
}
func (mc *MuxiAuditClient) UpdateItem(ac string, se string, req request.UploadReq) (response.Response, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return response.Response{}, err
	}
	resp, err := mc.httpServe(ac, se, data, UpdatePath, http.MethodPut)
	if err != nil {
		return response.Response{}, err
	}
	return resp, nil
}
func (mc *MuxiAuditClient) DeleteItem(ac string, se string, itemId uint) (response.Response, error) {
	path := fmt.Sprintf(DeletePath+"/%d", itemId)
	resp, err := mc.httpServe(ac, se, nil, path, http.MethodDelete)
	if err != nil {
		return response.Response{}, err
	}
	return resp, nil
}
func (mc *MuxiAuditClient) GetItem(ac string, se string, ids []uint) (response.Response, error) {
	query := url.Values{}
	query.Set("ids", uintSliceToString(ids))
	resp, err := mc.httpServe(ac, se, nil, GetPath+"?"+query.Encode(), http.MethodGet)
	if err != nil {
		return response.Response{}, err
	}
	return resp, nil
}
