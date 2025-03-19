package api

import (
	"encoding/json"
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

	// 记录请求JSON
	reqJSON, err := json.Marshal(record)
	if err != nil {
		logger.Logger.Error("JSON序列化失败", (err))
	} else {
		logger.Logger.Info("客户端请求数据", ("request"), (string(reqJSON)))
	}

	if err := c.ShouldBindJSON(&record); err != nil {
		logger.Logger.Error("参数解析失败", (err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "非法参数"})
		return
	}

	// 车场校验逻辑
	if !h.chargeService.ValidateParkID(record.ParkID) {
		logger.Logger.Warn("无效车场ID parkId:", record.ParkID)
		c.JSON(http.StatusForbidden, gin.H{"error": "未授权的车场"})
		return
	}

	// 处理充电信息
	//
	if err := h.chargeService.ProcessChargingAndDiscount(record); err != nil {
		logger.Logger.Error("处理充电和优惠失败", (err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "处理充电和优惠失败"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "充电记录处理成功"})

	// 异步处理优惠下发
	go func() {

	}()
}
