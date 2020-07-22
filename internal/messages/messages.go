package messages

// BaseMessage is a generic base type message
type BaseMessage struct {
	Type string `json:"type"`
}

// Coordinate represents a coordinate on the battleship board
type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}
