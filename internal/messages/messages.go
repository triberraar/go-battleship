package messages

type BaseMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
}

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}
