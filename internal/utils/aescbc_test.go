package utils

import (
	"os"
	"testing"
)

func TestCbc(t *testing.T) {

	// open a.txt and read the content to a string
	content, err := os.ReadFile("a.txt")
	if err != nil {
		t.Fatalf("failed to read a.txt: %v", err)
	}
	fileContent := string(content)
	t.Log("file content:", fileContent)

	secret := "1234567890abcdef"
	iv := "1234567890abcdef"
	// 生成签名
	sign, _ := CBCEncrypt_Base64((fileContent), secret, iv)
	//t.Log("data sign:", sign)
	wanted_str := "il7B0BSEjFdzpyKzfOFpvg/Se1CP802RItKYFPfSLRxJ3jf0bVl9hvYOEktPAYW2nd7S8MBcyHYyacHKbISq5iTmDzG+ivnR+SZJv3USNTYVMz9rCQVSxd0cLlqsJauko79NnwQJbzDTyLooYoIwz75qBOH2/xOMirpeEqRJrF/EQjWekJmGk9RtboXePu2rka+Xm51syBPhiXJAq0GfbfaFu9tNqs/e2Vjja/ltE1M0lqvxfXQ6da6HrThsm5id4ClZFIi0acRfrsPLRixS/IQYtksxghvJwbqOsbIsITail9Ayy4tKcogeEZiOO+4Ed264NSKmk7l3wKwJLAFjCFogBx8GE3OBz4pqcAn/ydA="
	if sign == wanted_str {
		t.Log("sign ok")
		return
	}

	t.Fatalf("sign error, got: %s", sign)
}
