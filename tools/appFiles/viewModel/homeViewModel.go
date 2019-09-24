package viewModel

import (
	"encoding/json"
)

type HomeViewModel struct {
	ReleaseDescriptionLines   []string `json:"ReleaseDescriptionLines"`
	ReleaseNotes              string   `json:"ReleaseNotes"`
	UserCount                 int      `json:"UserCount"`
	WebSocketConnectionsCount int      `json:"WebSocketConnectionsCount"`
}

func (self *HomeViewModel) LoadDefaultState() {
}

func (self *HomeViewModel) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
