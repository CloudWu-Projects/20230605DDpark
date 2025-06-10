package yianqiservice

type Message struct {
	Result      int    `json:"result"`
	Description string `json:"description"`
}

type QueryTokenRequest struct {
	OperatorID     string `json:"OperatorID"`
	OperatorSecret string `json:"OperatorSecret"`
}
type QueryTokenResponse struct {
	OperatorID         string `json:"OperatorID"`
	SuccStat           int    `json:"SuccStat"`
	AccessToken        string `json:"AccessToken"`
	TokenAvailableTime int    `json:"TokenAvailableTime"`
	FailReason         int    `json:"FailReason"`
}
type TotalResponse struct {
	Ret  int    `json:"Ret"`
	Msg  string `json:"Msg"`
	Data string `json:"Data"`
	Sig  string `json:"Sig"`
}
type YiAnqiRequest struct {
	OperatorID string `json:"OperatorID"`
	Data       string `json:"Data"`
	TimeStamp  string `json:"TimeStamp"`
	Seq        string `json:"Seq"`
	Sig        string `json:"Sig"`
}
