package api

import (
	"fmt"
	"jilaidian_go/config"
	"jilaidian_go/logger"
	"jilaidian_go/models"
	"jilaidian_go/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler API处理器
type Handler struct {
	chargeService *service.ChargeService
}

// NewHandler 创建新的API处理器
func NewHandler() *Handler {
	return &Handler{
		chargeService: service.NewChargeService(),
	}
}

// SetupRoutes 设置路由
func (h *Handler) SetupRoutes(router *gin.Engine) {
	// 充电记录接口
	router.POST("/chargePile/chargingRecord", h.HandleChargingRecord)
}

// HandleChargingRecord 处理充电记录
func (h *Handler) HandleChargingRecord(c *gin.Context) {
	var record models.ChargeInfo
	if err := c.ShouldBindJSON(&record); err != nil {
		logger.Logger.Error("参数解析失败", (err))
		c.JSON(http.StatusBadRequest, gin.H{"result": 1, "description": "非法参数"})
		return
	}

	// 车场校验逻辑
	parkinfo := config.GetParkInfo(record.ParkID)
	if parkinfo == nil {
		logger.Logger.Warnf(fmt.Sprintf("无效车场ID parkId: %+v", record))
		c.JSON(http.StatusForbidden, gin.H{"result": 1, "description": "未授权的车场"})
		return
	}

	// 处理充电信息
	//
	if err := h.chargeService.ProcessChargingAndDiscount(parkinfo, record); err != nil {
		logger.Logger.Error("处理充电和优惠失败", (err))
		c.JSON(http.StatusInternalServerError, gin.H{"result": 1, "description": "处理充电和优惠失败"})
	}
	c.JSON(http.StatusOK, gin.H{"result": 0, "description": "充电记录处理成功"})
}
