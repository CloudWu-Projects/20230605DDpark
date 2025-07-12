package main

import (
	"flag"
	"fmt"
	api "jilaidian_go/internal/api/handlers"
	"jilaidian_go/internal/config"
	"jilaidian_go/internal/utils"
	"jilaidian_go/internal/yianqiservice"
	"jilaidian_go/pkg/common"
	"jilaidian_go/pkg/logger"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	fmt.Println("\n欢迎使用", common.GetAppName(), os.Args[0])
	// 解析命令行参数
	install := flag.Bool("install", false, "安装服务")
	flag.Parse()

	// 初始化日志

	if *install {
		logger.Logger.Info("安装服务")
		// 写 utils.Simply_SystemdScript 到 /etc/systemd/system/jilaidian.service
		utils.InstallServer()
		return
	}
	logger.Logger.Info(config.Global.Version)
	logger.Logger.Info("ConfigPath:", common.GetConfigPath())
	// 如果是运行示例代码

	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// 设置路由
	handler := yianqiservice.NewHandler()
	handler.SetupRoutes(r)
	configHandler := api.NewConfigHandler()
	configHandler.SetupRoutes(r)
	// 启动服务器
	port := config.Global.ServerConfig.Port
	logger.Logger.Info("启动服务器 port:", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Logger.Error("启动服务器失败", err)
	}
}
