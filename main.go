package main

import (
	"embed"
	"flag"
	"jilaidian_go/api"
	"jilaidian_go/config"
	"jilaidian_go/logger"
	"jilaidian_go/test"
	"net/http"
)

//go:embed www/config.html
var content embed.FS

func main() {
	// 解析命令行参数

	runExample := flag.Bool("example", false, "运行示例代码")
	flag.Parse()

	// 初始化日志
	logger.Logger.Info("初始化日志")

	// 如果是运行示例代码
	if *runExample {
		logger.Logger.Info("运行示例代码")
		test.RunExample()
		return
	}

	// 初始化 Gin
	//router := gin.Default()

	// 设置路由
	handler := api.NewHandler()
	handler.SetupRoutes()
	configHandler := api.NewConfigHandler(content)
	configHandler.SetupRoutes()
	// 启动服务器
	port := config.Global.Server.Port
	logger.Logger.Info("启动服务器 port:", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Logger.Error("启动服务器失败", err)
	}
}
