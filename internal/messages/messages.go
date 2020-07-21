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
	BaseMessage
	Coordinate Coordinate `json:"coordinate"`
}

// MissMessage send to player when he misses stuff
type missMessage struct {
	BaseMessage
	Coordinate Coordinate `json:"coordinate"`
}

type shipDestroyedMessage struct {
	BaseMessage
	Coordinate Coordinate `json:"coordinate"`
	ShipSize   int        `json:"shipSize"`
	Vertical   bool       `json:"vertical"`
}

type VictoryMessage struct {
	BaseMessage
}

type LossMessage struct {
	BaseMessage
}

type boardMessage struct {
	BaseMessage
	ShipSizes []int `json:"shipSizes"`
}

// FireMessage sent by player when he fires somewhere
type FireMessage struct {
	BaseMessage
	Coordinate Coordinate `json:"coordinate"`
}

type PlayMessage struct {
	BaseMessage
}

// NewHitMessage constructor function
func NewHitMessage(coordinate Coordinate) hitMessage {
	return hitMessage{
		BaseMessage: BaseMessage{Type: "HIT"},
		Coordinate:  coordinate,
	}
}

// NewMissMessage constructor function
func NewMissMessage(coordinate Coordinate) missMessage {
	return missMessage{
		BaseMessage: BaseMessage{Type: "MISS"},
		Coordinate:  coordinate,
	}
}

func NewShipDestroyedMessage(coordinate Coordinate, shipSize int, vertical bool) shipDestroyedMessage {
	return shipDestroyedMessage{
		BaseMessage: BaseMessage{Type: "SHIP_DESTROYED"},
		Coordinate:  coordinate,
		ShipSize:    shipSize,
		Vertical:    vertical,
	}
}

func NewVictoryMessage() VictoryMessage {
	return VictoryMessage{
		BaseMessage{Type: "VICTORY"},
	}
}

func NewLossMessage() LossMessage {
	return LossMessage{
		BaseMessage{Type: "LOSS"},
	}
}

func NewBoardMessage(shipSizes []int) boardMessage {
	return boardMessage{
		BaseMessage: BaseMessage{Type: "BOARD"},
		ShipSizes:   shipSizes,
	}
}

type awaitingPlayersMessage struct {
	BaseMessage
}

func NewAwaitingPlayersMessage() awaitingPlayersMessage {
	return awaitingPlayersMessage{
		BaseMessage{Type: "AWAITING_PLAYERS"},
	}
}

type gameStartedMessage struct {
	BaseMessage
	Turn bool `json:"turn"`
}

func NewGameStartedMessage(turn bool) gameStartedMessage {
	return gameStartedMessage{
		BaseMessage: BaseMessage{Type: "GAME_STARTED"},
		Turn:        turn,
	}
}

type TurnMessage struct {
	BaseMessage
	Turn bool `json:"turn"`
}

func NewTurnMessage(turn bool) TurnMessage {
	return TurnMessage{
		BaseMessage: BaseMessage{Type: "TURN"},
		Turn:        turn,
	}
}
