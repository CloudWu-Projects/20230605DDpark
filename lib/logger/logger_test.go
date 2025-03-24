package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoggerInitialization(t *testing.T) {
	// 确保初始化后Logger不为nil
	if Logger == nil {
		t.Fatal("Logger was not initialized")
	}

	// 验证handler的日志级别设置
	//handler := Logger.Handler().(*slog.TextHandler)
	Logger.Info("init logger")
}

func TestLogOutput(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	oldFilename := LumberjackLogger.Filename
	defer func() { LumberjackLogger.Filename = oldFilename }()
	LumberjackLogger.Filename = filepath.Join(tempDir, "test.log")

	Logger = nil
	//init()

	testMsg := "TEST_LOG_MESSAGE_123"
	Logger.Info(testMsg)

	logFileContent, err := os.ReadFile(LumberjackLogger.Filename)
	if err != nil {
		t.Fatal("Failed to read log file:", err)
	}

	if !strings.Contains(string(logFileContent), testMsg) {
		t.Errorf("Log file missing test message: %s", testMsg)
	}
	if !strings.Contains(string(logFileContent), "level=INFO") {
		t.Error("Log missing level field")
	}
	if !strings.Contains(string(logFileContent), "msg="+testMsg) {
		t.Error("Log missing message field")
	}
}

func TestLogRotationConfig(t *testing.T) {
	expected := map[string]interface{}{
		"MaxSize":    10,
		"MaxBackups": 3,
		"MaxAge":     28,
		"Compress":   true,
	}

	if LumberjackLogger.MaxSize != expected["MaxSize"] {
		t.Errorf("MaxSize: expected %d, got %d", expected["MaxSize"], LumberjackLogger.MaxSize)
	}
	if LumberjackLogger.MaxBackups != expected["MaxBackups"] {
		t.Errorf("MaxBackups: expected %d, got %d", expected["MaxBackups"], LumberjackLogger.MaxBackups)
	}
	if LumberjackLogger.MaxAge != expected["MaxAge"] {
		t.Errorf("MaxAge: expected %d, got %d", expected["MaxAge"], LumberjackLogger.MaxAge)
	}
	if LumberjackLogger.Compress != expected["Compress"] {
		t.Errorf("Compress: expected %t, got %t", expected["Compress"], LumberjackLogger.Compress)
	}
}
