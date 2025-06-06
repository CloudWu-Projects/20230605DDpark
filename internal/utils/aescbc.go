package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

// PKCS7Padding pads the plaintext to be a multiple of block size.
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// CBCEncrypt encrypts plaintext using AES CBC mode.
// key must be 16, 24, or 32 bytes (AES-128, AES-192, AES-256).
// iv must be 16 bytes.
func CBCEncrypt(plaintext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != block.BlockSize() {
		return nil, errors.New("IV length must equal block size")
	}
	plaintext = PKCS7Padding(plaintext, block.BlockSize())
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func CBCEncrypt_Base64(plaintext, key, iv string) (string, error) {
	ciphertext, err := CBCEncrypt([]byte(plaintext), []byte(key), []byte(iv))
	if err != nil {
		return "", fmt.Errorf("CBCEncrypt failed: %w", err)
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func CBCDecrypt(ciphertext, key, iv []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != block.BlockSize() {
		return nil, errors.New("IV length must equal block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	// Remove padding
	plaintext = PKCS7UnPadding(plaintext, block.BlockSize())
	return plaintext, nil
}
func PKCS7UnPadding(data []byte, blockSize int) []byte {
	if len(data) == 0 {
		return data
	}
	padding := int(data[len(data)-1])
	if padding < 1 || padding > blockSize {
		return data // Invalid padding, return original data
	}
	return data[:len(data)-padding]
}
func CBCDecrypt_Base64(ciphertext, key, iv string) (string, error) {
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}
	plaintext, err := CBCDecrypt(decodedCiphertext, []byte(key), []byte(iv))
	if err != nil {
		return "", fmt.Errorf("CBCDecrypt failed: %w", err)
	}
	return string(plaintext), nil
}
