package service

import (
	"errors"
	"fmt"
	"jilaidian_go/internal/config"
	"jilaidian_go/pkg/client"
	"jilaidian_go/pkg/logger"
	"strings"
)

// ChargeService 充电服务
type ChargeService struct {
	apiClient *client.APIClient
}

type OrderLookupResult struct {
	ParkInfo *config.ParkInfo
	OrderID  string
	InTime   int64
}

// NewChargeService 创建新的充电服务
func NewChargeService() *ChargeService {
	return &ChargeService{
		apiClient: client.NewAPIClient(),
	}
}

// ValidateParkID 校验停车场ID

// ProcessChargingAndDiscount 处理充电和优惠
func (s *ChargeService) ProcessChargingAndDiscount(parkinfo *config.ParkInfo, PlateNum string) error {
	// 1. 处理充电信息

	// 2. 查询订单信息

	order_id, err := s.apiClient.QueryOrder(parkinfo.ParkID, PlateNum, parkinfo)
	if err != nil {
		return fmt.Errorf("查询订单信息失败: %v", err)
	}
	err = s.apiClient.SendDiscountNotice(parkinfo.ParkID, PlateNum, order_id, parkinfo)

	return err
}

func (s *ChargeService) QueryFirstMatchingOrder(plateNo string) (*OrderLookupResult, error) {
	plateNo = strings.ToUpper(strings.TrimSpace(plateNo))
	if plateNo == "" {
		return nil, errors.New("请输入车牌号")
	}

	attempted := 0
	errorCount := 0
	var lastErr error

	for i := range config.Global.Parks {
		park := &config.Global.Parks[i]
		if park.ParkID == 0 || strings.TrimSpace(park.GetUkey()) == "" {
			continue
		}

		attempted++
		response, err := s.apiClient.QueryOrderDetail(park.ParkID, plateNo, park)
		if err != nil {
			errorCount++
			lastErr = err
			logger.Logger.Errorf("查询停车订单失败 parkID=%d plate=%s err=%v", park.ParkID, plateNo, err)
			continue
		}

		if response.State != 0 && strings.TrimSpace(response.Data.OrderID) != "" {
			return &OrderLookupResult{
				ParkInfo: park,
				OrderID:  response.Data.OrderID,
				InTime:   response.Data.InTime,
			}, nil
		}
	}

	if attempted == 0 {
		return nil, errors.New("没有可参与查询的车场，请先配置 park_id 和 ukey")
	}
	if errorCount == attempted && lastErr != nil {
		return nil, fmt.Errorf("查询停车订单失败: %w", lastErr)
	}

	return nil, nil
}
