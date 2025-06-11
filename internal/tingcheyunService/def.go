package tingcheyunservice

/*
   {
     "service_name": "in_park",
     "park_id": 31270,
     "sign": "F88C2D64DFE56872D3C8C54AF67EE8C7",
     "data": {
       "car_number": "京GH0093",
       "in_time": 1490875218,
       "car_type": "大车",
       "c_type": "临时车",
       "uid": "325101",
       "operator_name":"系统管理员",
       "order_id": "325101",
       "empty_plot": 20,
       "in_channel_id": "A1",
       "worksite_id": 23,
       "in_remark": "入场信息备注",
       "force_update": 0，
       "license_color":4
     }
   }
*/
type In_out_Base struct {
	ServiceName string `json:"service_name"`
	ParkID      int    `json:"park_id"`
	Sign        string `json:"sign"`
}
type InPark struct {
	ServiceName string `json:"service_name"`
	ParkID      int    `json:"park_id"`
	Sign        string `json:"sign"`
	Data        struct {
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

/*
{
  "state": 1,
  "order_id": "325101",
  "park_id": 31270,
  "service_name": "in_park",
  "errmsg": " send success!"
}
*/
type InParkResponse struct {
	State       int    `json:"state"`
	OrderID     string `json:"order_id"`
	ParkID      int    `json:"park_id"`
	ServiceName string `json:"service_name"`
	Errmsg      string `json:"errmsg"`
}

/*
{
"service_name": "out_park",
"sign": "10836C5F35023C0689C49978B0525D66",
"park_id": "45126",
"data": {
"uid": "01",
"out_operator_name":"系统管理员",
"cash_pay": "0.0",
"pay_type": "cash",
"electronic_pay": "0.0",
"in_time": 1590042193,
"empty_plot": 1,
"in_channel_id": "83350",
"out_channel_id": "83351",
"order_id": "8526702",
"car_number": "临Q8BY97",
"auth_code": "",
"freereasons": "",
"c_type": "临时车",
"duration": 16,
"total": "0.0",
"out_time": 1590043207,
"car_type": "小型车",
"amount_receivable": "0.0",
"flag":0,
"reduction_rules":"10小时抵用券",
"charge_status":2,
"amount_receivable":"8.0",
"license_color":4
}
}
*/

type OutPark struct {
	ServiceName string `json:"service_name"`
	Sign        string `json:"sign"`
	ParkID      int    `json:"park_id"`
	Data        struct {
		UID              string `json:"uid"`
		OutOperatorName  string `json:"out_operator_name"`
		CashPay          string `json:"cash_pay"`
		PayType          string `json:"pay_type"`
		ElectronicPay    string `json:"electronic_pay"`
		InTime           int    `json:"in_time"`
		EmptyPlot        int    `json:"empty_plot"`
		InChannelID      string `json:"in_channel_id"`
		OutChannelID     string `json:"out_channel_id"`
		OrderID          string `json:"order_id"`
		CarNumber        string `json:"car_number"`
		AuthCode         string `json:"auth_code"`
		Freereasons      string `json:"freereasons"`
		CType            string `json:"c_type"`
		Duration         int    `json:"duration"`
		Total            string `json:"total"`
		OutTime          int    `json:"out_time"`
		CarType          string `json:"car_type"`
		AmountReceivable string `json:"amount_receivable"`
		Flag             int    `json:"flag"`
		ReductionRules   string `json:"reduction_rules"`
		ChargeStatus     int    `json:"charge_status"`
		LicenseColor     int    `json:"license_color"`
		PicAddr          string `json:"pic_addr"`
	}
}
