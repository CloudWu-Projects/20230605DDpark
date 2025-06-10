package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jilaidian_go/pkg/logger"
	"net/http"
	"time"
)

type HttpPostClient struct {
	client *http.Client
}

func NewHttpPostClient() *HttpPostClient {
	return &HttpPostClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// sendRequest 发送HTTP请求并解析响应
func (c *HttpPostClient) SendRequest(url string, request interface{}) ([]byte, error) {
	// 序列化请求体
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}
	logger.Logger.Debugf("sendRequest url:%s requestBody:%s", url, string(requestBody))
	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	logger.Logger.Debugf("sendRequest  body %s", body)

	return body, nil
}
