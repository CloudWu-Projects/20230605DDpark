package config

type YiAnqi struct {
	TokenURL       string `json:"tokenUrl" form:"YiAnqi.TokenURL"`
	OperatorID     string `json:"operatorID" form:"YiAnqi.OperatorID"`
	OperatorSecret string `json:"operatorSecret" form:"YiAnqi.OperatorSecret"`
	AesKey         string `json:"aeskey" form:"YiAnqi.AesKey"`
	AesIv          string `json:"aesiv" form:"YiAnqi.AesIv"`
	SignKey        string `json:"signKey" form:"YiAnqi.SignKey"`
}
