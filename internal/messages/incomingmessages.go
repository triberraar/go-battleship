package messages

type FireMessage struct {
	BaseMessage
	Coordinate Coordinate `json:"coordinate"`
}

type PlayMessage struct {
	BaseMessage
	Username string `json:"username"`
}
