package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"jilaidian_go/internal/config"
	"jilaidian_go/internal/models"
	"jilaidian_go/internal/utils"
	"jilaidian_go/pkg/logger"

	"github.com/Microsoft/go-winio/pkg/guid"
)

// APIClient 外部API客户端
type APIClient struct {
	TingCheYunUrl string
	client        *HttpPostClient
}

// NewAPIClient 创建新的API客户端
func NewAPIClient() *APIClient {
	return &APIClient{
		TingCheYunUrl: config.Global.ServerConfig.TingCheYunUrl,
		client:        NewHttpPostClient(),
	}
}

func (c *APIClient) QueryOrderDetail(parkID int, carNumber string, parkinfo *config.ParkInfo) (*models.QueryOrderResponse, error) {
	url := fmt.Sprintf("%s/order/queryOrder", c.TingCheYunUrl)

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

	return c.SendRequest(url, request)
}

// QueryOrder 查询订单
func (c *APIClient) QueryOrder(parkID int, carNumber string, parkinfo *config.ParkInfo) (string, error) {
	response, err := c.QueryOrderDetail(parkID, carNumber, parkinfo)
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
	url := fmt.Sprintf("%s/charge/discountNotice", c.TingCheYunUrl)

	uuid, _ := guid.NewV4()
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
		ReduceAmount:      0, //parkinfo.ReduceAmount,
		DeductionTime:     parkinfo.DeductionTime,
		DeductionMoney:    parkinfo.DeductionMoney,
		Duration:          parkinfo.Duration,
		Remark:            parkinfo.Remark,
		StartChargingTime: "2020-08-27 00:02:09",
		StopChargingTime:  "2020-08-27 00:25:07",
		UUID:              uuid.String(), //GUID//"de6c26a945c9478295d7cffa7631d7f9",
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
	response, err := c.SendRequest(url, request)

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
func (c *APIClient) SendRequest(url string, request interface{}) (*models.QueryOrderResponse, error) {

	body, err := c.client.SendRequest(url, request)
	if err != nil {
		logger.Logger.Error("发送请求失败", err)
		return nil, err
	}
	// 解析响应
	var response models.QueryOrderResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &response, nil
}
