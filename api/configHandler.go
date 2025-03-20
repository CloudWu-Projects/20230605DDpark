package api

import (
	"embed"
	"fmt"
	"html/template"
	"jilaidian_go/config"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler API处理器
type ConfigHandler struct {
	tmpl *template.Template
}

// NewHandler 创建新的API处理器
func NewConfigHandler(content embed.FS) *ConfigHandler {

	tmpl, err := template.ParseFS(content, "www/config.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		panic(err)
	}
	return &ConfigHandler{
		tmpl: tmpl,
	}
}

// SetupRoutes 设置路由
func (h *ConfigHandler) SetupRoutes(router *gin.Engine) {
	// 充电记录接口
	router.POST("/save", h.saveHandler)
	router.GET("/", h.indexHandler)
	router.POST("/add", h.addHandler)
}
func (h *ConfigHandler) HandleChargingRecord(c *gin.Context) {}

func (h *ConfigHandler) indexHandler(c *gin.Context) {
	ci := config.LoadConfig()
	h.tmpl.Execute(c.Writer, ci)
}
func (h *ConfigHandler) saveHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	serverPort := r.Form.Get("server.port")
	apiBaseURL := r.Form.Get("api.baseUrl")
	var parks []config.ParkInfo
	for key, _ := range r.Form {
		if strings.HasPrefix(key, "parks.parkid.") {
			parkID, err := strconv.Atoi(strings.TrimPrefix(key, "parks.parkid."))
			if err != nil {
				http.Error(w, "Invalid park ID", http.StatusBadRequest)
				return
			}
			ukey := r.Form.Get(fmt.Sprintf("parks.%d.ukey", parkID))
			deductionTime, _ := strconv.Atoi(r.Form.Get(fmt.Sprintf("parks.%d.deduction_time", parkID)))
			deductionMoney, _ := strconv.Atoi(r.Form.Get(fmt.Sprintf("parks.%d.deduction_money", parkID)))
			newParkid, _ := strconv.Atoi(r.Form.Get(key))
			parks = append(parks, config.ParkInfo{
				Parkid:          newParkid,
				Ukey:            ukey,
				Deduction_time:  deductionTime,
				Deduction_money: deductionMoney,
			})
		}
	}
	ci := &config.Config{
		Server: struct {
			Port string `json:"port"`
		}{Port: serverPort},
		Parks: parks,
		API: struct {
			BaseURL string `json:"baseUrl"`
		}{BaseURL: apiBaseURL},
	}

	if err := ci.SaveConfig(); err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *ConfigHandler) addHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request
	if r.Method != http.MethodPost {
		//http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		tmpl := template.Must(template.ParseFiles("add.html"))
		c := config.LoadConfig()
		tmpl.Execute(w, c)
		return
	}
	r.ParseForm()

	ukey := r.Form.Get("add.ukey")
	deductionTime, _ := strconv.Atoi(r.Form.Get("add.deduction_time"))
	deductionMoney, _ := strconv.Atoi(r.Form.Get("add.deduction_money"))
	parkID, _ := strconv.Atoi(r.Form.Get("add.parkid"))

	parkinfo := config.ParkInfo{
		Parkid:          parkID,
		Ukey:            ukey,
		Deduction_time:  deductionTime,
		Deduction_money: deductionMoney,
	}

	ci := config.LoadConfig()
	ci.Parks = append(ci.Parks, parkinfo)

	if err := ci.SaveConfig(); err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
