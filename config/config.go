package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jilaidian_go/logger"
	"os"
	"path/filepath"
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
var Global *Config

func createDefaultConfig() *Config {
	config := &Config{}
	config.Server.Port = "8080"
	config.Park.AllowedIDs = []string{""}
	config.API.BaseURL = "http://istparking.sciseetech.com/public"
	config.API.UKey = "7JPWIA1SGV9N17LE"
	return config
}
func createDirIfNotExist(filename string) error {
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 创建目录
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		fmt.Printf("Directory created: %s\n", dir)

	}
	return nil
}

// WriteConfigToFile writes the config to a file in YAML format
func WriteConfigToFile(config *Config, filename string) error {
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	createDirIfNotExist(filename)
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Load 从文件加载配置
func Load(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Logger.Error("没找到配置文件 重新创建")
			Global := createDefaultConfig()
			WriteConfigToFile(Global, filePath)
			return nil
		}
		return err
	}
	Global = createDefaultConfig()
	return json.Unmarshal(data, Global)
}

// GetAllowedParkIDs 获取允许的停车场ID列表
func GetAllowedParkIDs() []string {
	return Global.Park.AllowedIDs
}

func init() {

	if err := Load("config/config.json"); err != nil {
		logger.Logger.Error("加载配置失败", err)
		return
	}
}
