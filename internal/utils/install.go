package utils

import (
	"fmt"
	"html/template"
	"io"
	"jilaidian_go/pkg/common"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const Simply_SystemdScript = `[Unit]
Description={{.Description}}
After=network.target                                                                                                      
                                                                                                                          
[Service]                                                                                                                 
Type=simple                                                                                                               
Restart=always                                                                                                            
ExecStart={{.Path|cmdEscape}}                                          
                                                                                                                          
[Install]                                                                                                                 
WantedBy=multi-user.target
`

const SystemdScript = `[Unit]
Description={{.Description}}
ConditionFileIsExecutable={{.Path|cmdEscape}}
{{range $i, $dep := .Dependencies}} 
{{$dep}} {{end}}
[Service]
LimitNOFILE=65536
StartLimitInterval=5
StartLimitBurst=10
ExecStart={{.Path|cmdEscape}}{{range .Arguments}} {{.|cmd}}{{end}}
{{if .ChRoot}}RootDirectory={{.ChRoot|cmd}}{{end}}
{{if .WorkingDirectory}}WorkingDirectory={{.WorkingDirectory|cmdEscape}}{{end}}
{{if .UserName}}User={{.UserName}}{{end}}
{{if .ReloadSignal}}ExecReload=/bin/kill -{{.ReloadSignal}} "$MAINPID"{{end}}
{{if .PIDFile}}PIDFile={{.PIDFile|cmd}}{{end}}
{{if and .LogOutput .HasOutputFileSupport -}}
StandardOutput=file:/var/log/{{.Name}}.out
StandardError=file:/var/log/{{.Name}}.err
{{- end}}
Restart=always
RestartSec=120
[Install]
WantedBy=multi-user.target
`

// ServiceConfig 结构体用于存储 systemd 服务的配置信息
type ServiceConfig struct {
	Description          string
	Path                 string
	Dependencies         []string
	Arguments            []string
	ChRoot               string
	WorkingDirectory     string
	UserName             string
	ReloadSignal         string
	PIDFile              string
	LogOutput            bool
	HasOutputFileSupport bool
	Name                 string
}

// cmdEscape 函数用于转义命令中的特殊字符
func cmdEscape(s string) string {
	// 这里可以添加具体的转义逻辑
	return s
}

// WriteSystemdScript 将生成的 systemd 脚本写入指定路径
func WriteSystemdScript(targetPath string, config *ServiceConfig) error {
	// 检查目标路径的父目录是否存在
	parentDir := filepath.Dir(targetPath)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		return fmt.Errorf("parent directory %s does not exist", parentDir)
	}

	// 创建模板
	tmpl, err := template.New("systemd").Funcs(template.FuncMap{
		"cmdEscape": cmdEscape,
	}).Parse(Simply_SystemdScript)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	// 打开文件
	file, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", targetPath, err)
	}
	defer file.Close()

	// 执行模板并写入文件
	err = tmpl.Execute(file, config)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	fmt.Printf("Successfully wrote script to %s\n", targetPath)
	return nil
}
func copyStaticFile(srcPath, bin string) string {

	binPath, _ := filepath.Abs(os.Args[0])
	if !common.IsWindows() {
		if _, err := copyFile(filepath.Join(srcPath, bin), "/usr/bin/"+bin); err != nil {
			if _, err := copyFile(filepath.Join(srcPath, bin), "/usr/local/bin/"+bin); err != nil {
				log.Fatalln(err)
			} else {
				copyFile(filepath.Join(srcPath, bin), "/usr/local/bin/"+bin+"-update")
				chMod("/usr/local/bin/"+bin+"-update", 0755)
				binPath = "/usr/local/bin/" + bin
			}
		} else {
			copyFile(filepath.Join(srcPath, bin), "/usr/bin/"+bin+"-update")
			chMod("/usr/bin/"+bin+"-update", 0755)
			binPath = "/usr/bin/" + bin
		}
	} else {
		copyFile(filepath.Join(srcPath, bin+".exe"), filepath.Join(common.GetAppPath(), bin+"-update.exe"))
		copyFile(filepath.Join(srcPath, bin+".exe"), filepath.Join(common.GetAppPath(), bin+".exe"))
	}
	chMod(binPath, 0755)
	return binPath
}

// 生成目录并拷贝文件
func copyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()
	//分割path目录
	destSplitPathDirs := strings.Split(dest, string(filepath.Separator))

	//检测时候存在目录
	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + string(filepath.Separator)
			b, _ := pathExists(destSplitPath)
			if b == false {
				log.Println("mkdir:" + destSplitPath)
				//创建目录
				err := os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

// 检测文件夹路径时候存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func chMod(name string, mode os.FileMode) {
	if !common.IsWindows() {
		os.Chmod(name, mode)
	}
}
func InstallSystemd(binPath string) {
	config := &ServiceConfig{
		Description: "yilingshequ Service",
		Path:        binPath,
	}

	err := WriteSystemdScript("/etc/systemd/system/"+common.GetAppName()+".service", config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
func InstallServer() {
	path := common.GetInstallPath()
	binPath := copyStaticFile(common.GetAppPath(), common.GetAppName())
	log.Println("install ok!")
	log.Println("Static files and configuration files in the current directory will be useless")
	log.Println("The new configuration file is located in", path, "you can edit them")
	if !common.IsWindows() {
		log.Printf(`You can start with:
%s start|stop|restart|uninstall|update 
anywhere!\n`, common.GetAppName())
		InstallSystemd(binPath)
	} else {
		log.Printf(`You can copy executable files to any directory and start working with:
%s start|stop|restart|uninstall|update 
now!\n`, common.GetAppName())
	}
	chMod(common.GetLogPath(), 0777)
}
