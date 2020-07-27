package creator

import (
	"errors"

	"github.com/triberraar/go-battleship/internal/game"
	"github.com/triberraar/go-battleship/internal/game/battleship"
)

var gameCreators = map[string]game.GameCreator{"battleships": battleship.BattleshipGameCreator{}}

func NewGame(gameName string, playerID string) (game.Game, error) {
	if val, ok := gameCreators[gameName]; ok {
		return val.Game(playerID), nil
	}
	return nil, errors.New("unknown game")
}

func NewGameFromExistion(gameName string, game game.Game, playerID string) (game.Game, error) {
	if val, ok := gameCreators[gameName]; ok {
		return val.FromExisting(playerID, game)
	}
	return nil, errors.New("unknown game")

}

func NewGameDefinition(gameName string) (game.GameDefinition, error) {
	if val, ok := gameCreators[gameName]; ok {
		return val.GameDefinition(gameName), nil
	}
	return nil, errors.New("unknown game")
}
