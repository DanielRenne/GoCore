package socketViews

import "encoding/json"

type ClientStatus struct {
	Page string `json:"Page"`
}

func (self *ClientStatus) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
