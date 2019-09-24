package payloads

import "encoding/json"

type MarkupMiddleware struct {
	Html              string `xml:"html" json:"Html"`
	Json              string `xml:"json" json:"Json"`
	PageContent       string `xml:"pageContent" json:"PageContent"`
	Redirect          string `xml:"redirect" json:"Redirect"`
	GlobalMessage     string `xml:"globalMessage" json:"GlobalMessage"`
	Trace             string `xml:"trace" json:"Trace"`
	GlobalMessageType string `xml:"globalMessageType" json:"GlobalMessageType"`
}

func (self *MarkupMiddleware) Stringify() (value string, err error) {
	data, err := json.Marshal(self)
	value = string(data)
	return
}
