package payloads

import "encoding/json"

type ApiRequest struct {
	Action     string `json:"action" js:"action"`
	State      string `json:"state" js:"state"`
	Controller string `json:"controller"`
}

type ApiResponse struct {
	Redirect          string `json:"Redirect"`
	GlobalMessage     string `json:"GlobalMessage"`
	GlobalMessageType string `json:"GlobalMessageType"`
	Transactionid     string `json:"TransactionId"`
	Trace             string `json:"Trace"`
	State             string `json:"State"`
}

func (self *ApiResponse) Stringify() (value string, err error) {
	data, err := json.Marshal(self)
	value = string(data)
	return
}
