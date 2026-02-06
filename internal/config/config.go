package config

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"jilaidian_go/pkg/common"
	"jilaidian_go/pkg/logger"
	"jilaidian_go/version"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port          string `json:"port"`
	TingCheYunUrl string `json:"tingCheYunUrl"`
}

// Config 应用配置结构体
type Config struct {
	ServerConfig ServerConfig `json:"server"`
	Parks        []ParkInfo   `json:"parks"`

	YiAnqi YiAnqi `json:"yianqi"`

	Debug          bool           `json:"debug"`
	NeedQQLogin    bool           `json:"need_qq_login"`
	NanjingNengRui NanjingNengRui `json:"nanjingnengrui"`
	XinJunCheng    XinJunCheng    `json:"xinjuncheng"`
	HeiBeiProxy    HeiBeiProxy    `json:"heibieproxy"`
}

// Global 全局配置实例
var Global *Config
var GlobalVersion string

func init() {
	GlobalVersion = fmt.Sprintf("version:%s-%s build:%s go:%s", version.GitVersion, version.GitHash, version.BuildTime, version.GoVersion)
}

func createDefaultConfig() *Config {
	config := &Config{}
	config.ServerConfig = ServerConfig{
		Port:          "9090",
		TingCheYunUrl: "http://istparking.sciseetech.com/public",
	}

	config.Parks = []ParkInfo{
		{
			ParkID:         99999,
			Ukey:           "your_ukey",
			DeductionTime:  61,
			DeductionMoney: 62,
			ReduceAmount:   101,
			Duration:       0,
			Remark:         "备注",
			StationID:      "007B45C733038000", // 示例站点ID
			//ProxyUrl:       "http://127.0.0.1:28080/proxy",
			NpcPort: "33124",
		},
		{
			ParkID:         888888,
			Ukey:           "your_ukey",
			DeductionTime:  60,
			DeductionMoney: 60,
			ReduceAmount:   100,
			Duration:       0,
			Remark:         "备注",
			StationID:      "346790", // 示例站点ID
			//ProxyUrl:       "https://api.ddpark.fun:33124/api/yianqi/proxy",
			NpcPort: "33124",
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
	config.NanjingNengRui = NanjingNengRui{
		AppSercert: "I654HUNOU250AX45",
		AppId:      "cdist",
	}
	config.HeiBeiProxy = HeiBeiProxy{
		ProxyUrl: "https://api.ddpark.fun:{PORT}/v2/tripartite/queryEtcOrder",
	}
	config.XinJunCheng = XinJunCheng{
		VKey: "your_vkey_here",
	}

	config.Debug = os.Getenv("DEBUG") == "1"
	config.NeedQQLogin = os.Getenv("NEED_QQ_LOGIN") == "1"
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

func GetParkInfo_withParkid(parkID string) *ParkInfo {
	for i := range Global.Parks {
		if strconv.Itoa(Global.Parks[i].ParkID) == parkID {
			// 返回数组元素的地址，而不是临时变量的地址
			return &Global.Parks[i]
		}
	}
	return nil
}
func GetParkInfoByAppId(appId string) *ParkInfo {
	for i := range Global.Parks {
		if strconv.Itoa(Global.Parks[i].ParkID) == appId || Global.Parks[i].StationID == appId {
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
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}
	Global = &Config{}

	if err := Global.Load(); err != nil {
		logger.Logger.Error("初始化全局配置失败", err)
		// 不要重复记录相同的错误信息
		return
	}

	logger.Logger.Info("全局配置加载成功")
	logger.Logger.Info("全局配置:", Global)
	logger.Logger.Info("全局NeedQQLogin:", Global.NeedQQLogin)
}
