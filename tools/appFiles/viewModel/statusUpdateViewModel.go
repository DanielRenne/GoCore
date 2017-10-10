package viewModel

import "encoding/json"

type StatusUpdateViewModel struct {
	Error    bool     `json:"Error"`
	Message  string   `json:"Message"`
	Message2 string   `json:"Message2"`
	Info     []string `json:"Info"`
	Mode     string   `json:"Mode"`
}

func (this *StatusUpdateViewModel) LoadDefaultState() {

}

func (self *StatusUpdateViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
