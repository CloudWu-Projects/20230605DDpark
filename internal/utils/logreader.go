package utils

import (
	"bufio"
	"fmt"
	"jilaidian_go/pkg/common"
	"os"
	"strings"
)

// ReadLogFile reads the content of the log file and returns it as a string.
func ReadLogFile(lastBytes int64) (string, error) {
	logPath := common.GetLogPath()
	// 打开日志文件
	file, err := os.Open(logPath)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return "", err
	}
	defer file.Close()
	// get totallines from file
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return "", err
	}
	totalSize := fileInfo.Size()
	if lastBytes <= 0 || lastBytes > totalSize {
		lastBytes = totalSize // 如果lastLine不合法，重置为0
	}
	file.Seek(totalSize-lastBytes, 0)
	// 创建一个Scanner来逐行读取文件内容
	scanner := bufio.NewScanner(file)
	var content strings.Builder
	content.WriteString("totalSize: " + fmt.Sprintf("%d", totalSize) + "\n")
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
