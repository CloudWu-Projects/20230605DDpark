package yianqiservice

import (
	"testing"
	"time"
)

func TestQueryToken(t *testing.T) {
	q := NewTokenMgr()
	v := q.GetValidToken()
	if v != nil {
		t.Errorf("获取token失败: %v", v)
	}
	if q.ValidToken.IsValid() {
		t.Log("token有效")
	} else {
		t.Error("token无效")
	}
	time.Sleep(time.Hour)
}
