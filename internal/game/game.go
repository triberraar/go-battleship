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

type GameCreator interface {
	Game(playerID string) Game
	FromExisting(playerID string, game Game) (Game, error)
	GameDefinition(gameName string) GameDefinition
}
