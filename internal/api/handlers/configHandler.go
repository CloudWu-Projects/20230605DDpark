package api

import (
	"fmt"
	"html/template"
	"jilaidian_go/internal/config"
	"jilaidian_go/internal/utils"
	"jilaidian_go/pkg/logger"
	"jilaidian_go/www"
	"net/http"
	"strconv"
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
	configG := r.Group("/config", SessionAuthMiddleware())
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
		h.tmplIndexHtml.Execute(c.Writer, config.Global)
	})

	r.GET("/www/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		c.FileFromFS("www/"+filepath, http.FS(www.HtmlFS))
	})
	r.GET("/log/:lastbytes", h.logHandler)
	r.GET("/log/", h.logHandler)

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
		// {
		// 	GroupLabel: "逸安启配置",
		// 	Url:        "/config/save",
		// 	Fields:     structToStringMap(ci.YiAnqi, hiddenValueYianqi),
		// },
		// {
		// 	GroupLabel: "南京能瑞配置",
		// 	Url:        "/config/save",
		// 	Fields:     structToStringMap(ci.NanjingNengRui, hiddenValueNanjingNengrui),
		// },
		{
			GroupLabel: "河北配置",
			Url:        "/config/save",
			Fields:     structToStringMap(ci.HeiBeiProxy, hiddenValueHeiBeiProxy),
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
	deductionTime, _ := strconv.Atoi(r.Form.Get("add.deduction_time"))
	deductionMoney, _ := strconv.Atoi(r.Form.Get("add.deduction_money"))
	parkID, _ := strconv.Atoi(r.Form.Get("add.parkid"))
	StationID := r.Form.Get("add.station_id")
	//ProxyUrl := r.Form.Get("add.proxy_url")
	NpcPort := r.Form.Get("add.npc_port")
	Duration, _ := strconv.Atoi(r.Form.Get("add.Duration"))
	Remark := r.Form.Get("add.remark")
	//	Deduction:   100,
	//	Remark:        "备注",
	parkinfo := config.ParkInfo{
		ParkID:         parkID,
		Ukey:           ukey,
		DeductionTime:  deductionTime,
		DeductionMoney: deductionMoney,
		StationID:      StationID,
		Duration:       Duration,
		Remark:         Remark,
		NpcPort:        NpcPort,
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
