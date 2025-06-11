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
		authorized.POST("/inpark", h.In_out_parkHandler)
		authorized.POST("/outpark", h.In_out_parkHandler)

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

func (h *TingCheyunHandler) In_out_parkHandler(c *gin.Context) {
	// 处理请求并返回响应

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		h.MakeRepsonse(c, 1, "In_out_parkHandler", "Failed to read request body")
		return
	}
	iob := In_out_Base{}
	if err := json.Unmarshal(body, &iob); err != nil {
		logger.Logger.Errorf("JSON解析失败 %s  %v", body, err)
		h.MakeRepsonse(c, 1, "In_out_parkHandler", "JSON解析失败")
		return
	}
	logger.Logger.Info("Received request", "body", string(body))
	switch iob.ServiceName {
	case "in_park":
		{
			inPark := InPark{}
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
			thirdservice.PostInData(inData)
		}
	case "out_park":
		{
			outPark := OutPark{}
			if err := json.Unmarshal(body, &outPark); err != nil {
				logger.Logger.Errorf("JSON解析失败 %s  %v", body, err)
				h.MakeRepsonse(c, 1, "out_park", "JSON解析失败")
				return
			}
			outData := thirdservice.ThirdOutObj{
				OutTime:     outPark.Data.OutTime * 1000,
				PlateNumber: outPark.Data.CarNumber,
				ImgUrl:      outPark.Data.PicAddr,
				ParkingId:   strconv.Itoa(outPark.ParkID),
			}
			thirdservice.PostOutData(outData)
		}
	}
	h.MakeRepsonse(c, 0, iob.ServiceName, "Success")
}
