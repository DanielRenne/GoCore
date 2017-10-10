package viewModel

import (
	"encoding/json"
)

type LogsViewModel struct {
	Constants struct {
	} `json:"constants"`
	//AdditionalConstructs
	Id       string `json:"Id"`
	LongName string `json:"LongName"`
}

func (this *LogsViewModel) LoadDefaultState() {
	setConstants(this, "LOGS_CONST")
	this.LongName = ""
}

func (self *LogsViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
