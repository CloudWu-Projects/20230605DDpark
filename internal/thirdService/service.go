package thirdservice

import (
	"encoding/json"
	"fmt"
	"jilaidian_go/internal/config"
	"jilaidian_go/pkg/client"
)

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
type ThirdResponse struct {
	/*
	   {"code":0,"msg":"成功","data":""}
	*/
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func PostInData(data ThirdInObj) error {
	apiClient := client.NewHttpPostClient()
	// 构造请求数据
	body, err := apiClient.SendRequest(config.Global.ThridServer.In_url, data)
	if err != nil {
		return fmt.Errorf("PostInData 发送请求失败: %v data:%v", err, data)
	}
	var response ThirdResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("ThirdResponse 解析响应失败: %v %v", err, string(body))
	}
	if response.Code != 0 {
		return fmt.Errorf("ThirdResponse 返回失败: %v", response)
	}
	return err
}

func PostOutData(data ThirdOutObj) error {
	apiClient := client.NewHttpPostClient()
	// 构造请求数据
	body, err := apiClient.SendRequest(config.Global.ThridServer.Out_url, data)
	if err != nil {
		return fmt.Errorf("PostOutData 发送请求失败: %v data:%v", err, data)
	}
	var response ThirdResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("ThirdResponse 解析响应失败: %v %v", err, string(body))
	}
	if response.Code != 0 {
		return fmt.Errorf("ThirdResponse 返回失败: %v", response)
	}
	return err
}
