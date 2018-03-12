package socketViews

import "encoding/json"

type BroadcastLog struct {
	Id   string `json:"Id"`
	Data string `json"Data"`
}

func (self *BroadcastLog) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}

func (self *BroadcastLog) Stringify() (value string, err error) {
	data, err := json.Marshal(self)
	value = string(data)
	return
}

type BroadcastTime struct {
	Date string `json:"Date"`
	Time string `json"Time"`
}

func (self *BroadcastTime) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}

func (self *BroadcastTime) Stringify() (value string, err error) {
	data, err := json.Marshal(self)
	value = string(data)
	return
}
