package utils

import (
	"bufio"
	"fmt"
	"jilaidian_go/pkg/common"
	"os"
	"strings"
)

// ReadLogFile reads the content of the log file and returns it as a string.
func ReadLogFile() (string, error) {
	logPath := common.GetLogPath()
	// 打开日志文件
	file, err := os.Open(logPath)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return "", err
	}
	defer file.Close()

	// 创建一个Scanner来逐行读取文件内容
	scanner := bufio.NewScanner(file)
	var content strings.Builder

	// 逐行读取文件内容并累积到content中
	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}

	// 检查扫描过程中是否有错误
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading log file:", err)
		return "", err
	}

	return content.String(), nil
}
