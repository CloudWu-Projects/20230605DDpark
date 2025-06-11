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

	"github.com/gin-gonic/gin"
)

// Handler API处理器
type ConfigHandler struct {
	tmplConfigHtml *template.Template
}

// NewHandler 创建新的API处理器
func NewConfigHandler() *ConfigHandler {
	tmplConfigHtml, err := template.ParseFS(www.HtmlFS, "config.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		panic(err)
	}
	return &ConfigHandler{
		tmplConfigHtml: tmplConfigHtml,
	}
}

// SetupRoutes 设置路由
func (h *ConfigHandler) SetupRoutes(r *gin.Engine) {
	// 充电记录接口
	configG := r.Group("/config")
	{
		configG.GET("/", h.indexHandler)
	}
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, config.Global)
	})
	r.POST("/save", h.saveHandler)
	r.POST("/add", h.addHandler)
	r.GET("/log", h.logHandler)
	r.GET("/api/parks", func(c *gin.Context) {
		ci := config.LoadConfig()
		c.JSON(http.StatusOK, ci.Parks)
	})
	r.DELETE("/api/park/:parkid", func(c *gin.Context) {
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
	})
}

func (h *ConfigHandler) indexHandler(c *gin.Context) {
	ci := config.LoadConfig()
	h.tmplConfigHtml.Execute(c.Writer, ci)
}
func (h *ConfigHandler) saveHandler(c *gin.Context) {

	r := c.Request
	r.ParseForm()

	serverPort := r.Form.Get("server.port")
	outUrl := r.Form.Get("ThridServer.Out_url")
	inUrl := r.Form.Get("ThridServer.In_url")

	var parks []config.ParkInfo
	// for key, _ := range r.Form {
	// 	if strings.HasPrefix(key, "parks.parkid.") {
	// 		parkID, err := strconv.Atoi(strings.TrimPrefix(key, "parks.parkid."))

	// 		if err != nil {
	// 			//http.Error(w, "Invalid park ID", http.StatusBadRequest)
	// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid park ID"})
	// 			return
	// 		}
	// 		ukey := r.Form.Get(fmt.Sprintf("parks.%d.ukey", parkID))
	// 		deductionTime, _ := strconv.Atoi(r.Form.Get(fmt.Sprintf("parks.%d.deduction_time", parkID)))
	// 		deductionMoney, _ := strconv.Atoi(r.Form.Get(fmt.Sprintf("parks.%d.deduction_money", parkID)))
	// 		newParkid, _ := strconv.Atoi(r.Form.Get(key))

	// 		StationID := r.Form.Get(fmt.Sprintf("parks.%d.station_id", parkID))
	// 		parks = append(parks, config.ParkInfo{
	// 			ParkID:         newParkid,
	// 			Ukey:           ukey,
	// 			DeductionTime:  deductionTime,
	// 			DeductionMoney: deductionMoney,
	// 			StationID:      StationID,
	// 		})
	// 	}
	// }
	fmt.Println(parks)
	fmt.Println(len(parks))
	ci := config.LoadConfig()
	if serverPort != "" {
		ci.Server.Port = serverPort
	}
	if outUrl != "" {
		ci.ThridServer.Out_url = outUrl
	}
	if inUrl != "" {
		ci.ThridServer.In_url = inUrl
	}

	if len(parks) > 0 {
		ci.Parks = parks
	}

	if err := ci.SaveConfig(); err != nil {
		//http.Error(c.Writer, "Failed to save config", http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/config/")
}

func (h *ConfigHandler) logHandler(c *gin.Context) {
	// Read log file content using the utils package
	logContent, err := utils.ReadLogFile()
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
