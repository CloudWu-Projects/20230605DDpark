package main

import (
	"embed"
	"flag"
	"jilaidian_go/lib/api"
	"jilaidian_go/lib/config"
	"jilaidian_go/lib/logger"
	"jilaidian_go/test"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed www/config.html
var content embed.FS

func main() {
	// 解析命令行参数
	install := flag.Bool("install", false, "安装服务")
	runExample := flag.Bool("example", false, "运行示例代码")
	flag.Parse()

	// 初始化日志
	logger.Logger.Info("初始化日志")

	if *install {
		logger.Logger.Info("安装服务")

	}
	// 如果是运行示例代码
	if *runExample {
		logger.Logger.Info("运行示例代码")
		test.RunExample()
		return
	}

	r := mux.NewRouter()
	// 设置路由
	handler := api.NewHandler()
	handler.SetupRoutes(r)
	configHandler := api.NewConfigHandler(content)
	configHandler.SetupRoutes(r)
	// 启动服务器
	port := config.Global.Server.Port
	logger.Logger.Info("启动服务器 port:", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Logger.Error("启动服务器失败", err)
	}
}
