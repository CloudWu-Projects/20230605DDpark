package client

import (
	"testing"
)

func TestClientQueryOrder(t *testing.T) {

	parkID := 10051557
	carNumber := "鲁BFZ7606"

	nc := NewAPIClient()
	orderId, err := nc.QueryOrder(parkID, carNumber)
	if err != nil {
		t.Errorf("Error querying order:\n %v", err)
		return
	}
	//t.Errorf("Query order success orderId %s %v", orderId, err)

	err = nc.SendDiscountNotice(parkID, carNumber, orderId, 8.0, 4, 5)
	if err != nil {
		t.Errorf("Error sending discount notice:\n %v", err)
		return
	}
}
func TestClientSendDiscountNotice(t *testing.T) {

}
