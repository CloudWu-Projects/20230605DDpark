package xinjuncheng

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
)

// GenerateSign 生成签名
// params: 请求参数字典
// key: 由停车场方提供的key
// 返回签名字符串（大写）
func GenerateSign(params map[string]string, key string) string {
	// 1. 过滤掉值为空的参数
	filteredParams := make(map[string]string)
	for k, v := range params {
		if v != "" {
			filteredParams[k] = v
		}
	}

	// 2. 按参数名ASCII码从小到大排序
	keys := make([]string, 0, len(filteredParams))
	for k := range filteredParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 3. 拼接键值对字符串
	var dataBuilder strings.Builder
	for i, k := range keys {
		if i > 0 {
			dataBuilder.WriteString("&")
		}
		dataBuilder.WriteString(k)
		dataBuilder.WriteString("=")
		dataBuilder.WriteString(filteredParams[k])
	}
	dataString := dataBuilder.String()

	// 4. 拼接key得到signString
	signString := dataString + "&key=" + key

	// 5. 计算MD5并转为大写
	md5Sum := md5.Sum([]byte(signString))
	sign := fmt.Sprintf("%X", md5Sum)

	return sign
}
