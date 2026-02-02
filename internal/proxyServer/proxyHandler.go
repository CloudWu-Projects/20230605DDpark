package proxyserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jilaidian_go/internal/config"
	"jilaidian_go/internal/debugHandler"
	"log"
	"net/http"
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

	r.POST("/v2/heetc/payNotify", h.HandleProxy)

}

// APIResponse represents a generic API response
type APIResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
type Response struct {
	State    int    `json:"state"`
	ErrorMsg string `json:"errormsg"`
}

// SendErrorResponse sends an error response to the client
func SendErrorResponse(c *gin.Context, msg string) {
	response := APIResponse{Code: 1, Msg: msg}
	c.JSON(http.StatusOK, response)
}

// SendSuccessResponse sends a success response to the client
func SendSuccessResponse(c *gin.Context, msg string) {
	response := APIResponse{Code: 0, Msg: msg}
	c.JSON(http.StatusOK, response)
}

type Visitor struct {
	ParkID string `json:"park_id"`
	// 其他字段根据需要添加
}

func (v *Visitor) UnmarshalJSON(data []byte) error {
	// 定义一个临时 map 来解析 JSON
	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// 尝试获取 "park_id"
	if id, ok := temp["park_id"].(string); ok {
		v.ParkID = id
		return nil
	}

	// 如果没有 "park_id"，尝试获取 "appid"
	if id, ok := temp["appId"].(string); ok {
		v.ParkID = id
		return nil
	}

	// 都没有则返回错误
	return fmt.Errorf("json: missing park_id or appid field")
}
func (h *Handler) HandleProxy(c *gin.Context) {
	//w := c.Writer
	r := c.Request
	//check if r.URL.Path  in  Config.proxy

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		SendErrorResponse(c, "Failed to read request body")
		return
	}

	var visitor Visitor
	if err := visitor.UnmarshalJSON(bodyBytes); err != nil {
		SendErrorResponse(c, fmt.Sprintf("Failed to parse JSON: %v", err))
		return
	}

	if visitor.ParkID == "" {
		SendErrorResponse(c, "Missing park_id")
		return
	}
	parkInfo := config.GetParkInfoByAppId(visitor.ParkID)

	if parkInfo == nil {
		SendErrorResponse(c, fmt.Sprintf("parkInfo [%s] not found", visitor.ParkID))
		return
	}
	fmt.Println(parkInfo)
	targetURL := parkInfo.ProxyUrl

	log.Printf("Forwarding visitor[%s] request to: %s", visitor.ParkID, targetURL)

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
	//c.JSON(http.StatusOK, string(body))
}
