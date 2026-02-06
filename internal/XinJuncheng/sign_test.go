package xinjuncheng

import (
	"fmt"
	"testing"
)

// GenerateSignExample 使用示例数据生成签名（与文档中的示例匹配）
func GenerateSignExample() string {
	params := map[string]string{
		"discountNumber": "169157199496923925786",
		"plateNumber":    "浙A6Y557",
		"discountTime":   "120",
		"timestamp":      "1691573927454",
		"platformNo":     "test001",
	}
	key := "123456"
	return GenerateSign(params, key)
}

func TestSign(t *testing.T) {
	// 测试示例
	sign := GenerateSignExample()
	fmt.Println("生成的签名:", sign)
	fmt.Println("预期签名: 9B2643F390139D67DB98F6D8CE34B6A2")
	fmt.Println("是否一致:", sign == "9B2643F390139D67DB98F6D8CE34B6A2")

	// 其他测试用例
	fmt.Println("\n其他测试用例:")

	// 示例1：包含空值的参数
	params1 := map[string]string{
		"name":  "John",
		"age":   "30",
		"email": "",
		"city":  "New York",
		"phone": "",
	}
	sign1 := GenerateSign(params1, "mySecretKey")
	fmt.Println("签名1:", sign1)

	// 示例2：区分大小写
	params2 := map[string]string{
		"Name": "John",
		"name": "Doe",
		"AGE":  "25",
		"age":  "30",
	}
	sign2 := GenerateSign(params2, "testKey")
	fmt.Println("签名2:", sign2)
}
