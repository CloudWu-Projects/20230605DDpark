package models

// ChargeInfo 充电信息结构体
type ChargeInfo struct {
	ParkID  string `json:"parkId"`
	PlateNo string `json:"plateNo"`
}

// QueryOrderRequest 查询订单请求结构体
type QueryOrderRequest struct {
	ServiceName string `json:"service_name"`
	Sign        string `json:"sign"`
	ParkID      int    `json:"park_id"`
	Data        struct {
		CarNumber string `json:"car_number"`
	} `json:"data"`
}

// QueryOrderResponse 查询订单响应结构体
type QueryOrderResponse struct {
	ServiceName string `json:"service_name"`
	Sign        string `json:"sign"`
	ParkID      int    `json:"park_id"`
	Data        struct {
		CarNumber string `json:"car_number"`
		OrderID   string `json:"order_id"`
		InTime    int64  `json:"in_time"`
	} `json:"data"`
	ErrMsg string `json:"errmsg"`
	State  int    `json:"state"`
}

// DiscountNoticeRequest 下发优惠信息请求结构体
type DiscountNoticeRequest struct {
	ServiceName string `json:"service_name"`
	Sign        string `json:"sign"`
	ParkID      int    `json:"park_id"`
	Data        struct {
		CarNumber         string  `json:"car_number"`
		OrderID           string  `json:"order_id"`
		ReduceAmount      float64 `json:"reduce_amount"`
		DeductionTime     int     `json:"deduction_time"`
		DeductionMoney    int     `json:"deduction_money"`
		Duration          int     `json:"duration"`
		Remark            string  `json:"remark"`
		StartChargingTime string  `json:"start_charging_time"`
		StopChargingTime  string  `json:"stop_charging_time"`
		UUID              string  `json:"uuid"`
	} `json:"data"`
}
