package utils

import (
	"encoding/json"
	"testing"
)

func TestSign(t *testing.T) {
	carNumber := "京A5566TT"

	ukey := "dny75"
	// 构造请求数据
	data := struct {
		CarNumber string `json:"car_number"`
	}{
		CarNumber: carNumber,
	}
	dataBytes, _ := json.Marshal(data)
	dataStr := string(dataBytes)
	// 生成签名
	sign := GenerateSignString(dataStr, ukey)

	t.Log("data sign:", sign)
	if sign != "9BBD511E3178A91E66AD4C6A043CEC55" {
		t.Error("sign error")
	}
}
