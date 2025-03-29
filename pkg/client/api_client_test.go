package client

import (
	"jilaidian_go/internal/config"
	"testing"
)

func TestClientQueryOrder(t *testing.T) {

	parkID := 10051557
	carNumber := "鲁BFZ7606"
	parkInfo := config.ParkInfo{
		DeductionMoney: 100,
		DeductionTime:  60,
		ParkID:         10051557,
		Ukey:           "your_ukey",
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
