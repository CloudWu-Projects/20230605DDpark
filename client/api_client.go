package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jilaidian_go/config"
	"jilaidian_go/models"
	"jilaidian_go/utils"
	"net/http"
	"time"
)

// APIClient 外部API客户端
type APIClient struct {
	baseURL string
	ukey    string
	client  *http.Client
}

// NewAPIClient 创建新的API客户端
func NewAPIClient() *APIClient {
	return &APIClient{
		baseURL: config.Global.API.BaseURL,
		ukey:    config.Global.API.UKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// QueryOrder 查询订单
func (c *APIClient) QueryOrder(parkID int, carNumber string) (*models.QueryOrderResponse, error) {
	url := fmt.Sprintf("%s/order/queryOrder", c.baseURL)

	// 构造请求数据
	data := map[string]interface{}{
		"car_number": carNumber,
		"query_time": time.Now().Unix(),
	}

	// 生成签名
	sign := utils.GenerateSign(data, c.ukey)

	// 构造请求体
	request := models.QueryOrderRequest{
		ServiceName: "query_order",
		Sign:        sign,
		ParkID:      parkID,
		Data: struct {
			CarNumber string `json:"car_number"`
			QueryTime int64  `json:"query_time"`
		}{
			CarNumber: carNumber,
			QueryTime: time.Now().Unix(),
		},
	}

	// 发送请求
	return c.sendRequest(url, request)
}

// SendDiscountNotice 下发优惠信息
func (c *APIClient) SendDiscountNotice(parkID int, carNumber, orderID string, reduceAmount float64, deductionTime, deductionMoney int) error {
	url := fmt.Sprintf("%s/charge/discountNotice", c.baseURL)

	// 构造请求数据
	data := map[string]interface{}{
		"car_number":          carNumber,
		"order_id":            orderID,
		"reduce_amount":       reduceAmount,
		"deduction_time":      deductionTime,
		"deduction_money":     deductionMoney,
		"duration":            20,
		"remark":              "备注",
		"start_charging_time": "2020-08-27 00:02:09",
		"stop_charging_time":  "2020-08-27 00:25:07",
		"uuid":                "de6c26a945c9478295d7cffa7631d7f9",
	}

	// 生成签名
	sign := utils.GenerateSign(data, c.ukey)

	// 构造请求体
	request := models.DiscountNoticeRequest{
		ServiceName: "charge_discount_notice",
		Sign:        sign,
		ParkID:      parkID,
		Data: struct {
			CarNumber         string  `json:"car_number"`
			OrderID           string  `json:"order_id"`
			ReduceAmount      float64 `json:"reduce_amount"`
			DeductionTime     int     `json:"deduction_time"`
			DeductionMoney    int     `json:"deduction_money"`
			Duration          int     `json:"duration"`
			Remark            string  `json:"remark"`
			StartChargingTime string  `json:"start_charging_time"`
			StopChargingTime  string  `json:"stop_charging_time"`
			UUID              string  `json:"uuid"`
		}{
			CarNumber:         carNumber,
			OrderID:           orderID,
			ReduceAmount:      reduceAmount,
			DeductionTime:     deductionTime,
			DeductionMoney:    deductionMoney,
			Duration:          20,
			Remark:            "备注",
			StartChargingTime: "2020-08-27 00:02:09",
			StopChargingTime:  "2020-08-27 00:25:07",
			UUID:              "de6c26a945c9478295d7cffa7631d7f9",
		},
	}

	// 发送请求
	_, err := c.sendRequest(url, request)
	return err
}

// sendRequest 发送HTTP请求并解析响应
func (c *APIClient) sendRequest(url string, request interface{}) (*models.QueryOrderResponse, error) {
	// 序列化请求体
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

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

	// 解析响应
	var response models.QueryOrderResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &response, nil
}
