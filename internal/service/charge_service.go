package service

import (
	"fmt"
	"jilaidian_go/internal/config"
	"jilaidian_go/internal/models"
	"jilaidian_go/pkg/client"
	"strconv"
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
func (s *ChargeService) ProcessChargingAndDiscount(parkinfo *config.ParkInfo, chargeData models.Notification_charge_end_order_info_Request) error {
	// 1. 处理充电信息

	// 2. 查询订单信息
	parkID, err := strconv.Atoi(chargeData.ParkID)
	if err != nil {
		return fmt.Errorf("停车场ID格式错误: %v", err)
	}

	order_id, err := s.apiClient.QueryOrder(parkID, chargeData.PlateNum, parkinfo)
	if err != nil {
		return fmt.Errorf("查询订单信息失败: %v", err)
	}
	err = s.apiClient.SendDiscountNotice(parkID, chargeData.PlateNum, order_id, parkinfo)

	return err
}
