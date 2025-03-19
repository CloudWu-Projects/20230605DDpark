package test

import (
	"fmt"
	"jilaidian_go/client"
	"jilaidian_go/models"
	"jilaidian_go/service"
)

// RunExample 运行示例代码
func RunExample() {
	// 创建服务和客户端
	chargeService := service.NewChargeService()
	apiClient := client.NewAPIClient()

	// 1. 接收吉来电的充电信息
	chargeData := models.ChargeInfo{
		ParkID:  "10051557",
		PlateNo: "京A5566TT",
	}
	if err := chargeService.HandleChargeInfo(chargeData); err != nil {
		fmt.Println("Error handling charge info:", err)
		return
	}

	// 2. 通过阿里云查询订单信息
	parkID := 10033791
	carNumber := "京A5566TT"
	orderID, err := apiClient.QueryOrder(parkID, carNumber)
	if err != nil {
		fmt.Println("Error querying order:", err)
		return
	}

	// 3. 下发优惠信息
	reduceAmount := 8.0
	deductionTime := 4
	deductionMoney := 5
	if err := apiClient.SendDiscountNotice(parkID, carNumber, orderID, reduceAmount, deductionTime, deductionMoney); err != nil {
		fmt.Println("Error sending discount notice:", err)
		return
	}
	fmt.Println("Discount notice sent successfully!")
}
