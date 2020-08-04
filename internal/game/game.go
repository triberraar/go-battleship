package game

import (
	"github.com/triberraar/go-battleship/internal/messages"
)

type Game interface {
	OutMessages() chan messages.GameMessage
	InMessages() chan []byte
	Rejoin()
}

type GameDefinition interface {
	TurnDuration() int
	GameName() string
}

type GameCreator interface {
	Game(username string) Game
	FromExisting(username string, game Game) (Game, error)
	GameDefinition(gameName string) GameDefinition
}
