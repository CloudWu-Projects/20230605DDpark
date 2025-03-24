package common

import (
	"os"
	"path/filepath"
	"runtime"
)

// Get the absolute path to the running directory
func GetAppPath() string {
	if path, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		return path
	}
	return os.Args[0]
}

// Determine whether the current system is a Windows system?
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// interface log file path
func GetLogPath() string {
	var path string
	if IsWindows() {
		path = filepath.Join(GetAppPath(), "jilaidian.log")
	} else {
		path = "/var/log/jilaidian.log"
	}
	return path
}

// Different systems get different installation paths
func GetInstallPath() string {

	if IsWindows() {
		return `C:\Program Files\nps`
	}
	return "/etc/nps"

}
