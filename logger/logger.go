package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/petermattis/goid"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	pid := os.Getpid()
	ggid := goid.Get()
	logMessage := fmt.Sprintf("%s [PID: %d][GID: %d] [%s] %s\n",
		entry.Time.Format("2006-01-02 15:04:05"),
		pid, ggid, entry.Level, entry.Message)
	//logMessage = strings.TrimSpace(logMessage)
	return []byte(logMessage), nil
}

var LumberjackLogger = &lumberjack.Logger{
	Filename:   "app.log",
	MaxSize:    10,
	MaxBackups: 3,
	MaxAge:     28,
	Compress:   true,
}

func init() {
	_logInstance := logrus.New()
	multiWriter := io.MultiWriter(os.Stdout, LumberjackLogger)
	_logInstance.SetOutput(multiWriter)
	_logInstance.SetFormatter(&CustomFormatter{})
	_logInstance.SetLevel(logrus.DebugLevel)
	// 创建 Logger
	Logger = _logInstance
}
