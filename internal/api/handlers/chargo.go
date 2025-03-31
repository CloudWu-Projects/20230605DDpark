package api

import (
	"encoding/json"
	"jilaidian_go/internal/config"
	"jilaidian_go/internal/models"
	"jilaidian_go/internal/service"
	"jilaidian_go/pkg/logger"
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
func (h *Handler) SetupRoutes(r *gin.Engine) {
	// 充电记录接口
	r.POST("/chargePile/chargingRecord", h.HandleChargingRecord)
}

type Message struct {
	Result      int    `json:"result"`
	Description string `json:"description"`
}

func httpResponse(message Message, w http.ResponseWriter) {

	//httpResponse(Message{Result: 1, Description: "非法参数"}, w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}

// HandleChargingRecord 处理充电记录
func (h *Handler) HandleChargingRecord(c *gin.Context) {
	var record models.ChargeInfo

	// 解析请求体
	if err := c.ShouldBindJSON(&record); err != nil {
		logger.Logger.Error("参数解析失败", err)
		c.JSON(http.StatusBadRequest, gin.H{"Result": 1, "Description": "非法参数"})
		return
	}

	// 添加输入验证
	if record.ParkID == "" || record.PlateNo == "" {
		logger.Logger.Warn("请求缺少必要参数", "parkID", record.ParkID, "plateNo", record.PlateNo)
		c.JSON(http.StatusBadRequest, gin.H{"Result": 1, "Description": "缺少必要参数"})
		return
	}

	// 车场校验逻辑
	parkinfo := config.GetParkInfo(record.ParkID)
	if parkinfo == nil {
		logger.Logger.Warn("无效车场ID", "parkID", record.ParkID, "plateNo", record.PlateNo)
		c.JSON(http.StatusBadRequest, gin.H{"Result": 1, "Description": "未授权的车场"})
		return
	}

	// 处理充电信息
	if err := h.chargeService.ProcessChargingAndDiscount(parkinfo, record); err != nil {
		logger.Logger.Error("处理充电和优惠失败", err, "parkID", record.ParkID, "plateNo", record.PlateNo)
		c.JSON(http.StatusInternalServerError, gin.H{"Result": 1, "Description": "处理充电和优惠失败"})
		return
	}

	logger.Logger.Info("充电记录处理成功", "parkID", record.ParkID, "plateNo", record.PlateNo)
	c.JSON(http.StatusOK, gin.H{"Result": 0, "Description": "充电记录处理成功"})
}
