package shaaXiETC

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jilaidian_go/internal/config"
	"jilaidian_go/internal/debugHandler"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// Handler API处理器
type Handler struct {
	//chargeService *service.ChargeService
	DebugInfo *debugHandler.DebugInfo
}

var (
	instance *Handler
	once     sync.Once
)

// NewHandler 创建新的API处理器
func NewHandler(r *gin.Engine) *Handler {
	once.Do(func() {
		instance = &Handler{
			//	chargeService: service.NewChargeService(),
		}
		instance.SetupRoutes(r)
		instance.DebugInfo = debugHandler.GetDebugInfo("ProxyServer")
	})
	return instance
}

func (h *Handler) SetupRoutes(r *gin.Engine) {

	r.POST("/etcThirdPark/billConfirmV1", h.Handle214_billConfirmV1)
	r.POST("/etcThirdPark/billsDownV1", h.Handle215_billsDownV1)

}

// APIResponse represents a generic API response
type APIResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// SendErrorResponse sends an error response to the client
func SendErrorResponse(c *gin.Context, msg string) {
	response := APIResponse{Code: 1, Msg: msg}
	c.JSON(http.StatusOK, response)
}

// SendSuccessResponse sends a success response to the client
func SendSuccessResponse(c *gin.Context, msg string) {
	response := APIResponse{Code: 200, Msg: msg}
	c.JSON(http.StatusOK, response)
}

type RequestItem214 struct {
	ParkID      string `json:"parkId"`
	SplitDate   string `json:"splitDate"`
	PayNum      int32  `json:"payNum"`
	TotalToll   int    `json:"totalToll"`
	BillsDetail string `json:"billsDetail"`
}

func (v *RequestItem214) UnmarshalJSON(data []byte) error {
	// 定义一个临时 map 来解析 JSON
	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// 尝试获取 "parkId"
	if id, ok := temp["parkId"].(string); ok {
		v.ParkID = id
		return nil
	}

	// 都没有则返回错误
	return fmt.Errorf("json: missing parkId field")
}
func (h *Handler) Handle214_billConfirmV1(c *gin.Context) {
	//w := c.Writer
	r := c.Request
	//check if r.URL.Path  in  Config.proxy

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		SendErrorResponse(c, "Failed to read request body")
		return
	}

	var requestItem RequestItem214
	if err := requestItem.UnmarshalJSON(bodyBytes); err != nil {
		SendErrorResponse(c, fmt.Sprintf("Failed to parse JSON: %v", err))
		return
	}

	if requestItem.ParkID == "" {
		SendErrorResponse(c, "Missing parkId")
		return
	}
	parkInfo := config.GetParkInfoByAppId(requestItem.ParkID)

	if parkInfo == nil {
		SendErrorResponse(c, fmt.Sprintf("parkInfo [%s] not found", requestItem.ParkID))
		return
	}
	fmt.Println(parkInfo)
	//https://api.ddpark.fun:{PORT}/api/yianqi/proxy
	proxyUrlTemplate := config.Global.ShaaXiEtcProxy.ProxyUrl214
	//replace {PORT} with parkInfo.NpcPort
	// 使用strings.Replace替换{PORT}占位符
	targetURL := strings.Replace(proxyUrlTemplate, "{PORT}", parkInfo.NpcPort, 1)

	log.Printf("Forwarding 214 [%s] request to: %s", requestItem.ParkID, targetURL)

	payload := bytes.NewReader(bodyBytes)
	req, err := http.NewRequest("POST", targetURL, payload)
	if err != nil {
		SendErrorResponse(c, fmt.Sprintf("Failed to create request: %v", err))
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		SendErrorResponse(c, fmt.Sprintf("Failed to send request: %v", err))
		return
	}
	defer res.Body.Close()
	//io.Copy(c.Writer, res.Body)

	bodyBytes, err = io.ReadAll(res.Body)
	if err != nil {
		SendErrorResponse(c, fmt.Sprintf("Failed to read response body: %v", err))
		return
	}
	log.Printf("Response from %s: %s", targetURL, string(bodyBytes))
	c.Writer.WriteHeader(res.StatusCode)
	c.Writer.Write(bodyBytes)
}

func (h *Handler) Handle215_billsDownV1(c *gin.Context) {
	//w := c.Writer
	r := c.Request
	//check if r.URL.Path  in  Config.proxy

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		SendErrorResponse(c, "Failed to read request body")
		return
	}

	var requestItem RequestItem214
	if err := requestItem.UnmarshalJSON(bodyBytes); err != nil {
		SendErrorResponse(c, fmt.Sprintf("Failed to parse JSON: %v", err))
		return
	}

	if requestItem.ParkID == "" {
		SendErrorResponse(c, "Missing parkId")
		return
	}
	parkInfo := config.GetParkInfoByAppId(requestItem.ParkID)

	if parkInfo == nil {
		SendErrorResponse(c, fmt.Sprintf("parkInfo [%s] not found", requestItem.ParkID))
		return
	}
	fmt.Println(parkInfo)
	//https://api.ddpark.fun:{PORT}/api/yianqi/proxy
	proxyUrlTemplate := config.Global.ShaaXiEtcProxy.ProxyUrl215
	//replace {PORT} with parkInfo.NpcPort
	// 使用strings.Replace替换{PORT}占位符
	targetURL := strings.Replace(proxyUrlTemplate, "{PORT}", parkInfo.NpcPort, 1)

	log.Printf("Forwarding 215 [%s] request to: %s", requestItem.ParkID, targetURL)

	payload := bytes.NewReader(bodyBytes)
	req, err := http.NewRequest("POST", targetURL, payload)
	if err != nil {
		SendErrorResponse(c, fmt.Sprintf("Failed to create request: %v", err))
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		SendErrorResponse(c, fmt.Sprintf("Failed to send request: %v", err))
		return
	}
	defer res.Body.Close()
	//io.Copy(c.Writer, res.Body)

	bodyBytes, err = io.ReadAll(res.Body)
	if err != nil {
		SendErrorResponse(c, fmt.Sprintf("Failed to read response body: %v", err))
		return
	}
	log.Printf("Response from %s: %s", targetURL, string(bodyBytes))
	c.Writer.WriteHeader(res.StatusCode)
	c.Writer.Write(bodyBytes)
}
