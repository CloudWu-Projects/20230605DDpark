package xinjuncheng

import (
	"crypto/md5"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

// 注意：根据文档示例，platformNo 应该是小写的
type XjcRequest struct {
	DiscountNumber string `json:"discountNumber"` // 优惠流水号
	PlateNumber    string `json:"plateNumber"`    // 车牌号
	PlatformNo     string `json:"platformNo"`     // 平台编号 - 注意是小写！
	DiscountTime   int    `json:"discountTime"`   // 优惠时长
	DiscountMoney  int    `json:"discountMoney"`  // 优惠金额
	Timestamp      int64  `json:"timestamp"`
	Sign           string `json:"sign"` // 签名本身不参与签名计算
}

// GetSignParams 获取参与签名的参数（过滤空值）
func (r *XjcRequest) GetSignParams() map[string]string {
	params := make(map[string]string)

	// 按照文档示例的字段名（注意大小写）
	if r.DiscountNumber != "" {
		params["discountNumber"] = r.DiscountNumber
	}
	if r.PlateNumber != "" {
		params["plateNumber"] = r.PlateNumber
	}
	if r.PlatformNo != "" {
		params["platformNo"] = r.PlatformNo // 注意：小写！
	}
	if r.DiscountTime != 0 {
		params["discountTime"] = strconv.Itoa(r.DiscountTime)
	}
	if r.DiscountMoney != 0 {
		params["discountMoney"] = strconv.Itoa(r.DiscountMoney)
	}
	if r.Timestamp != 0 {
		params["timestamp"] = strconv.FormatInt(r.Timestamp, 10)
	}

	return params
}

// CalculateSign 计算签名
func (r *XjcRequest) CalculateSign(key string) string {
	params := r.GetSignParams()
	return CalculateSignFromMap(params, key)
}

// CalculateSignFromMap 从map计算签名
func CalculateSignFromMap(params map[string]string, key string) string {
	// 1. 过滤空值参数
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
	var builder strings.Builder
	for i, k := range keys {
		if i > 0 {
			builder.WriteString("&")
		}
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(filteredParams[k])
	}
	dataString := builder.String()

	// 4. 拼接key得到signString
	signString := dataString + "&key=" + key

	// 5. 计算MD5并转为大写
	md5Sum := md5.Sum([]byte(signString))
	sign := fmt.Sprintf("%X", md5Sum)

	return sign
}

// VerifySign 验证签名
func (r *XjcRequest) VerifySign(key string) bool {
	savedSign := r.Sign
	calculatedSign := r.CalculateSign(key)
	log.Println("savedSign:", savedSign)
	log.Println("calculatedSign:", calculatedSign)
	return strings.EqualFold(savedSign, calculatedSign)
}
