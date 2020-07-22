package game

import (
	"errors"

	"github.com/triberraar/go-battleship/internal/messages"

	"github.com/triberraar/go-battleship/internal/game/battleship"
)

type Game interface {
	GetOutMessages() chan messages.GameMessage
	GetInMessages() chan []byte
}

func NewGame(gameName string, playerID string) (Game, error) {
	switch gameName {
	case "battleship":
		return battleship.NewBattleship(playerID), nil
	}
	return nil, errors.New("unknown game")
}

func NewGameFromExistion(game Game, playerID string) (Game, error) {
	switch g := game.(type) {
	case *battleship.Battleship:
		return g.NewBattleshipFromExisting(playerID), nil
	default:
		return nil, errors.New("unknown game")
	}
}
