package service

import (
	"fmt"
	"jilaidian_go/internal/config"
	"jilaidian_go/pkg/client"
)

// ChargeService 充电服务
type ChargeService struct {
	apiClient *client.APIClient
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
