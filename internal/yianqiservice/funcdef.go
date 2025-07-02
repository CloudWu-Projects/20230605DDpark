package yianqiservice

import (
	"encoding/json"
	"fmt"
	"jilaidian_go/internal/config"
	"jilaidian_go/pkg/logger"
	"net/http"
	"strconv"
	"time"

	. "jilaidian_go/internal/utils"

	"github.com/gin-gonic/gin"
)

func (tr *TotalResponse) MakeSig() {
	/*
		Sig（签名）采用HMAC-MD5算法，采用MD5作为散列函数，
		通过SigSecret（签名密钥）对整个消息主体各参数的值拼接后进行加密，
		入参拼接顺序为：OperatorID（运营商标识）、Data（参数内容）、TimeStamp（时间戳）、Seq（自增序列），
		出参拼接顺序为：Ret（返回值）、Msg（返回信息）、Data（参数内容），
		然后采用MD5信息摘要的方式形成新密文，参数签名必须大写，详见附录B。
	*/
	prestr := fmt.Sprintf("%d%s%s", tr.Ret, tr.Msg, tr.Data)

	tr.Sig = GenerateSignString(prestr, config.Global.YiAnqi.SignKey)
}
func (tr *YiAnqiRequest) MakeSig() {
	/*
		Sig（签名）采用HMAC-MD5算法，采用MD5作为散列函数，
		通过SigSecret（签名密钥）对整个消息主体各参数的值拼接后进行加密，
		入参拼接顺序为：OperatorID（运营商标识）、Data（参数内容）、TimeStamp（时间戳）、Seq（自增序列），
		出参拼接顺序为：Ret（返回值）、Msg（返回信息）、Data（参数内容），
		然后采用MD5信息摘要的方式形成新密文，参数签名必须大写，详见附录B。
	*/
	prestr := fmt.Sprintf("%s%s%s%s", tr.OperatorID, tr.Data, tr.TimeStamp, tr.Seq)

	tr.Sig = GenerateSignString(prestr, config.Global.YiAnqi.SignKey)
}

var seq int

func init() {
	seq = 0
}

func GetSeq() string {
	seq++
	return strconv.Itoa(seq)
}

func MakeRequest(data interface{}) YiAnqiRequest {
	var jsonData []byte
	switch v := data.(type) {
	case string:
		jsonData = []byte(v)
	default:
		jsonData, _ = json.Marshal(data)
	}

	encodedStr, _ := CBCEncrypt_Base64(string(jsonData), config.Global.YiAnqi.AesKey, config.Global.YiAnqi.AesIv)

	tr := YiAnqiRequest{
		OperatorID: config.Global.YiAnqi.OperatorID,
		Data:       encodedStr,
		TimeStamp:  time.Now().Format("20060102150405"),
		Seq:        GetSeq(),
	}
	tr.MakeSig()
	return tr
}

func (h *Handler) MakeRepsonse(c *gin.Context, result int, description string, data interface{}) {

	var jsonData []byte
	switch v := data.(type) {
	case string:
		jsonData = []byte(v)
	default:
		jsonData, _ = json.Marshal(data)
	}

	logger.Logger.Debug("yianqi MakeRepsonse Data :", data)
	encodedStr, _ := CBCEncrypt_Base64(string(jsonData), config.Global.YiAnqi.AesKey, config.Global.YiAnqi.AesIv)

	tr := TotalResponse{
		Ret:  result,
		Msg:  description,
		Data: encodedStr,
		Sig:  "",
	}
	tr.MakeSig()
	logger.Logger.Debugf("yianqi MakeRepsonse jsonData : %s", string(jsonData))
	logger.Logger.Debug("yianqi MakeRepsonse TotalResponse : ", tr)

	c.JSON(http.StatusOK, tr)
}
