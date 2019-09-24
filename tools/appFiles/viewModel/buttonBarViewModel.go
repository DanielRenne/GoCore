package viewModel

type ButtonBar struct {
	Config struct {
		CurrentTab       string            `json:"CurrentTab"`
		VisibleTabs      map[string]string `json:"VisibleTabs"` // key references only which will be used for translations and identifying the tab ID only
		TabOrder         []string          `json:"TabOrder"`
		TabActions       []string          `json:"TabActions"`
		TabControllers   []string          `json:"TabControllers"`
		TabIsVisible     []bool            `json:"TabIsVisible"`
		OtherTabSelected []string          `json:"OtherTabSelected"`
	} `json:"Config"`
}

func (this *ButtonBar) AddTab(key string, action string, controller string, order string, visible bool, otherTabSelected string) {

	if this.Config.VisibleTabs == nil {
		this.Config.VisibleTabs = make(map[string]string)
	}

	this.Config.VisibleTabs[key] = key
	this.Config.TabActions = append(this.Config.TabActions, action)
	this.Config.TabControllers = append(this.Config.TabControllers, controller)
	this.Config.TabOrder = append(this.Config.TabOrder, key)
	this.Config.TabIsVisible = append(this.Config.TabIsVisible, visible)
	this.Config.OtherTabSelected = append(this.Config.OtherTabSelected, otherTabSelected)
}
