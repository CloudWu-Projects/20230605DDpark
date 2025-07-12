package api

import (
	"jilaidian_go/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ConfigHandler) saveYianqiConfigHandler(c *gin.Context) {
	var yianqiConfig config.YiAnqi
	if err := c.ShouldBindJSON(&yianqiConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"aaaa error": err.Error()})
		return
	}
	ci := config.LoadConfig()
	ci.YiAnqi = yianqiConfig
	if err := ci.SaveConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Config saved successfully"})

	//c.Redirect(http.StatusSeeOther, "/config/")
}
