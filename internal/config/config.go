package config

import (
	"encoding/json"
	"fmt"

	"jilaidian_go/pkg/common"
	"jilaidian_go/pkg/logger"
	"os"
	"path/filepath"
)

type ParkInfo struct {
	ParkID         int    `json:"parkid"`     // 改为驼峰式
	StationID      string `json:"station_id"` // 改为驼峰式
	Ukey           string `json:"ukey"`
	DeductionTime  int    `json:"deduction_time"`  // 改为驼峰式
	DeductionMoney int    `json:"deduction_money"` // 改为驼峰式
	ReduceAmount   int    `json:"reduceAmount"`
	Duration       int    `json:"Duration"`
	Remark         string `json:"remark"`
}
type YiAnqi struct {
	TokenURL       string `json:"tokenUrl" form:"YiAnqi.TokenURL"`
	OperatorID     string `json:"operatorID" form:"YiAnqi.OperatorID"`
	OperatorSecret string `json:"operatorSecret" form:"YiAnqi.OperatorSecret"`
	AesKey         string `json:"aeskey" form:"YiAnqi.AesKey"`
	AesIv          string `json:"aesiv" form:"YiAnqi.AesIv"`
	SignKey        string `json:"signKey" form:"YiAnqi.SignKey"`
}

// 添加一个方法用于获取解密后的ukey
func (p *ParkInfo) GetUkey() string {
	// 这里可以添加解密逻辑，或从环境变量获取
	// 简单示例：如果环境变量中有对应的ukey，则使用环境变量中的值
	envKey := os.Getenv(fmt.Sprintf("PARK_UKEY_%d", p.ParkID))
	if envKey != "" {
		return envKey
	}
	return p.Ukey
}

// Config 应用配置结构体
type Config struct {
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
	Parks []ParkInfo `json:"parks"`

	YiAnqi YiAnqi `json:"yianqi"`

	API struct {
		TingCheYunUrl string `json:"tingCheYunUrl"`
	} `json:"api"`
	Debug bool
}

// Global 全局配置实例
var Global *Config

func createDefaultConfig() *Config {
	config := &Config{}
	config.Server.Port = "8081"
	config.Parks = []ParkInfo{
		{
			ParkID:         99999,
			Ukey:           "your_ukey",
			DeductionTime:  61,
			DeductionMoney: 62,
			ReduceAmount:   101,
			Duration:       0,
			Remark:         "备注",
			StationID:      "1234567890", // 示例站点ID
		},
		{
			ParkID:         888888,
			Ukey:           "your_ukey",
			DeductionTime:  60,
			DeductionMoney: 60,
			ReduceAmount:   100,
			Duration:       0,
			Remark:         "备注",
			StationID:      "1234567890", // 示例站点ID
		},
	}
	config.YiAnqi = YiAnqi{
		AesIv:          "sExIBhW3Y4mYr0ne",
		AesKey:         "2RU7xUBhsyW5SnCY",
		OperatorID:     "cdist2025",
		OperatorSecret: "u8wLRTOTkjKmVpoP",
		SignKey:        "Q2YAboqxfmrtbfsw",
		TokenURL:       "https://api.ddpark.fun/api/yianqi/token",
	}

	config.API.TingCheYunUrl = "http://istparking.sciseetech.com/public"
	config.Debug = os.Getenv("DEBUG") == "1"
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
	// 使用os.WriteFile替代ioutil.WriteFile
	err = os.WriteFile(filename, data, 0644)
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
	err := c.WriteConfigToFile(common.GetConfigPath())
	if err != nil {
		return err
	}
	Global = c // 更新全局配置
	return nil
}
func GetParkInfo(stationID string) *ParkInfo {
	for i := range Global.Parks {
		if Global.Parks[i].StationID == stationID {
			// 返回数组元素的地址，而不是临时变量的地址
			return &Global.Parks[i]
		}
	}
	return nil
}
func LoadConfig() *Config {
	cc := &Config{}
	if err := cc.Load(); err != nil {
		logger.Logger.Error("加载配置失败", err)
		return nil
	}
	return cc
}

func init() {
	Global = &Config{}

	if err := Global.Load(); err != nil {
		logger.Logger.Error("初始化全局配置失败", err)
		// 不要重复记录相同的错误信息
		return
	}

	logger.Logger.Info("全局配置加载成功")
}
