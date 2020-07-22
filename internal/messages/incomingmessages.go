package messages

// FireMessage sent by player when he fires somewhere
type FireMessage struct {
	BaseMessage
	Coordinate Coordinate `json:"coordinate"`
}

type PlayMessage struct {
	BaseMessage
}
