package NanjingNengRui

import (
	"encoding/json"
	"io"
	"jilaidian_go/internal/config"
	"jilaidian_go/internal/service"
	"jilaidian_go/internal/utils"
	"jilaidian_go/pkg/logger"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// Handler API处理器
type Handler struct {
	chargeService *service.ChargeService
}

var (
	instance *Handler
	once     sync.Once
)

// NewHandler 创建新的API处理器
func NewHandler(r *gin.Engine) *Handler {
	once.Do(func() {
		instance = &Handler{
			chargeService: service.NewChargeService(),
		}
		instance.SetupRoutes(r)
	})
	return instance
}

func (h *Handler) SetupRoutes(r *gin.Engine) {
	// 充电记录接口

	r.POST("/api/wec/SyncChargePilePay", h.SyncChargePilePay)

}

type SyncChargePilePayRequest struct {
	// 	appId	string	是	用户id，不允许为空
	AppId string `json:"appId" binding:"required"`
	// key	string	是	验证码，参照安全性验证
	Key string `json:"key" binding:"required"`
	// parkId	int	是	车场Id
	ParkId int `json:"parkId" binding:"required"`
	// serviceCode	string	否	syncChargePilePay
	ServiceCode string `json:"serviceCode"`
	// ts	string	否	时间戳
	Ts string `json:"ts"`
	// reqId	String	否	每次请求的唯一标识
	ReqId string `json:"reqId"`
	// orderNo	string	是	订单号(充电桩订单号),
	OrderNo string `json:"orderNo" binding:"required"`
	// 建议车场保证唯一
	// 不允许为空
	// plateNo	string	是	车牌号(全车牌)
	PlateNo string `json:"plateNo" binding:"required"`
	// 不允许为空
	// startTime	string	是	开始时间 不允许为空
	StartTime string `json:"startTime" binding:"required"`
	// yyyy-MM-dd HH:mm:ss
	// endTime	string	是	结束时间 不允许为空
	EndTime string `json:"endTime" binding:"required"`
	// yyyy-MM-dd HH:mm:ss
	// stationId	string	是	电站ID
	StationId string `json:"stationId" binding:"required"`
	// stationName	string	是	电站名称
	StationName string `json:"stationName" binding:"required"`
	// deviceId	string	是	设备ID
	DeviceId string `json:"deviceId" binding:"required"`
	// deviceName	string	是	设备名称
	DeviceName string `json:"deviceName" binding:"required"`
	// spaceNo	String	是	车位号,停放位置
	SpaceNo string `json:"spaceNo" binding:"required"`
	// power	float	是	充电量 单位：度，小数点后2位
	Power float64 `json:"power" binding:"required"`
	// elecMoney	int	是	电费 单位：分
	ElecMoney int `json:"elecMoney" binding:"required"`
	// seviceMoney	int	是	服务费/附加费 单位：分
	SeviceMoney int `json:"seviceMoney" binding:"required"`
	// totalMoney	int	是	总费用 单位：分
	TotalMoney int `json:"totalMoney" binding:"required"`
	// freeType	int	是	值：0或1
	FreeType int `json:"freeType" binding:"required"`
	// 0：根据freemoeny、freeTime减免
	// 1：本地车场配置减免规则，根据充电时长换算，freemoney、freetime可以为0
	// freeMoney	int	否	减免停车金额 单位：分
	FreeMoney int `json:"freeMoney"`
	// freeTime	int	否	减免停车时长单位：秒
	FreeTime int `json:"freeTime"`
}
type SyncChargePilePayResponse struct {
	// 	resCode	string	响应代码 0=成功，其他失败，详见返回码枚举
	ResCode string `json:"resCode"`
	// resMsg	string	响应消息
	ResMsg string `json:"resMsg"`
	// data	string	响应内容，json字符串
	Data string `json:"data"`
}

func checkSign(jsonBody []byte, req SyncChargePilePayRequest, appSercert string) bool {
	// 	1、将json的所有非空属性名(属性值不为null且不为空字符串以及appId与appSercert除外)ASCII码从小到大排序（字典序），使用URL键值对的格式（即key1=value1&key2=value2…）拼接成字符串stringA；注意（参数为空值、json对象、数组的，不参与加密）
	stringA := ""
	if req.AppId != "" {
		stringA += "appId=" + req.AppId + "&"
	}
	params := map[string]string{}

	//req to  map
	json.Unmarshal([]byte(jsonBody), &params)
	//属性值不为null且不为空字符串以及appId与appSercert除外
	for key, value := range params {
		if value != "" && key != "appId" && key != "appSercert" {
			stringA += key + "=" + value + "&"
		}
	}

	// 2、stringA用“&”拼接上appSercert得到stringSignTemp字符串，并对stringSignTemp进行MD5运算，再将得到的字符串所有字符转换为大写，得到sign值signValue。
	stringSignTemp := strings.Join([]string{stringA, config.Global.NanjingNengRui.AppSercert}, "&")

	signValue := utils.MD5(stringSignTemp)
	// 3、将signValue与请求参数中的sign进行比较，如果相同则认为请求合法。
	if strings.EqualFold(signValue, req.Key) {
		return true
	}
	logger.Logger.Warn("签名校验失败", "expectedSign", signValue, "receivedSign", req.Key)
	logger.Logger.Warn("请求参数", stringSignTemp)
	return false
}
func (h *Handler) SyncChargePilePay(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}
	var req SyncChargePilePayRequest

	err = json.Unmarshal(body, &req)
	if err != nil {
		h.MakeRepsonse(c, 1201, "json解析失败"+err.Error(), nil)
		return
	}
	//sign check
	if !checkSign(body, req, config.Global.NanjingNengRui.AppSercert) {
		h.MakeRepsonse(c, 1101, "签名校验失败", nil)
		return
	}
	// 车场校验逻辑
	parkinfo := config.GetParkInfo_withParkid(req.ParkId)
	if parkinfo == nil {
		logger.Logger.Warn("无效车场ID", "parkID", req.ParkId, "plateNo", req.PlateNo)
		h.MakeRepsonse(c, 1102, "未授权的车场", req)
		return
	}

	// 处理充电信息
	if err := h.chargeService.ProcessChargingAndDiscount(parkinfo, req.PlateNo); err != nil {
		logger.Logger.Error("处理充电和优惠失败", err, "parkID", parkinfo.ParkID, "plateNo", req.PlateNo)
		h.MakeRepsonse(c, 1, "处理充电和优惠失败", nil)
		return
	}

	logger.Logger.Info("充电记录处理成功", "parkID", parkinfo.ParkID, "plateNo", req.PlateNo)

	h.MakeRepsonse(c, 0, "充电记录处理成功", nil)
}
func (h *Handler) MakeRepsonse(c *gin.Context, result int, description string, data interface{}) {

	res := &SyncChargePilePayResponse{
		ResCode: strconv.Itoa(result),
		ResMsg:  description,
	}
	c.JSON(http.StatusOK, res)
}
