package yianqiservice

import (
	"encoding/json"
	"jilaidian_go/internal/config"
	"jilaidian_go/pkg/client"
	"jilaidian_go/pkg/logger"
	"time"

	. "jilaidian_go/internal/utils"
)

type TokenMgr struct {
	validTokensA ValidToken

	apiClient *client.HttpPostClient
}

var tokenMgr TokenMgr

func NewTokenMgr() *TokenMgr {
	return &tokenMgr

}
func init() {
	// start a thread to update the token every hour

	tm := NewTokenMgr()
	tm.apiClient = client.NewHttpPostClient()
	tm.validTokensA.ExpirationTime = 0
	go func() {
		tm := NewTokenMgr()
		for {
			if !tm.validTokensA.IsValid() {

				err := tm.GetValidToken()
				if err != nil {
					logger.Logger.Errorf("获取token失败: %v", err)
				} else {
					//	validTokensA.update()
					logger.Logger.Infof("获取到新的token: %s, 过期时间: %d",
						tm.validTokensA.AccessToken, tm.validTokensA.ExpirationTime)
				}
			} else {
				logger.Logger.Infof("token未过期: %s, 过期时间: %d",
					tm.validTokensA.AccessToken, tm.validTokensA.ExpirationTime)
			}
			time.Sleep(time.Minute) // 每小时更新一次
		}
	}()

}

func (tm *TokenMgr) GetValidToken() error {

	q := QueryTokenRequest{
		OperatorID:     config.Global.YiAnqi.OperatorID,
		OperatorSecret: config.Global.YiAnqi.OperatorSecret,
	}
	request := MakeRequest(q)
	body, err := tm.apiClient.SendRequest(config.Global.YiAnqi.TokenURL, request)
	if err != nil {
		logger.Logger.Errorf("获取token失败: %v", err)
		return err
	}
	TotalResponse := TotalResponse{}
	err = json.Unmarshal(body, &TotalResponse)

	decodedStr, err := CBCDecrypt_Base64(string(TotalResponse.Data), config.Global.YiAnqi.AesKey,
		config.Global.YiAnqi.AesIv)

	if err != nil {
		logger.Logger.Errorf("解密失败 %s %v", body, err)
		return err
	}

	err = json.Unmarshal([]byte(decodedStr), &tm.validTokensA)
	if err != nil {
		logger.Logger.Errorf("解析失败 %s %v", decodedStr, err)
		return err
	}
	tm.validTokensA.ExpirationTime = time.Now().Add(time.Second * time.Duration(tm.validTokensA.TokenAvailableTime)).Unix() // 1小时后过期
	return err
}
