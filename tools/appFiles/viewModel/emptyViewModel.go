package viewModel

import (
	"encoding/json"
)

type EmptyViewModel struct{}

func (this *EmptyViewModel) LoadDefaultState() {

}

func (self *EmptyViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
