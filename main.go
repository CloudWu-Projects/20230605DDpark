package main

import (
	"flag"
	"jilaidian_go/api"
	"jilaidian_go/config"
	"jilaidian_go/logger"
	"jilaidian_go/test"

	"github.com/gin-gonic/gin"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config/config.json", "配置文件路径")
	runExample := flag.Bool("example", false, "运行示例代码")
	flag.Parse()

	// 初始化日志
	logger.Logger.Info("初始化日志")

	// 加载配置
	if err := config.Load(*configPath); err != nil {
		logger.Logger.Error("加载配置失败", err)
		return
	}
	logger.Logger.Info("配置加载成功")

	// 如果是运行示例代码
	if *runExample {
		logger.Logger.Info("运行示例代码")
		test.RunExample()
		return
	}

	// 初始化 Gin
	router := gin.Default()

	// 设置路由
	handler := api.NewHandler()
	handler.SetupRoutes(router)

	// 启动服务器
	port := config.Global.Server.Port
	logger.Logger.Info("启动服务器 port:", port)
	if err := router.Run(":" + port); err != nil {
		logger.Logger.Error("服务器启动失败", err)
	}
}
