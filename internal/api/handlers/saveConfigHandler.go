package api

import (
	"encoding/json"
	"fmt"
	"io"
	"jilaidian_go/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	hiddenValueBaseServer     = "baseserver"
	hiddenValueYianqi         = "yianqi"
	hiddenValueNanjingNengrui = "nanjingnengrui"
)

type HiddenObject struct {
	HiddenValue string `json:"hiddenValue"`
}

func (h *ConfigHandler) saveSeverConfigHandler(c *gin.Context) {

	var serverConfig config.ServerConfig
	if err := c.ShouldBindJSON(&serverConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ci := config.LoadConfig()
	if serverConfig.Port != "" {
		ci.ServerConfig.Port = serverConfig.Port

	}
	if serverConfig.TingCheYunUrl != "" {
		ci.ServerConfig.TingCheYunUrl = serverConfig.TingCheYunUrl
	}
	if err := ci.SaveConfig(); err != nil {
		//http.Error(c.Writer, "Failed to save config", http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Config saved successfully"})
}
func (h *ConfigHandler) saveConfigHandler(c *gin.Context) {
	// c.json have a item  "hiddenValue" which is not used in the code, so we can remove it
	// var hiddenValue string
	// we should follow hiddenvalue to select the correct config to save
	//need debug  to see the post value
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}

	var hiddenValue HiddenObject
	err = json.Unmarshal(body, &hiddenValue)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ci := config.LoadConfig()
	var unmarshalErr error
	switch hiddenValue.HiddenValue {
	case hiddenValueBaseServer:
		var serverConfig config.ServerConfig
		unmarshalErr = json.Unmarshal(body, &serverConfig)
		if unmarshalErr == nil {
			if serverConfig.Port != "" {
				ci.ServerConfig.Port = serverConfig.Port

			}
			if serverConfig.TingCheYunUrl != "" {
				ci.ServerConfig.TingCheYunUrl = serverConfig.TingCheYunUrl
			}
		}
	case hiddenValueYianqi:
		err = json.Unmarshal(body, &ci.YiAnqi)
	case hiddenValueNanjingNengrui:
		err = json.Unmarshal(body, &ci.NanjingNengRui)
	}
	if unmarshalErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"aaaa error": err.Error()})
		return
	}
	if err := ci.SaveConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Config saved successfully"})

	//c.Redirect(http.StatusSeeOther, "/config/")
}
