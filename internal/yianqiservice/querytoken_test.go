package yianqiservice

import (
	"testing"
	"time"
)

func TestQueryToken(t *testing.T) {
	q := NewTokenMgr()
	_, v := q.GetValidToken()
	if v.isValidToken("") {
		t.Log("token有效")
	} else {
		t.Error("token无效")
	}
	time.Sleep(time.Hour)
}
