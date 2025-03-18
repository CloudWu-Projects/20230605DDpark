package config

import (
	"encoding/json"
	"os"
)

// Config 应用配置结构体
type Config struct {
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
	Park struct {
		AllowedIDs []string `json:"allowedIds"`
	} `json:"park"`
	API struct {
		BaseURL string `json:"baseUrl"`
		UKey    string `json:"uKey"`
	} `json:"api"`
}

// Global 全局配置实例
var Global Config

// Load 从文件加载配置
func Load(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &Global)
}

// GetAllowedParkIDs 获取允许的停车场ID列表
func GetAllowedParkIDs() []string {
	return Global.Park.AllowedIDs
}
