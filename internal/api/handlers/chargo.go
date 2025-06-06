package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jilaidian_go/internal/config"
	"jilaidian_go/internal/models"
	"jilaidian_go/internal/service"
	"jilaidian_go/internal/utils"
	"jilaidian_go/pkg/logger"
	"net/http"
	"strings"
	"time"

	. "jilaidian_go/internal/utils"

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

type ValidToken struct {
	Token              string `json:"token"`
	ExpirationTime     int64  `json:"expirationTime"`
	TokenAvailableTime int    `json:"tokenAvailableTime"` // 可用时间，单位秒
}

var validTokens ValidToken

func init() {
	validTokens.ExpirationTime = 0
	validTokens.update()
}

func (vt *ValidToken) update() {
	if vt.ExpirationTime < time.Now().Unix() {
		vt.Token = utils.GenerateSignString(time.Now().Format("2006-01-02 15:04:05"), config.Global.YiAnqi.SignKey)
		vt.ExpirationTime = time.Now().Add(time.Hour).Unix() // 1小时后过期
		vt.TokenAvailableTime = 3600                         // 可用时间，单位秒
	}
	logger.Logger.Info("Updated valid token:", vt.Token, " Expiration Time:", vt.ExpirationTime, " Available Time:", vt.TokenAvailableTime)
}
func (vt *ValidToken) isValidToken(token string) bool {
	return vt.Token == token && vt.ExpirationTime > time.Now().Unix()
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
		if !validTokens.isValidToken(tokenString) {
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

type Message struct {
	Result      int    `json:"result"`
	Description string `json:"description"`
}

type QueryTokenRequest struct {
	OperatorID     string `json:"OperatorID"`
	OperatorSecret string `json:"OperatorSecret"`
}
type QueryTokenResponse struct {
	OperatorID         string `json:"OperatorID"`
	SuccStat           int    `json:"SuccStat"`
	AccessToken        string `json:"AccessToken"`
	TokenAvailableTime int    `json:"TokenAvailableTime"`
	FailReason         int    `json:"FailReason"`
}
type TotalResponse struct {
	Ret  int    `json:"Ret"`
	Msg  string `json:"Msg"`
	Data string `json:"Data"`
	Sig  string `json:"Sig"`
}

func (tr *TotalResponse) MakeSig() {
	/*
		Sig（签名）采用HMAC-MD5算法，采用MD5作为散列函数，
		通过SigSecret（签名密钥）对整个消息主体各参数的值拼接后进行加密，
		入参拼接顺序为：OperatorID（运营商标识）、Data（参数内容）、TimeStamp（时间戳）、Seq（自增序列），
		出参拼接顺序为：Ret（返回值）、Msg（返回信息）、Data（参数内容），
		然后采用MD5信息摘要的方式形成新密文，参数签名必须大写，详见附录B。
	*/
	prestr := fmt.Sprintf("%d%s%s", tr.Ret, tr.Msg, tr.Data)

	tr.Sig = GenerateSignString(prestr, config.Global.YiAnqi.SignKey)
}
func (h *Handler) MakeRepsonse(c *gin.Context, result int, description string, data interface{}) {

	var jsonData []byte
	switch v := data.(type) {
	case string:
		jsonData = []byte(v)
	default:
		jsonData, _ = json.Marshal(data)
	}

	encodedStr, _ := CBCEncrypt_Base64(string(jsonData), config.Global.YiAnqi.AesKey, config.Global.YiAnqi.AesIv)

	tr := TotalResponse{
		Ret:  result,
		Msg:  description,
		Data: encodedStr,
		Sig:  "",
	}
	tr.MakeSig()
	c.JSON(http.StatusOK, tr)
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

// query_token 查询token接口
func (h *Handler) query_token(c *gin.Context) {
	decodedStr, err := h.extractRequest(c)
	if err != nil {
		return
	}
	var queryToken QueryTokenRequest
	if err := json.Unmarshal([]byte(decodedStr), &queryToken); err != nil {
		logger.Logger.Errorf("JSON解析失败 %s  %v", decodedStr, err)
		h.MakeRepsonse(c, 1, "JSON解析失败", "")
		return
	}
	validTokens.update()
	response := QueryTokenResponse{
		OperatorID:         queryToken.OperatorID,
		SuccStat:           1,
		AccessToken:        validTokens.Token,
		TokenAvailableTime: validTokens.TokenAvailableTime,
		FailReason:         0,
	}
	// 将 struct 转为 JSON 字符串
	jsonData, _ := json.Marshal(response)

	fmt.Println(string(jsonData))
	h.MakeRepsonse(c, 0, "查询成功", string(jsonData))
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
	parkinfo := config.GetParkInfo(requestItem.ParkID)
	if parkinfo == nil {
		logger.Logger.Warn("无效车场ID", "parkID", requestItem.ParkID, "plateNo", requestItem.PlateNum)
		h.MakeRepsonse(c, 1, "未授权的车场", respose)
		return
	}

	// 处理充电信息
	if err := h.chargeService.ProcessChargingAndDiscount(parkinfo, requestItem); err != nil {
		logger.Logger.Error("处理充电和优惠失败", err, "parkID", requestItem.ParkID, "plateNo", requestItem.PlateNum)
		h.MakeRepsonse(c, 1, "处理充电和优惠失败", respose)
		return
	}

	logger.Logger.Info("充电记录处理成功", "parkID", requestItem.ParkID, "plateNo", requestItem.PlateNum)
	respose.ConfirmResult = 0  // 假设处理成功后返回确认结果为0
	respose.PlateAutResult = 1 // 假设车牌自动识别结果为1
	h.MakeRepsonse(c, 1, "充电记录处理成功", respose)
}
