package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"jilaidian_go/lib/config"
	"jilaidian_go/lib/logger"
	"jilaidian_go/lib/models"
	"jilaidian_go/lib/utils"
	"net/http"
	"time"
)

// APIClient 外部API客户端
type APIClient struct {
	baseURL string
	client  *http.Client
}

// NewAPIClient 创建新的API客户端
func NewAPIClient() *APIClient {
	return &APIClient{
		baseURL: config.Global.API.BaseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// QueryOrder 查询订单
func (c *APIClient) QueryOrder(parkID int, carNumber string, parkinfo *config.ParkInfo) (string, error) {
	url := fmt.Sprintf("%s/order/queryOrder", c.baseURL)

	// 构造请求数据
	data := struct {
		CarNumber string `json:"car_number"`
	}{
		CarNumber: carNumber,
	}
	dataBytes, _ := json.Marshal(data)
	dataStr := string(dataBytes)
	// 生成签名
	sign := utils.GenerateSignString(dataStr, parkinfo.Ukey)

	// 构造请求体
	request := models.QueryOrderRequest{
		ServiceName: "query_order",
		Sign:        sign,
		ParkID:      parkID,
		Data:        data,
	}

	// 发送请求
	response, err := c.sendRequest(url, request)
	if err != nil {
		return "", err
	}
	if response.State == 0 {
		return "", fmt.Errorf("查询订单失败，状态码：%d", response.State)
	}
	return response.Data.OrderID, nil
}

// SendDiscountNotice 下发优惠信息
func (c *APIClient) SendDiscountNotice(parkID int, carNumber, orderID string, parkinfo *config.ParkInfo) error {
	url := fmt.Sprintf("%s/charge/discountNotice", c.baseURL)

	// 构造请求数据
	data := struct {
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
		ReduceAmount:      0,
		DeductionTime:     parkinfo.Deduction_time,
		DeductionMoney:    parkinfo.Deduction_money,
		Duration:          20,
		Remark:            "备注",
		StartChargingTime: "2020-08-27 00:02:09",
		StopChargingTime:  "2020-08-27 00:25:07",
		UUID:              "de6c26a945c9478295d7cffa7631d7f9",
	}

	dataBytes, _ := json.Marshal(data)
	dataStr := string(dataBytes)
	logger.Logger.Info("dataStr:", dataStr)
	// 生成签名
	sign := utils.GenerateSignString(dataStr, parkinfo.Ukey)

	// 构造请求体
	request := models.DiscountNoticeRequest{
		ServiceName: "charge_discount_notice",
		Sign:        sign,
		ParkID:      parkID,
		Data:        data,
	}
	// 发送请求
	response, err := c.sendRequest(url, request)

	if err != nil {
		logger.Logger.Error("发送请求失败", (err))
		return err
	}
	logger.Logger.Debugf(fmt.Sprintf("SendDiscountNotice response %+v", response))
	if response.State == 0 {
		logger.Logger.Error("处理优惠失败", (err))
		return errors.New("处理优惠失败")
	}
	return nil
}

// sendRequest 发送HTTP请求并解析响应
func (c *APIClient) sendRequest(url string, request interface{}) (*models.QueryOrderResponse, error) {
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

	// 解析响应
	var response models.QueryOrderResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	logger.Logger.Debugf("sendRequest  body %s", body)

	return &response, nil
}
