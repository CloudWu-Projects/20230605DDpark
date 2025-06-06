package main

import (
	"fmt"
	. "jilaidian_go/internal/utils"
	"os"
)

func main() {
	// Example usage
	key := "1234567890abcdef"
	iv := "1234567890abcdef"

	plaintext := "This is a secret message to be encrypted with padding-5"

	fmt.Printf("Key: %s\n", key)
	fmt.Printf("IV: %s\n", iv)
	fmt.Printf("Plaintext: %s\n", plaintext)
	cc, _ := CBCEncrypt_Base64(plaintext, key, iv)
	// base64 encoded ciphertext
	// Example ciphertext:
	// 0f8c1b2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b
	fmt.Printf("Ciphertext: %s\n", cc)
	//K3z9fiEYzvO7wctaoLCIUOW1xqODo9UQxivdJa/NCIArpt4/fUF2ymnt58velWegWGaKKiNuNSggCiRfUQEREw==

	// open a.txt and read the content to a string
	content, err := os.ReadFile(".\\internal\\utils\\a.txt")
	if err != nil {
		fmt.Errorf("failed to read a.txt: %v", err)
	}
	fileContent := string(content)

	// 生成签名
	cc, _ = CBCEncrypt_Base64(fileContent, key, iv)

	fmt.Printf("Ciphertext: %s\n", cc)
	fmt.Println("il7B0BSEjFdzpyKzfOFpvg/Se1CP802RItKYFPfSLRxJ3jf0bVl9hvYOEktPAYW2nd7S8MBcyHYyacHKbISq5iTmDzG+ivnR+SZJv3USNTYVMz9rCQVSxd0cLlqsJauko79NnwQJbzDTyLooYoIwz75qBOH2/xOMirpeEqRJrF/EQjWekJmGk9RtboXePu2rka+Xm51syBPhiXJAq0GfbfaFu9tNqs/e2Vjja/ltE1M0lqvxfXQ6da6HrThsm5id4ClZFIi0acRfrsPLRixS/IQYtksxghvJwbqOsbIsITail9Ayy4tKcogeEZiOO+4Ed264NSKmk7l3wKwJLAFjCFogBx8GE3OBz4pqcAn/ydA=")

	bb, err := CBCDecrypt_Base64(cc, key, iv)
	if err != nil {
		fmt.Errorf("failed to decrypt: %v", err)
	}
	fmt.Printf("Decrypted: %s\n", string(bb))
}
