package utils

import (
	"testing"
)

func TestSign(t *testing.T) {
	/*
		Sig（签名）采用HMAC-MD5算法，采用MD5作为散列函数，
		通过SigSecret（签名密钥）对整个消息主体各参数的值拼接后进行加密，
		入参拼接顺序为：OperatorID（运营商标识）、Data（参数内容）、TimeStamp（时间戳）、Seq（自增序列），
		出参拼接顺序为：Ret（返回值）、Msg（返回信息）、Data（参数内容），
		然后采用MD5信息摘要的方式形成新密文，参数签名必须大写，详见附录B。
	*/
	/*
	   示例签名密钥：1234567890abcdef
	   示例运营商标识（OperatorID）：123456789
	   示例参数信息（Data）： il7B0BSEjFdzpyKzfOFpvg/Se1CP802RItKYFPfSLRxJ3jf0bVl9hvYOEktPAYW2nd7S8MBcyHYyacHKbISq5iTmDzG+ivnR+SZJv3USNTYVMz9rCQVSxd0cLlqsJauko79NnwQJbzDTyLooYoIwz75qBOH2/xOMirpeEqRJrF/EQjWekJmGk9RtboXePu2rka+Xm51syBPhiXJAq0GfbfaFu9tNqs/e2Vjja/ltE1M0lqvxfXQ6da6HrThsm5id4ClZFIi0acRfrsPLRixS/IQYtksxghvJwbqOsbIsITail9Ayy4tKcogeEZiOO+4Ed264NSKmk7l3wKwJLAFjCFogBx8GE3OBz4pqcAn/ydA=
	   示例时间戳（TimeStamp）：20160729142400
	   示例自增序列（Seq）：0001
	   示例签名（Sig）：745166E8C43C84D37FFEC0F529C4136F
	*/
	// 示例数据
	ukey := "1234567890abcdef"
	operatorID := "123456789"
	data := "il7B0BSEjFdzpyKzfOFpvg/Se1CP802RItKYFPfSLRxJ3jf0bVl9hvYOEktPAYW2nd7S8MBcyHYyacHKbISq5iTmDzG+ivnR+SZJv3USNTYVMz9rCQVSxd0cLlqsJauko79NnwQJbzDTyLooYoIwz75qBOH2/xOMirpeEqRJrF/EQjWekJmGk9RtboXePu2rka+Xm51syBPhiXJAq0GfbfaFu9tNqs/e2Vjja/ltE1M0lqvxfXQ6da6HrThsm5id4ClZFIi0acRfrsPLRixS/IQYtksxghvJwbqOsbIsITail9Ayy4tKcogeEZiOO+4Ed264NSKmk7l3wKwJLAFjCFogBx8GE3OBz4pqcAn/ydA="
	timeStamp := "20160729142400"
	seq := "0001"
	want_sig := "745166E8C43C84D37FFEC0F529C4136F"

	prestr := operatorID + data + timeStamp + seq
	t.Log("prestr:", prestr)
	// 生成签名
	sign := GenerateSignString(prestr, ukey)

	t.Log("data sign:", sign)
	t.Log("want sign:", want_sig)
	if sign != want_sig {
		t.Error("sign error")
	}
}
