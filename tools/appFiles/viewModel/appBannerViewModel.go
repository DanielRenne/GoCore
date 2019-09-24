package viewModel

type Banner struct {
	Color              string `json:"Color"`
	HistoryLength      int    `json:"HistoryLength"`
	AccountName        string `json:"AccountName"`
	IsSecondaryAccount bool   `json:"IsSecondaryAccount"`
	CurrentHistoryIdx  int    `json:"CurrentHistoryIdx"`
	LatestHistoryLen   int    `json:"LatestHistoryLen"`
}

func (self *Banner) LoadDefaultState() {

}
