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
type hitMessage struct {
	Type       string     `json:"type"`
	Coordinate Coordinate `json:"coordinate"`
}

// MissMessage send to player when he misses stuff
type missMessage struct {
	Type       string     `json:"type"`
	Coordinate Coordinate `json:"coordinate"`
}

// FireMessage sent by player when he fires somewhere
type FireMessage struct {
	Type       string     `json:"type"`
	Coordinate Coordinate `json:"coordinate"`
}

// NewHitMessage constructor function
func NewHitMessage(coordinate Coordinate) hitMessage {
	return hitMessage{
		Type:       "HIT",
		Coordinate: coordinate,
	}
}

// NewMissMessage constructor function
func NewMissMessage(coordinate Coordinate) missMessage {
	return missMessage{
		Type:       "MISS",
		Coordinate: coordinate,
	}
}
