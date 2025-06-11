package tingcheyunservice

import (
	"encoding/json"
	"io"
	thirdservice "jilaidian_go/internal/thirdService"
	"jilaidian_go/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler API处理器
type TingCheyunHandler struct {
}

// NewHandler 创建新的API处理器
func NewHandler() *TingCheyunHandler {
	return &TingCheyunHandler{}
}

// SetupRoutes 设置路由
func (h *TingCheyunHandler) SetupRoutes(r *gin.Engine) {
	// 充电记录接口
	authorized := r.Group("/tingcheyun")
	//authorized.Use(h.AuthMiddleware()) // 使用中间件进行身份验证
	{
		authorized.POST("/inpark", h.InparkHandler)
		authorized.POST("/outpark", h.OutparkHandler)

	}

}

func (h *TingCheyunHandler) MakeRepsonse(c *gin.Context, result int, ServiceName string, description string) {

	tr := InParkResponse{
		State:       result,
		Errmsg:      description,
		ServiceName: ServiceName,
	}
	c.JSON(http.StatusOK, tr)
}

func (h *TingCheyunHandler) InparkHandler(c *gin.Context) {
	// 处理请求并返回响应
	inPark := InPark{}

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		h.MakeRepsonse(c, 1, "InParkHandler", "Failed to read request body")
		return
	}
	logger.Logger.Info("InParkHandler  Received:", string(body))

	if err := json.Unmarshal(body, &inPark); err != nil {
		logger.Logger.Errorf("JSON解析失败 %s  %v", body, err)
		h.MakeRepsonse(c, 1, "in_park", "JSON解析失败")
		return
	}
	inData := thirdservice.ThirdInObj{
		InTime:      inPark.Data.InTime,
		PlateNumber: inPark.Data.CarNumber,
		ImgUrl:      inPark.Data.PicAddr,
		ParkingId:   strconv.Itoa(inPark.ParkID),
	}
	err = thirdservice.PostInData(inData)
	if err != nil {
		logger.Logger.Errorf("in_park  发送失败 %v  %v", inData, err)
		h.MakeRepsonse(c, 1, "in_park", "发送失败")
		return
	}
	h.MakeRepsonse(c, 0, inPark.ServiceName, "Success")
}

func (h *TingCheyunHandler) OutparkHandler(c *gin.Context) {
	// 处理请求并返回响应

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		h.MakeRepsonse(c, 1, "OutparkHandler", "Failed to read request body")
		return
	}

	logger.Logger.Info("OutparkHandler  Received:", string(body))

	outPark := OutPark{}
	if err := json.Unmarshal(body, &outPark); err != nil {
		logger.Logger.Errorf("out_park JSON解析失败 %s  %v", body, err)
		h.MakeRepsonse(c, 1, "out_park", "JSON解析失败")
		return
	}
	outData := thirdservice.ThirdOutObj{
		OutTime:     outPark.Data.OutTime * 1000,
		PlateNumber: outPark.Data.CarNumber,
		ImgUrl:      outPark.Data.PicAddr,
		ParkingId:   strconv.Itoa(outPark.ParkID),
	}
	err = thirdservice.PostOutData(outData)
	if err != nil {
		logger.Logger.Errorf("out_park  发送失败 %v  %v", outData, err)
		h.MakeRepsonse(c, 1, "out_park", "发送失败")
		return
	}
	h.MakeRepsonse(c, 0, outPark.ServiceName, "Success")
}
