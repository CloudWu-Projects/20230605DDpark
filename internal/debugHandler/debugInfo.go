package debugHandler

import (
	"sync"

	"github.com/gin-gonic/gin"
)

type DebugInfo struct {
	LastReuqestJson  interface{} `json:"last_request_json"`
	LastResponseJson interface{} `json:"last_response_json"`
}

type DebugInfoHandler struct {
	DebugInfo map[string]*DebugInfo `json:"debug_info"`
}

var (
	instance = DebugInfoHandler{DebugInfo: make(map[string]*DebugInfo)}
	once     sync.Once
)

// NewHandler 创建新的API处理器
func NewHandler(r *gin.Engine) DebugInfoHandler {
	once.Do(func() {

		instance.SetupRoutes(r)
	})
	return instance
}
func (h *DebugInfoHandler) SetupRoutes(r *gin.Engine) {
	r.GET("/debug/info/:serviceName", func(c *gin.Context) {
		serviceName := c.Param("serviceName")
		if info, exists := h.DebugInfo[serviceName]; exists {
			c.JSON(200, info)
		} else {
			c.JSON(404, gin.H{"error": "Debug info not found for service: " + serviceName})
		}
	})
	r.GET("/debug/all", func(c *gin.Context) {
		c.JSON(200, h.DebugInfo)
	})
}

func GetDebugInfo(serviceName string) *DebugInfo {
	if info, exists := instance.DebugInfo[serviceName]; exists {
		return info
	}
	info := &DebugInfo{
		LastReuqestJson:  "",
		LastResponseJson: "",
	}
	instance.DebugInfo[serviceName] = info
	return info
}
