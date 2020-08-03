package creator

import (
	"errors"

	"github.com/triberraar/go-battleship/internal/game"
	"github.com/triberraar/go-battleship/internal/game/battleship"
)

var gameCreators = map[string]game.GameCreator{"battleships": battleship.BattleshipGameCreator{}}

func NewGame(gameName string, connectionID string) (game.Game, error) {
	if val, ok := gameCreators[gameName]; ok {
		return val.Game(connectionID), nil
	}
	return nil, errors.New("unknown game")
}

func NewGameFromExistion(gameName string, game game.Game, connectionID string) (game.Game, error) {
	if val, ok := gameCreators[gameName]; ok {
		return val.FromExisting(connectionID, game)
	}
	return nil, errors.New("unknown game")

}

func NewGameDefinition(gameName string) (game.GameDefinition, error) {
	if val, ok := gameCreators[gameName]; ok {
		return val.GameDefinition(gameName), nil
	}
	return nil, errors.New("unknown game")
}
