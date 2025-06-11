package thirdservice

import "jilaidian_go/pkg/client"

// Handler API处理器
type Thirdservice struct {
}
type ThirdInObj struct {
	Api         string `json:"api"`
	ImgUrl      string `json:"imgUrl"`
	InTime      int    `json:"inTime"`
	ParkingId   string `json:"parkingId"`
	PlateNumber string `json:"plateNumber"`
}
type ThirdOutObj struct {
	Api         string `json:"api"`
	ImgUrl      string `json:"imgUrl"`
	OutTime     int    `json:"outTime"`
	ParkingId   string `json:"parkingId"`
	PlateNumber string `json:"plateNumber"`
}

func PostInData(data ThirdInObj) error {
	apiClient := client.NewAPIClient()
	// 构造请求数据
	_, err := apiClient.SendRequest("", data)
	return err
}

func PostOutData(data ThirdOutObj) error {
	apiClient := client.NewAPIClient()
	// 构造请求数据
	_, err := apiClient.SendRequest("", data)
	return err
}
