package models

// 请求结构体
type QueryRequest struct {
	OperatorID string `json:"OperatorID"`
	Data       string `json:"Data"`
	TimeStamp  string `json:"TimeStamp"`
	Seq        string `json:"Seq"`
	Sig        string `json:"Sig"`
}

// ChargeInfo 充电信息结构体
type ChargeInfo struct {
	ParkID  string `json:"parkId"`
	PlateNo string `json:"plateNo"`
}

// 充电记录数据结构
type ChargingRecord struct {
	PortName     string `json:"portName"`
	ParkId       string `json:"parkId"`
	OrderNo      string `json:"orderNo"`
	PlateNo      string `json:"plateNo"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	StationId    string `json:"stationId"`
	StationName  string `json:"stationName"`
	DeviceId     string `json:"deviceId"`
	DeviceName   string `json:"deviceName"`
	SpaceNo      string `json:"spaceNo"`
	Power        string `json:"power"`
	ElecMoney    string `json:"elecMoney"`
	ServiceMoney string `json:"seviceMoney"`
	TotalMoney   string `json:"totalMoney"`
	Sign         string `json:"sign"`
	Time         string `json:"time"`
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
