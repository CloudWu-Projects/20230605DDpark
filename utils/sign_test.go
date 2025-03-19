package utils

import (
	"sort"
	"testing"
)

func TestSign(t *testing.T) {
	carNumber := "京A5566TT"
	query_time := 1672129283
	ukey := "dny75"
	// 构造请求数据
	data := map[string]interface{}{
		"query_time": query_time,
		"car_number": carNumber,
	}
	// 提取 keys 并进行反向排序
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys))) // 反向排序

	// 按反向排序的 key 重新构造 map
	sortedData := make(map[string]interface{})
	for _, k := range keys {
		sortedData[k] = data[k]
	}

	// 生成签名
	datasign := GenerateSign(data, ukey)
	t.Log("data sign:", datasign)
	sign := GenerateSign(sortedData, ukey)
	t.Log("sortedData sign:", sign)
	if sign != "9BBD511E3178A91E66AD4C6A043CEC55" {
		t.Error("sign error")
	}
}
