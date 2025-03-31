package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type QiuConfig struct {
	AccessKey  string `json:"accessKey"`
	SecretKey  string `json:"secretKey"`
	BucketPath string `json:"bucketPath"`
	TargetPath string `json:"targetPath"`
	BaseURL    string `json:"baseUrl"`
}

// getExecutableDir returns directory path of current executable
func getExecutableDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	return filepath.Dir(exePath), nil
}

// WriteConfigToFile writes the config to a file in YAML format
func (config *QiuConfig) WriteConfigToFile(filename string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	//createDirIfNotExist(filename)
	// 使用os.WriteFile替代ioutil.WriteFile
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
func (c *QiuConfig) LoadConfig() error {
	exeDir, err := getExecutableDir()
	if err != nil {
		return fmt.Errorf("get executable dir failed: %w", err)
	}
	filePath := filepath.Join(exeDir, "config.json")
	// 创建默认配置如果不存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		defaultConfig := &QiuConfig{
			AccessKey:  "YOUR_ACCESS_KEY",
			SecretKey:  "YOUR_SECRET_KEY",
			BucketPath: "default-bucket",
			TargetPath: "",
			BaseURL:    "http://7niu.hyman.store/",
		}
		if err := defaultConfig.WriteConfigToFile(filePath); err != nil {
			return fmt.Errorf("创建默认配置文件失败: %w", err)
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开配置文件失败: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		return err
	}
	return nil
}

func prepareURL(baseUrl string, config *QiuConfig) string {
	expire := time.Now().Unix() + 3600
	u, _ := url.Parse(baseUrl)
	query := u.Query()
	query.Set("e", fmt.Sprintf("%d", expire))
	u.RawQuery = query.Encode()

	accessKey := config.AccessKey
	secretKey := config.SecretKey

	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(u.String()))
	sign := mac.Sum(nil)
	encodedSign := base64.URLEncoding.EncodeToString(sign)

	token := fmt.Sprintf("%s:%s", accessKey, encodedSign)
	u.RawQuery += "&token=" + token

	return u.String()
}

func download(urlStr, targetPath string) error {
	fmt.Printf("正在下载... %s\n", urlStr)
	fmt.Printf("保存路径: %s\n", targetPath)

	resp, err := http.Get(urlStr)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("无效状态码: %d %s", resp.StatusCode, resp.Status)
	}

	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	out, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	fmt.Println("下载完成")
	return nil
}

// getMD5Hash returns the MD5 hash of the file
func getMD5Hash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
func downloadFromArgs(config *QiuConfig) error {

	downloadURL := config.BaseURL + config.BucketPath
	targetPath := config.TargetPath
	if targetPath == "" {
		targetPath = filepath.Base(config.BucketPath)
	}
	log.Println("downloadURL: ", downloadURL)
	log.Println("targetPath: ", targetPath)
	if err := download(prepareURL(downloadURL, config), targetPath); err != nil {
		return fmt.Errorf("文件下载失败: %w", err)
	}

	// Get MD5 hash of the downloaded file
	md5Hash, err := getMD5Hash(targetPath)
	if err != nil {
		return fmt.Errorf("获取文件MD5哈希失败: %w", err)
	}
	log.Printf("文件 %s 的 MD5 哈希: %s\n", targetPath, md5Hash)

	return nil
}

func main() {

	log.SetFlags(log.Ldate | log.Ltime)
	log.SetPrefix("[七牛下载] ")
	// 解析命令行参数
	bucketPath := flag.String("bucket", "", "bucketPath")
	showHelp := flag.Bool("help", false, "显示帮助信息")

	flag.Parse()

	if *showHelp {
		fmt.Println("七牛云文件下载工具")
		fmt.Println("用法:")
		flag.PrintDefaults()
		return
	}
	qiuc := &QiuConfig{}
	if err := qiuc.LoadConfig(); err != nil {
		log.Fatalf("加载配置失败: %v\n", err)
	}

	if *bucketPath != "" {
		log.Printf("bucketPath: %s\n", *bucketPath)
		qiuc.BucketPath = *bucketPath
	}

	if err := downloadFromArgs(qiuc); err != nil {
		fmt.Printf("下载失败: %v\n", err)
		os.Exit(1)
	}
}
