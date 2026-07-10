package api

import (
	"fmt"
	"html/template"
	"jilaidian_go/internal/config"
	"jilaidian_go/internal/service"
	"jilaidian_go/internal/utils"
	"jilaidian_go/pkg/logger"
	"jilaidian_go/www"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	configHandlerInstance *ConfigHandler
	configHandlerOnce     sync.Once
)

// Handler API处理器
type ConfigHandler struct {
	tmplConfigHtml *template.Template
	tmplIndexHtml  *template.Template
	chargeService  *service.ChargeService
}

type QueryPageData struct {
	Version string
	Plate   string
	Result  *QueryResultView
	Error   string
}

type QueryResultView struct {
	Found         bool
	ParkName      string
	Plate         string
	OrderID       string
	ParkedMinutes int64
	InTimeText    string
	Message       string
}

// NewHandler 创建新的API处理器
func NewConfigHandler(r *gin.Engine) *ConfigHandler {
	configHandlerOnce.Do(func() {

		tmplConfigHtml, err := template.ParseFS(www.HtmlFS, "www/config.html", "www/modalForm.html")
		if err != nil {
			fmt.Println("Error parsing template:", err)
			panic(err)
		}
		tmplIndexHtml, err := template.ParseFS(www.HtmlFS, "www/index.html")
		configHandlerInstance = &ConfigHandler{
			tmplConfigHtml: tmplConfigHtml,
			tmplIndexHtml:  tmplIndexHtml,
			chargeService:  service.NewChargeService(),
		}
		configHandlerInstance.setupRoutes(r)
	})
	return configHandlerInstance
}

// 认证中间件
func CheckSessionValid(c *gin.Context) bool {
	if !config.Global.NeedQQLogin {
		return true
	}
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		return false
	}
	userSession, ok := user.(UserSession)
	if !ok || !userSession.IsLogin {
		return false
	}
	if time.Since(userSession.LastActive) > MaxSessionDuration {
		return false
	}
	fmt.Println("Session valid")
	return true
}

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !CheckSessionValid(c) {
			scheme := "http"
			if c.Request.TLS != nil {
				scheme = "https"
			}
			redirectURL := fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, c.Request.RequestURI)
			if c.Request.Method != http.MethodGet {
				redirectURL = fmt.Sprintf("%s://%s%s", scheme, c.Request.Host)
			}
			c.Redirect(http.StatusFound, "/toLogin?redirect_url="+redirectURL)
			c.Abort()
			return
		}
		c.Next()
	}
}

// SetupRoutes 设置路由
func (h *ConfigHandler) setupRoutes(r *gin.Engine) {
	// 充电记录接口
	configG := r.Group("/config")
	{
		configG.GET("/c", func(c *gin.Context) {
			c.JSON(http.StatusOK, config.Global)
		})
		configG.GET("/", h.indexHandler)
		configG.POST("/baseserver", h.saveSeverConfigHandler)
		configG.POST("/save", h.saveConfigHandler)
		configG.POST("/add", h.addHandler)
		configG.DELETE("/park/:parkid", h.deleteParkHandler)
		configG.GET("/parks", func(c *gin.Context) {
			ci := config.LoadConfig()
			c.JSON(http.StatusOK, ci.Parks)
		})
		// Add endpoint to list all routes
		configG.GET("/r", func(c *gin.Context) {
			routeshtml, _ := template.ParseFS(www.HtmlFS, "www/routes.html")
			routeshtml.Execute(c.Writer, r.Routes())
		})
	}

	r.GET("/", func(c *gin.Context) {
		h.renderQueryPage(c, QueryPageData{Version: config.GlobalVersion})
	})
	r.POST("/query", h.queryOrderHandler)

	r.GET("/www/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		c.FileFromFS("www/"+filepath, http.FS(www.HtmlFS))
	})
	r.GET("/log/:lastbytes", h.logHandler)
	r.GET("/log/", h.logHandler)

}

func (h *ConfigHandler) renderQueryPage(c *gin.Context, data QueryPageData) {
	if err := h.tmplIndexHtml.Execute(c.Writer, data); err != nil {
		logger.Logger.Errorf("渲染查询页面失败: %v", err)
		c.String(http.StatusInternalServerError, "render query page failed")
	}
}

