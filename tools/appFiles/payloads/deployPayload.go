package payloads

type DeployPayload struct {
	Type    int    `json:"type"`
	Payload string `json:"payload"`
}

type DeployEmail struct {
	Type    int    `json:"type"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}
