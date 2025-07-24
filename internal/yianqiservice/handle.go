package yianqiservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jilaidian_go/internal/config"
	"jilaidian_go/internal/debugHandler"
	"jilaidian_go/internal/models"
	"jilaidian_go/internal/service"
	"jilaidian_go/pkg/logger"
	"net/http"
	"strings"
	"sync"

	. "jilaidian_go/internal/utils"

	"github.com/gin-gonic/gin"
)

// Handler API处理器
type Handler struct {
	chargeService *service.ChargeService
	DebugInfo     debugHandler.DebugInfo
}

var (
	myValidToken ValidToken
	instance     *Handler
	once         sync.Once
)

// NewHandler 创建新的API处理器
func NewHandler(r *gin.Engine) *Handler {
	once.Do(func() {
		instance = &Handler{
			chargeService: service.NewChargeService(),
		}
		instance.SetupRoutes(r)
		instance.DebugInfo = debugHandler.GetDebugInfo("YianqiService")

	})
	return instance
}

func init() {
	myValidToken.ExpirationTime = 0
	myValidToken.update()
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header (e.g., Authorization: Bearer <token>)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			h.MakeRepsonse(c, 4002, "Missing authorization token", "")
			return
		}

		// Check if the header starts with Bearer
		if !strings.HasPrefix(authHeader, "Bearer ") {
			h.MakeRepsonse(c, 4002, "Invalid authorization format", "")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		fmt.Println("Received token:", tokenString)
		// Validate token (example: check against a list of valid tokens or use JWT)
		if !myValidToken.isValidToken(tokenString) {
			h.MakeRepsonse(c, 4002, "Invalid or expired token", "")
			return
		}

		// Continue processing the request
		c.Next()
	}
}

// SetupRoutes 设置路由
func (h *Handler) SetupRoutes(r *gin.Engine) {
	// 充电记录接口
	authorized := r.Group("/")
	authorized.Use(h.AuthMiddleware()) // 使用中间件进行身份验证
	{
		authorized.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Test endpoint"})
		})

		authorized.POST("/notification_charge_end_order_info", h.notification_charge_end_order_info)
	}

	r.POST("/query_token", h.query_token)

}

func (h *Handler) extractRequest(c *gin.Context) (string, error) {
	var queryRequest models.QueryRequest

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Logger.Error("参数解析失败", err)
		h.MakeRepsonse(c, 1, "非法参数", "")
		return "", err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // 重置请求体
	logger.Logger.Debugf("请求体内容: %s", body)
	if err := c.ShouldBindJSON(&queryRequest); err != nil {
		logger.Logger.Error("参数解析失败", err)
		h.MakeRepsonse(c, 1, "非法参数", "")
		return "", err
	}
	if queryRequest.Data == "" {
		logger.Logger.Errorf("queryRequest.Data is empty")
		h.MakeRepsonse(c, 1, "参数解析失败 Data 为空", "")
		return "", err
	}
	decodedStr, err := CBCDecrypt_Base64(queryRequest.Data, config.Global.YiAnqi.AesKey,
		config.Global.YiAnqi.AesIv)
	if err != nil {
		logger.Logger.Errorf("解密失败 %s %v", queryRequest.Data, err)
		h.MakeRepsonse(c, 1, "解密失败", "")
		return "", err
	}
	logger.Logger.Debugf("解密内容: %s", decodedStr)
	return decodedStr, nil
}

// HandleChargingRecord 处理充电记录
func (h *Handler) notification_charge_end_order_info(c *gin.Context) {
	//1.2代理云根据充电信息的车牌号，向停车云3.10接口查询订单，得到orderid后；向3.13接口推送优惠信息；
	// 解析请求体
	decodedStr, err := h.extractRequest(c)
	if err != nil {
		return
	}

	var requestItem models.Notification_charge_end_order_info_Request
	if err := json.Unmarshal([]byte(decodedStr), &requestItem); err != nil {
		logger.Logger.Errorf("JSON解析失败 %s  %v", decodedStr, err)
		h.MakeRepsonse(c, 1, "JSON解析失败", "")
		return
	}
	h.DebugInfo.LastReuqestJson = decodedStr
	respose := models.Notification_charge_end_order_info_Response{
		StartChargeSeq: requestItem.StartChargeSeq,
		ConfirmResult:  1,
		PlateAutResult: 0,
	}
	// 添加输入验证
	if requestItem.PlateNum == "" {
		logger.Logger.Warn("请求缺少必要参数", "plateNo", requestItem.PlateNum)
		h.MakeRepsonse(c, 1, "缺少必要参数", respose)
		return
	}

	// 车场校验逻辑
	parkinfo := config.GetParkInfo(requestItem.StationID)
	if parkinfo == nil {
		logger.Logger.Warn("无效车场ID", "parkID", requestItem.StationID, "plateNo", requestItem.PlateNum)
		h.MakeRepsonse(c, 1, "未授权的车场", respose)
		return
	}

	// 处理充电信息
	if err := h.chargeService.ProcessChargingAndDiscount(parkinfo, requestItem.PlateNum); err != nil {
		logger.Logger.Error("处理充电和优惠失败", err, "parkID", parkinfo.ParkID, "plateNo", requestItem.PlateNum)
		h.MakeRepsonse(c, 1, "处理充电和优惠失败", respose)
		return
	}

	logger.Logger.Info("充电记录处理成功", "parkID", parkinfo.ParkID, "plateNo", requestItem.PlateNum)
	respose.ConfirmResult = 0  // 假设处理成功后返回确认结果为0
	respose.PlateAutResult = 1 // 假设车牌自动识别结果为1
	h.MakeRepsonse(c, 0, "充电记录处理成功", respose)
}

// query_token 查询token接口
func (h *Handler) query_token(c *gin.Context) {
	decodedStr, err := h.extractRequest(c)
	if err != nil {
		return
	}
	var queryToken QueryTokenRequest
	if err := json.Unmarshal([]byte(decodedStr), &queryToken); err != nil {
		logger.Logger.Errorf("query_token JSON解析失败 %s  %v", decodedStr, err)
		h.MakeRepsonse(c, 1, "query_token JSON解析失败", "")
		return
	}
	myValidToken.update()
	response := QueryTokenResponse{
		OperatorID:         queryToken.OperatorID,
		SuccStat:           0,
		AccessToken:        myValidToken.AccessToken,
		TokenAvailableTime: myValidToken.TokenAvailableTime,
		FailReason:         0,
	}
	// 将 struct 转为 JSON 字符串
	jsonData, _ := json.Marshal(response)

	fmt.Println(string(jsonData))
	h.MakeRepsonse(c, 0, "查询成功", string(jsonData))
}
