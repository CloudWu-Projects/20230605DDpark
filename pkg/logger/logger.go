package logger

import (
	"fmt"
	"io"
	"jilaidian_go/pkg/common"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/petermattis/goid"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	pid := os.Getpid()
	ggid := goid.Get()

	// 添加文件和行号信息
	var fileInfo string
	if entry.Caller != nil {
		fileInfo = fmt.Sprintf("[%s:%d]", filepath.Base(entry.Caller.File), entry.Caller.Line)
	}

	logMessage := fmt.Sprintf("%s [PID: %d][GID: %d] %s [%s] %s\n",
		entry.Time.Format("2006-01-02 15:04:05"),
		pid, ggid, fileInfo, entry.Level, entry.Message)

	// 添加字段信息
	if len(entry.Data) > 0 {
		fields := make([]string, 0, len(entry.Data))
		for k, v := range entry.Data {
			fields = append(fields, fmt.Sprintf("%s=%v", k, v))
		}
		logMessage = strings.TrimSuffix(logMessage, "\n") + " " + strings.Join(fields, " ") + "\n"
	}

	return []byte(logMessage), nil
}

var LumberjackLogger = &lumberjack.Logger{
	Filename:   common.GetLogPath(),
	MaxSize:    1,
	MaxBackups: 5,
	MaxAge:     28,
	Compress:   true,
}

func init() {

	_logInstance := logrus.New()

	fmt.Println("log path:", common.GetLogPath())
	multiWriter := io.MultiWriter(os.Stdout, LumberjackLogger)
	_logInstance.SetOutput(multiWriter)
	_logInstance.SetFormatter(&CustomFormatter{})
	_logInstance.SetLevel(logrus.DebugLevel)
	log.SetFlags(0) // 移除标准库的默认标志
	log.SetOutput(_logInstance.Writer())

	// 创建 Logger
	Logger = _logInstance
}
