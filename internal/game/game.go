package game

import (
	"github.com/triberraar/go-battleship/internal/messages"
)

type Game interface {
	OutMessages() chan messages.GameMessage
	InMessages() chan []byte
}

type GameDefinition interface {
	TurnDuration() int
	GameName() string
}
