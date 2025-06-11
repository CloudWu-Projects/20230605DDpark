package tingcheyunservice

type InOutBase struct {
	ServiceName string `json:"service_name"`
	ParkID      int    `json:"park_id"`
	Sign        string `json:"sign"`
}
type InPark struct {
	InOutBase
	Data struct {
		CarNumber    string `json:"car_number"`
		InTime       int    `json:"in_time"`
		CarType      string `json:"car_type"`
		CType        string `json:"c_type"`
		UID          string `json:"uid"`
		OperatorName string `json:"operator_name"`
		OrderID      string `json:"order_id"`
		EmptyPlot    int    `json:"empty_plot"`
		InChannelID  string `json:"in_channel_id"`
		WorksiteID   int    `json:"worksite_id"`
		InRemark     string `json:"in_remark"`
		ForceUpdate  int    `json:"force_update"`
		LicenseColor int    `json:"license_color"`
		PicAddr      string `json:"pic_addr"`
	}
}

type InParkResponse struct {
	State       int    `json:"state"`
	OrderID     string `json:"order_id"`
	ParkID      int    `json:"park_id"`
	ServiceName string `json:"service_name"`
	Errmsg      string `json:"errmsg"`
}

type OutPark struct {
	InOutBase
	Data struct {
		ServiceName      string  `json:"service_name"`
		CarNumber        string  `json:"car_number"`
		OutTime          int     `json:"out_time"`
		OutType          string  `json:"out_type"`
		CType            string  `json:"c_type"`
		EmptyPlot        int     `json:"empty_plot"`
		CarType          string  `json:"car_type"`
		Duration         int     `json:"duration"`
		UID              string  `json:"uid"`
		LicenseColor     int     `json:"license_color"`
		OutUID           string  `json:"out_uid"`
		Freereasons      string  `json:"freereasons"`
		OutOperatorName  string  `json:"out_operator_name"`
		InRemark         string  `json:"in_remark"`
		PayType          string  `json:"pay_type"`
		ReduceAmount     string  `json:"reduce_amount"`
		PicAddr          string  `json:"pic_addr"`
		ElectronicPay    string  `json:"electronic_pay"`
		Remark           string  `json:"remark"`
		Total            float32 `json:"total"`
		LicenseType      int     `json:"licence_type"`
		CashPay          string  `json:"cash_pay"`
		ParkID           string  `json:"park_id"`
		InChannelID      string  `json:"in_channel_id"`
		DataTarget       string  `json:"data_target"`
		TicketID         string  `json:"ticket_id"`
		AuthCode         string  `json:"auth_code"`
		OperatorName     string  `json:"operator_name"`
		InTime           int     `json:"in_time"`
		AmountReceivable string  `json:"amount_receivable"`
		OutChannelID     string  `json:"out_channel_id"`
		OrderID          string  `json:"order_id"`
		ForceUpdate      int     `json:"force_update"`
	}
}
