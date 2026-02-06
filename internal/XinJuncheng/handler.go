package xinjuncheng

import (
	"fmt"
	"jilaidian_go/internal/config"
	"jilaidian_go/internal/debugHandler"
	"jilaidian_go/internal/service"
	"jilaidian_go/pkg/logger"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	chargeService *service.ChargeService
	DebugInfo     *debugHandler.DebugInfo
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
		instance.DebugInfo = debugHandler.GetDebugInfo("YianqiService")

	})
	return instance
}

func init() {
}

type XjcResponse struct {
	Result         int    `json:"result"`
	ResultMsg      string `json:"resultMsg"`
	DiscountNumber string `json:"discountNumber"` // 优惠流水号
}

func (h *Handler) MakeRepsonse(c *gin.Context, tr XjcResponse) {
	h.DebugInfo.LastResponseJson = tr
	c.JSON(http.StatusOK, tr)
}

// SetupRoutes 设置路由
func (h *Handler) SetupRoutes(r *gin.Engine) {
	// 充电记录接口
	// authorized := r.Group("/")
	// authorized.Use(h.AuthMiddleware()) // 使用中间件进行身份验证
	// {
	// 	authorized.POST("/test", func(c *gin.Context) {
	// 		c.JSON(http.StatusOK, gin.H{"message": "Test endpoint"})
	// 	})

	// 	authorized.POST("/api/wec/XJCdiscount", h.XJCdiscount)
	// }

	r.POST("/api/wec/XJCdiscount", h.XJCdiscount)

}
func (h *Handler) XJCdiscount(c *gin.Context) {

	var request XjcRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.MakeRepsonse(c, XjcResponse{
			Result:    0,
			ResultMsg: fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}
	h.DebugInfo.LastReuqestJson = request
	logger.Logger.Infof("Received XJCdiscount request: %+v", request)

	//check sign here if needed
	if !request.VerifySign(config.Global.XinJunCheng.VKey) {
		h.MakeRepsonse(c, XjcResponse{Result: 0, ResultMsg: "Invalid sign", DiscountNumber: request.DiscountNumber})
		return
	}
	// 车场校验逻辑
	parkinfo := config.GetParkInfo_withParkid(request.PlatformNo)
	if parkinfo == nil {
		logger.Logger.Warn("无效车场ID", "parkID", request.PlatformNo, "plateNo", request.PlateNumber)
		h.MakeRepsonse(c, XjcResponse{
			Result:         0,
			ResultMsg:      "未授权的车场",
			DiscountNumber: request.DiscountNumber,
		})
		return
	}
	// Here you would process the discount logic
	// For demonstration, we'll assume it's always successful
	// 3.2.2 代理云通过platformNo、 plateNumber；,调用停车云3.10接口获取order_id
	// park_id=10063417
	// ukey=YCE6JPHTH3D3ZSMD
	// url 如下
	// http://istparking.sciseetech.com/public/order/queryOrder
	// 阿里云向停车云查询订单，park_id网页配置；
	// 处理充电信息
	if err := h.chargeService.ProcessChargingAndDiscount(parkinfo, request.PlateNumber); err != nil {
		logger.Logger.Error("处理充电和优惠失败", err, "parkID", parkinfo.ParkID, "plateNo", request.PlateNumber)
		h.MakeRepsonse(c, XjcResponse{
			Result:         0,
			ResultMsg:      "处理充电和优惠失败",
			DiscountNumber: request.DiscountNumber,
		})
		return
	}

	logger.Logger.Info("充电记录处理成功", "parkID", parkinfo.ParkID, "plateNo", request.PlateNumber)

	h.MakeRepsonse(c, XjcResponse{
		Result:         1,
		ResultMsg:      "success",
		DiscountNumber: request.DiscountNumber,
	})
}
