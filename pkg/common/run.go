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
func GetAppName() string {
	a := filepath.Base(os.Args[0])
	if len(a) < 2 {
		return filepath.Base(GetAppPath())
	}
	return a
}

// Determine whether the current system is a Windows system?
func IsWindows() bool {
	return runtime.GOOS == "windows"
}
func IsLinux() bool {
	return runtime.GOOS == "linux"
}
func IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

// interface log file path
func GetLogPath() string {
	var path string
	if IsWindows() {
		path = filepath.Join(GetAppPath(), GetAppName()+".log")
	} else if IsLinux() {
		path = "/var/log/" + GetAppName() + ".log"
	} else if IsMacOS() {
		path = "~/Library/Logs/" + GetAppName() + ".log"
	}
	return path
}
func GetConfigPath() string {
	var path string
	if IsWindows() {
		path = filepath.Join(GetAppPath(), GetAppName()+".json")
	} else if IsLinux() {
		path = "/etc/" + GetAppName() + ".json"
	} else if IsMacOS() {
		path = "~/Library/Logs/" + GetAppName() + ".json"
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
