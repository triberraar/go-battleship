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

// HitMessage send to player when he hits stuff
type HitMessage struct {
	Type       string     `json:"type"`
	Coordinate Coordinate `json:"coordinate"`
}

// MissMessage send to player when he misses stuff
type MissMessage struct {
	Type       string     `json:"type"`
	Coordinate Coordinate `json:"coordinate"`
}

// FireMessage sent by player when he fires somewhere
type FireMessage struct {
	Type       string     `json:"type"`
	Coordinate Coordinate `json:"coordinate"`
}
