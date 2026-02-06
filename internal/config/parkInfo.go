package config

import (
	"fmt"
	"os"
)

type ParkInfo struct {
	ParkID         int    `json:"parkid"`     // 改为驼峰式
	StationID      string `json:"station_id"` // 改为驼峰式
	Ukey           string `json:"ukey"`
	DeductionTime  int    `json:"deduction_time"`  // 改为驼峰式
	DeductionMoney int    `json:"deduction_money"` // 改为驼峰式
	ReduceAmount   int    `json:"reduceAmount"`
	Duration       int    `json:"Duration"`
	Remark         string `json:"remark"`
	NpcPort        string `json:"npc_port"`
}

// 添加一个方法用于获取解密后的ukey
func (p *ParkInfo) GetUkey() string {
	// 这里可以添加解密逻辑，或从环境变量获取
	// 简单示例：如果环境变量中有对应的ukey，则使用环境变量中的值
	envKey := os.Getenv(fmt.Sprintf("PARK_UKEY_%d", p.ParkID))
	if envKey != "" {
		return envKey
	}
	return p.Ukey
}
