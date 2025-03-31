package api

import (
	"fmt"
	"html/template"
	"jilaidian_go/internal/config"
	"jilaidian_go/www"
	"net/http"
	"strconv"
	"strings"

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
	r.GET("/", h.indexHandler)
	r.POST("/save", h.saveHandler)
	r.POST("/add", h.addHandler)
	r.GET("/api/parks", func(c *gin.Context) {
		ci := config.LoadConfig()
		c.JSON(http.StatusOK, ci.Parks)
	})

	r.POST("/api/deletepark/{parkid}", func(c *gin.Context) {

		parkIDStr := strings.TrimPrefix(c.Request.URL.Path, "/api/deletepark/")
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
	apiBaseURL := r.Form.Get("api.baseUrl")
	var parks []config.ParkInfo
	for key, _ := range r.Form {
		if strings.HasPrefix(key, "parks.parkid.") {
			parkID, err := strconv.Atoi(strings.TrimPrefix(key, "parks.parkid."))

			if err != nil {
				//http.Error(w, "Invalid park ID", http.StatusBadRequest)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid park ID"})
				return
			}
			ukey := r.Form.Get(fmt.Sprintf("parks.%d.ukey", parkID))
			deductionTime, _ := strconv.Atoi(r.Form.Get(fmt.Sprintf("parks.%d.deduction_time", parkID)))
			deductionMoney, _ := strconv.Atoi(r.Form.Get(fmt.Sprintf("parks.%d.deduction_money", parkID)))
			newParkid, _ := strconv.Atoi(r.Form.Get(key))
			parks = append(parks, config.ParkInfo{
				ParkID:         newParkid,
				Ukey:           ukey,
				DeductionTime:  deductionTime,
				DeductionMoney: deductionMoney,
			})
		}
	}
	fmt.Println(parks)
	fmt.Println(len(parks))
	ci := config.LoadConfig()
	if serverPort != "" {
		ci.Server.Port = serverPort
	}
	if apiBaseURL != "" {
		ci.API.BaseURL = apiBaseURL
	}
	if len(parks) > 0 {
		ci.Parks = parks
	}

	if err := ci.SaveConfig(); err != nil {
		//http.Error(c.Writer, "Failed to save config", http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}

	///http.Redirect(c.Writer, r, "/", http.StatusSeeOther)
	c.Redirect(http.StatusSeeOther, "/")
}

func (h *ConfigHandler) addHandler(c *gin.Context) {
	r := c.Request
	r.ParseForm()

	ukey := r.Form.Get("add.ukey")
	deductionTime, _ := strconv.Atoi(r.Form.Get("add.deduction_time"))
	deductionMoney, _ := strconv.Atoi(r.Form.Get("add.deduction_money"))
	parkID, _ := strconv.Atoi(r.Form.Get("add.parkid"))

	parkinfo := config.ParkInfo{
		ParkID:         parkID,
		Ukey:           ukey,
		DeductionTime:  deductionTime,
		DeductionMoney: deductionMoney,
	}

	ci := config.LoadConfig()
	ci.Parks = append(ci.Parks, parkinfo)

	if err := ci.SaveConfig(); err != nil {
		//http.Error(c.Writer, "Failed to save config", http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}

	http.Redirect(c.Writer, r, "/", http.StatusSeeOther)
}