func (h *ConfigHandler) queryOrderHandler(c *gin.Context) {
	plate := strings.ToUpper(strings.TrimSpace(c.PostForm("plate")))
	pageData := QueryPageData{
		Version: config.GlobalVersion,
		Plate:   plate,
	}

	result, err := h.chargeService.QueryFirstMatchingOrder(plate)
	if err != nil {
		pageData.Error = err.Error()
		h.renderQueryPage(c, pageData)
		return
	}
	if result == nil {
		pageData.Result = &QueryResultView{
			Found:   false,
			Plate:   plate,
			Message: "未查询到当前在场订单",
		}
		h.renderQueryPage(c, pageData)
		return
	}

	parkName := strings.TrimSpace(result.ParkInfo.Remark)
	if parkName == "" {
		parkName = fmt.Sprintf("车场 %d", result.ParkInfo.ParkID)
	}

	parkedMinutes := int64(0)
	if result.InTime > 0 {
		parkedMinutes = time.Now().Unix() - result.InTime
		if parkedMinutes < 0 {
			parkedMinutes = 0
		}
		parkedMinutes /= 60
	}

	pageData.Result = &QueryResultView{
		Found:         true,
		ParkName:      parkName,
		Plate:         plate,
		OrderID:       result.OrderID,
		ParkedMinutes: parkedMinutes,
		InTimeText:    time.Unix(result.InTime, 0).Format("2006-01-02 15:04:05"),
	}
	h.renderQueryPage(c, pageData)
}

func (h *ConfigHandler) deleteParkHandler(c *gin.Context) {
	parkIDStr := c.Param("parkid")
	logger.Logger.Error("deletepark parkIDStr:", parkIDStr)
	parkID, err := strconv.Atoi(parkIDStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid park ID"})
		return
	}

	fmt.Printf("\ndeletepark parkID: %d \n", parkID)
	ci := config.LoadConfig()

	for i, park := range ci.Parks {
		if park.ParkID == parkID {
			ci.Parks = append(ci.Parks[:i], ci.Parks[i+1:]...)
			break
		}
	}
	if err := ci.SaveConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Park deleted successfully"})
}

func (h *ConfigHandler) indexHandler(c *gin.Context) {
	ci := config.LoadConfig()

	groups := []FieldGroup{
		{
			GroupLabel: "系统配置",
			Url:        "/config/save",
			Fields:     structToStringMap(ci.ServerConfig, hiddenValueBaseServer),
		},
	}
	data := struct {
		Groups  []FieldGroup
		Debug   bool
		Version string
	}{
		Groups:  groups,
		Debug:   ci.Debug,
		Version: config.GlobalVersion,
	}

	h.tmplConfigHtml.Execute(c.Writer, data)
}

func (h *ConfigHandler) logHandler(c *gin.Context) {
	lastBytesStr := c.Param("lastbytes")
	lastBytes, _ := strconv.ParseInt(lastBytesStr, 10, 64)
	logContent, err := utils.ReadLogFile(lastBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read log file"})
		return
	}

	c.Writer.WriteString(logContent)
}
func (h *ConfigHandler) addHandler(c *gin.Context) {
	r := c.Request
	r.ParseForm()

	ukey := r.Form.Get("add.ukey")
	oldparkid := r.Form.Get("add.old_parkid")
	parkID, _ := strconv.Atoi(r.Form.Get("add.parkid"))
	Remark := r.Form.Get("add.remark")
	parkinfo := config.ParkInfo{
		ParkID: parkID,
		Ukey:   ukey,
		Remark: Remark,
	}
	fmt.Println("addHandler parkinfo:", oldparkid, parkinfo)
	ci := config.LoadConfig()
	if oldparkid == "" {
		ci.Parks = append(ci.Parks, parkinfo)
	} else {
		oldparkidInt, _ := strconv.Atoi(oldparkid)
		for i, park := range ci.Parks {
			if park.ParkID == oldparkidInt {
				ci.Parks[i] = parkinfo
				break
			}
		}
	}

	if err := ci.SaveConfig(); err != nil {
		//http.Error(c.Writer, "Failed to save config", http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/config/")
}
