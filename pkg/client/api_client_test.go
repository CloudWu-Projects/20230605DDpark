package client

import (
	"jilaidian_go/internal/config"
	"testing"
)

func TestClientQueryOrder(t *testing.T) {

	parkID := 10051834
	carNumber := "京A11116"
	parkInfo := config.ParkInfo{
		DeductionMoney: 100,
		DeductionTime:  60,
		ParkID:         10051834,
		Ukey:           "J445O54V3REDI0NT",
	}
	nc := NewAPIClient()
	orderId, err := nc.QueryOrder(parkID, carNumber, &parkInfo)
	if err != nil {
		t.Errorf("Error querying order:\n %v", err)
		return
	}
	//t.Errorf("Query order success orderId %s %v", orderId, err)

	err = nc.SendDiscountNotice(parkID, carNumber, orderId, &parkInfo)
	if err != nil {
		t.Errorf("Error sending discount notice:\n %v", err)
		return
	}
}
func TestClientSendDiscountNotice(t *testing.T) {

}
