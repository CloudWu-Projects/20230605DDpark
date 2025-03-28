package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jilaidian_go/lib/common"
	"jilaidian_go/lib/logger"
	"os"
	"path/filepath"
	"strconv"
)

type ParkInfo struct {
	Parkid          int    `json:"parkid"`
	Ukey            string `json:"ukey"`
	Deduction_time  int    `json:"deduction_time"`
	Deduction_money int    `json:"deduction_money"`
	ReduceAmount    int    `json:"reduceAmount"`
}

// Config 应用配置结构体
type Config struct {
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
	Parks []ParkInfo `json:"parks"`

	API struct {
		BaseURL string `json:"baseUrl"`
	} `json:"api"`
}

// Global 全局配置实例
var Global *Config

func createDefaultConfig() *Config {
	config := &Config{}
	config.Server.Port = "8080"
	config.Parks = []ParkInfo{
		{
			Parkid:          99999,
			Ukey:            "your_ukey",
			Deduction_time:  60,
			Deduction_money: 60,
			ReduceAmount:    100,
		},
		{
			Parkid:          888888,
			Ukey:            "your_ukey",
			Deduction_time:  60,
			Deduction_money: 60,
			ReduceAmount:    100,
		},
	}
	config.API.BaseURL = "http://istparking.sciseetech.com/public"
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
func (config *Config) WriteConfigToFile(filename string) error {
	data, err := json.MarshalIndent(config, "", "  ")
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
func (c *Config) Load() error {
	filePath := common.GetConfigPath()
	// 读取配置文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// 创建默认配置并写入文件
			defaultConfig := createDefaultConfig()
			err := defaultConfig.WriteConfigToFile(filePath)
			if err != nil {
				logger.Logger.Errorf("创建配置文件失败: %v", err)
				return err
			}
			// 将默认配置赋值给当前对象
			*c = *defaultConfig
			return nil
		}
		logger.Logger.Errorf("读取配置文件失败: %v", err)
		return err
	}

	// 解析配置文件内容
	if err := json.Unmarshal(data, c); err != nil {
		logger.Logger.Errorf("解析配置文件失败: %v", err)
		return err
	}
	return nil
}
func (c *Config) SaveConfig() error {
	return c.WriteConfigToFile(common.GetConfigPath())
}
func GetParkInfo(parkID string) *ParkInfo {
	for _, id := range Global.Parks {
		if strconv.Itoa(id.Parkid) == parkID {
			return &id
		}
	}
	return nil
}
func LoadConfig() *Config {
	cc := &Config{}
	if err := cc.Load(); err != nil {
		logger.Logger.Error("加载配置失败xxxx", err)
		return nil
	}
	return cc
}
func init() {
	Global = &Config{}

	if err := Global.Load(); err != nil {
		logger.Logger.Error("加载配置失败xxxxxaaaa", err)
		return
	}
}
