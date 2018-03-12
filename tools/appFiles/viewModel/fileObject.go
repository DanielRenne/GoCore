package viewModel

import (
	"encoding/json"
	"time"
)

type FileObject struct {
	Name         string    `json:"Name"`
	Content      string    `json:"Content"`
	Size         int       `json:"Size"`
	Type         string    `json:"Type"`
	ModifiedUnix int       `json:"ModifiedUnix"`
	Modified     time.Time `json:"Modified"`
	Meta         struct {
		CompleteFailure    bool     `json:"CompleteFailure"`
		RowsSkipped        int      `json:"RowsSkipped"`
		RowsSkippedInfo    string   `json:"RowsSkippedInfo"`
		RowsCommitted      int      `json:"RowsCommitted"`
		RowsCommittedInfo  string   `json:"RowsCommittedInfo"`
		RowsSkippedDetails []string `json:"RowsSkippedDetails"`
		FileErrors         []string `json:"FileErrors"`
	} `json:"Meta"`
}

func (self *FileObject) LoadDefaultState() {
}

func (self *FileObject) Parse(data string) {
	json.Unmarshal([]byte(data), &self)
}
