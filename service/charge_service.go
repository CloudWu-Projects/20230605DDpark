package service

import (
	"fmt"
	"jilaidian_go/client"
	"jilaidian_go/config"
	"jilaidian_go/logger"
	"jilaidian_go/models"
	"strconv"

	"go.uber.org/zap"
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
func (s *ChargeService) ValidateParkID(parkID string) bool {
	allowedIDs := config.GetAllowedParkIDs()
	for _, id := range allowedIDs {
		if id == parkID {
			return true
		}
	}
	return false
}

// HandleChargeInfo 处理充电信息
func (s *ChargeService) HandleChargeInfo(chargeData models.ChargeInfo) error {
	// 校验 parkId
	if !s.ValidateParkID(chargeData.ParkID) {
		err := fmt.Errorf("无效的停车场ID: %s", chargeData.ParkID)
		logger.Logger.Error("校验停车场ID失败", zap.String("parkId", chargeData.ParkID))
		return err
	}

	// 记录充电信息
	logger.Logger.Info("接收到充电信息",
		zap.String("parkId", chargeData.ParkID),
		zap.String("plateNo", chargeData.PlateNo))

	// 这里可以添加更多的处理逻辑，如存储到数据库等
	return nil
}

// ProcessChargingAndDiscount 处理充电和优惠
func (s *ChargeService) ProcessChargingAndDiscount(chargeData models.ChargeInfo) error {
	// 1. 处理充电信息
	if err := s.HandleChargeInfo(chargeData); err != nil {
		return err
	}

	// 2. 查询订单信息
	parkID, err := strconv.Atoi(chargeData.ParkID)
	if err != nil {
		return fmt.Errorf("停车场ID格式错误: %v", err)
	}

	orderInfo, err := s.apiClient.QueryOrder(parkID, chargeData.PlateNo)
	if err != nil {
		logger.Logger.Error("查询订单失败", zap.Error(err))
		return err
	}

	// 3. 下发优惠信息
	if orderInfo.State == 1 { // 查询成功
		reduceAmount := 8.0
		deductionTime := 4
		deductionMoney := 5

		err := s.apiClient.SendDiscountNotice(
			parkID,
			chargeData.PlateNo,
			orderInfo.Data.OrderID,
			reduceAmount,
			deductionTime,
			deductionMoney,
		)
		if err != nil {
			logger.Logger.Error("下发优惠失败", zap.Error(err))
			return err
		}
		logger.Logger.Info("优惠下发成功",
			zap.String("plateNo", chargeData.PlateNo),
			zap.String("orderId", orderInfo.Data.OrderID))
	} else {
		logger.Logger.Warn("订单查询未成功，不下发优惠",
			zap.Int("state", orderInfo.State),
			zap.String("errMsg", orderInfo.ErrMsg))
	}

	return nil
}
