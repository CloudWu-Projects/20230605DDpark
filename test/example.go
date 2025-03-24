package test

import (
	"jilaidian_go/lib/config"
	"jilaidian_go/lib/models"
	"jilaidian_go/lib/service"
)

// RunExample 运行示例代码
func RunExample() {
	// 创建服务和客户端
	chargeService := service.NewChargeService()
	//apiClient := client.NewAPIClient()

	parkInfo := config.ParkInfo{
		Deduction_money: 100,
		Deduction_time:  60,
		Parkid:          10051557,
		Ukey:            "your_ukey",
	}
	chargeService.ProcessChargingAndDiscount(&parkInfo, models.ChargeInfo{PlateNo: "京A5566TT", ParkID: "10033791"})

}
