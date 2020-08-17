package messages

type HitMessage struct {
	BaseMessage
	Coordinate Coordinate `json:"coordinate"`
}

func NewHitMessage(username string, coordinate Coordinate) HitMessage {
	return HitMessage{
		BaseMessage: BaseMessage{Username: username, Type: "HIT"},
		Coordinate:  coordinate,
	}
}

type MissMessage struct {
	BaseMessage
	Coordinate Coordinate `json:"coordinate"`
}

func NewMissMessage(username string, coordinate Coordinate) MissMessage {
	return MissMessage{
		BaseMessage: BaseMessage{Username: username, Type: "MISS"},
		Coordinate:  coordinate,
	}
}

type ShipDestroyedMessage struct {
	BaseMessage
	Coordinate Coordinate `json:"coordinate"`
	ShipSize   int        `json:"shipSize"`
	Vertical   bool       `json:"vertical"`
}

func NewShipDestroyedMessage(username string, coordinate Coordinate, shipSize int, vertical bool) ShipDestroyedMessage {
	return ShipDestroyedMessage{
		BaseMessage: BaseMessage{Username: username, Type: "SHIP_DESTROYED"},
		Coordinate:  coordinate,
		ShipSize:    shipSize,
		Vertical:    vertical,
	}
}

type OpponentDestroyedShip struct {
	BaseMessage
}

func NewOpponentDestroyedShipMessage(username string) OpponentDestroyedShip {
	return OpponentDestroyedShip{
		BaseMessage: BaseMessage{Username: username, Type: "OPPONENT_DESTROYED_SHIP"},
	}
}

type VictoryMessage struct {
	BaseMessage
}

func NewVictoryMessage(username string) VictoryMessage {
	return VictoryMessage{
		BaseMessage{Username: username, Type: "VICTORY"},
	}
}

type LossMessage struct {
	BaseMessage
}

func NewLossMessage(username string) LossMessage {
	return LossMessage{
		BaseMessage{Username: username, Type: "LOSS"},
	}
}

type BoardMessage struct {
	BaseMessage
	ShipSizes []int `json:"shipSizes"`
}

func NewBoardMessage(username string, shipSizes []int) BoardMessage {
	return BoardMessage{
		BaseMessage: BaseMessage{Username: username, Type: "BOARD"},
		ShipSizes:   shipSizes,
	}
}

type awaitingPlayersMessage struct {
	BaseMessage
}

func NewAwaitingPlayersMessage(username string) awaitingPlayersMessage {
	return awaitingPlayersMessage{
		BaseMessage{Username: username, Type: "AWAITING_PLAYERS"},
	}
}

type gameStartedMessage struct {
	BaseMessage
	TurnMessage
	Usernames []string `json:"usernames"`
}

func NewGameStartedMessage(username string, turn bool, duration int, usernames []string) gameStartedMessage {
	return gameStartedMessage{
		BaseMessage: BaseMessage{Username: username, Type: "GAME_STARTED"},
		TurnMessage: TurnMessage{Turn: turn, Duration: duration},
		Usernames:   usernames,
	}
}

type TurnMessage struct {
	BaseMessage
	Turn     bool `json:"turn"`
	Duration int  `json:"duration"`
}

func NewTurnMessage(username string, turn bool, duration int) TurnMessage {
	return TurnMessage{
		BaseMessage: BaseMessage{Username: username, Type: "TURN"},
		Turn:        turn,
		Duration:    duration,
	}
}

type TurnExtendedMessage struct {
	BaseMessage
	Turn     bool `json:"turn"`
	Duration int  `json:"duration"`
}

func NewTurnExtendedMessage(username string, duration int) TurnExtendedMessage {
	return TurnExtendedMessage{
		BaseMessage: BaseMessage{Username: username, Type: "TURN_EXTENDED"},
		Turn:        true,
		Duration:    duration,
	}
}

type boardStateMessage struct {
	BaseMessage
	Hits    []HitMessage           `json:"hits"`
	Misses  []MissMessage          `json:"misses"`
	Destoys []ShipDestroyedMessage `json:"destroys"`
	Board   BoardMessage           `json:"board"`
}

func NewBoardStateMessage(username string, hits []HitMessage, misses []MissMessage, destroys []ShipDestroyedMessage, board BoardMessage) boardStateMessage {
	return boardStateMessage{
		BaseMessage: BaseMessage{Username: username, Type: "BOARD_STATE"},
		Board:       board,
		Hits:        hits,
		Misses:      misses,
		Destoys:     destroys,
	}
}

type cancelledMessage struct {
	BaseMessage
}

func NewCancelledMessage() cancelledMessage {
	return cancelledMessage{
		BaseMessage{Type: "CANCELLED"},
	}
}
