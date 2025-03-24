package api

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"jilaidian_go/lib/config"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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
func (h *ConfigHandler) SetupRoutes(r *mux.Router) {
	// 充电记录接口
	r.HandleFunc("/", h.indexHandler).Methods("GET")
	r.HandleFunc("/save", h.saveHandler).Methods("POST")
	r.HandleFunc("/add", h.addHandler).Methods("POST")
	r.HandleFunc("/api/parks", func(w http.ResponseWriter, r *http.Request) {
		ci := config.LoadConfig()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ci.Parks)
	}).Methods("GET")

	r.HandleFunc("/api/deletepark/{parkid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		parkID, err := strconv.Atoi(vars["parkid"])
		if err != nil {
			http.Error(w, "Invalid park ID", http.StatusBadRequest)
			return
		}

		fmt.Printf("\ndeletepark parkID: %d \n", parkID)
		ci := config.LoadConfig()

		for i, park := range ci.Parks {
			if park.Parkid == parkID {
				ci.Parks = append(ci.Parks[:i], ci.Parks[i+1:]...)
				break
			}
		}
		if err := ci.SaveConfig(); err != nil {
			http.Error(w, "Failed to save config", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}).Methods("POST")
}

func (h *ConfigHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	ci := config.LoadConfig()
	h.tmpl.Execute(w, ci)
}
func (h *ConfigHandler) saveHandler(w http.ResponseWriter, r *http.Request) {

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
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *ConfigHandler) addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
