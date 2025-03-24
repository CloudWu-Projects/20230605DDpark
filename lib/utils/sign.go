package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strings"
)

// GenerateSign 生成签名
func generateSign(data map[string]interface{}, ukey string) string {
	// 将 data 转换为 JSON 字符串
	dataBytes, _ := json.Marshal(data)
	dataStr := string(dataBytes)
	return GenerateSignString(dataStr, ukey)
}
func GenerateSignString(dataStr string, ukey string) string {
	// 拼接 key=ukey
	signStr := dataStr + "key=" + ukey
	// 计算 MD5 并转为大写
	hash := md5.Sum([]byte(signStr))
	hexStr := hex.EncodeToString(hash[:])
	return strings.ToUpper(hexStr)
}
