package models

type Notification_charge_end_order_info_Request struct {
	/*
	   {
	         "StartChargeSeq":"123456789201704111645012",
	         "StationID":"3702120244",
	         "StationName":"北京中环集中办公区充电站",
	         "ConnectorID":"1101020190301",
	   "ParkID":"3502060030004",
	         "ConnectorName":"301号交流",
	         "PlateNum":"京AD06088",
	         "ParkNo":"017",
	         "FreeParkingTimes":30.00
	   }
	*/
	StartChargeSeq   string  `json:"StartChargeSeq"`
	StationID        string  `json:"StationID"`
	StationName      string  `json:"StationName"`
	ConnectorID      string  `json:"ConnectorID"`
	ParkID           string  `json:"ParkID"`
	ConnectorName    string  `json:"ConnectorName"`
	PlateNum         string  `json:"PlateNum"`
	ParkNo           string  `json:"ParkNo"`
	FreeParkingTimes float64 `json:"FreeParkingTimes"`
}

type Notification_charge_end_order_info_Response struct {
	/*
	   {
	   "StartChargeSeq":"201801120000052917",
	   "ConfirmResult":0
	   "PlateAutResult":1
	   }
	*/
	StartChargeSeq string `json:"StartChargeSeq"`
	ConfirmResult  int    `json:"ConfirmResult"`
	PlateAutResult int    `json:"PlateAutResult"`
}
