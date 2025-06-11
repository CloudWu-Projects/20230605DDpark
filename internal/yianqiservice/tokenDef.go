package yianqiservice

import (
	"jilaidian_go/internal/config"
	"jilaidian_go/internal/utils"
	"jilaidian_go/pkg/logger"
	"time"
)

type ValidToken struct {
	OperatorID         string `json:"OperatorID"`
	SuccStat           int    `json:"SuccStat"` // 成功状态
	AccessToken        string `json:"AccessToken"`
	TokenAvailableTime int    `json:"TokenAvailableTime"`
	FailReason         int    `json:"FailReason"` // 失败原因
	ExpirationTime     int64  `json:"UpdateTime"` // 更新时间
}

func (vt *ValidToken) update() {
	if vt.ExpirationTime < time.Now().Unix() {
		vt.AccessToken = utils.GenerateSignString(time.Now().Format("2006-01-02 15:04:05"), config.Global.YiAnqi.SignKey)
		vt.ExpirationTime = time.Now().Add(time.Hour).Unix() // 1小时后过期
		vt.TokenAvailableTime = 3600                         // 可用时间，单位秒
	}
	logger.Logger.Info("Updated valid token:", vt.AccessToken, " Expiration Time:", vt.ExpirationTime, " Available Time:", vt.TokenAvailableTime)
}
func (vt *ValidToken) isValidToken(token string) bool {
	return vt.AccessToken == token && vt.ExpirationTime > time.Now().Unix()
}
func (vt *ValidToken) IsValid() bool {
	return vt.AccessToken != "" && vt.ExpirationTime > time.Now().Unix()
}
