package main

import (
	"embed"
	"fmt"
	"html/template"
	"jilaidian_go/config"
	"net/http"
	"strconv"
	"strings"
)

//go:embed config.html
var content embed.FS

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/add", addHandler)
	fmt.Println("Server started at :8081")
	http.ListenAndServe(":8081", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFS(content, "config.html"))
	c := config.LoadConfig()
	tmpl.Execute(w, c)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
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

			parks = append(parks, config.ParkInfo{
				Parkid:          parkID,
				Ukey:            ukey,
				Deduction_time:  deductionTime,
				Deduction_money: deductionMoney,
			})
		}
	}
	config := &config.Config{
		Server: struct {
			Port string `json:"port"`
		}{Port: serverPort},
		Parks: parks,
		API: struct {
			BaseURL string `json:"baseUrl"`
		}{BaseURL: apiBaseURL},
	}

	if err := config.WriteConfigToFile("config/config.json"); err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func addHandler(w http.ResponseWriter, r *http.Request) {

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

	c := config.LoadConfig()
	c.Parks = append(c.Parks, parkinfo)

	if err := c.SaveConfig(); err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
