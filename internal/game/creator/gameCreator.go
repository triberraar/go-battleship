package creator

import (
	"errors"

	"github.com/triberraar/go-battleship/internal/game"
	"github.com/triberraar/go-battleship/internal/game/battleship"
)

type GameCreator interface {
	game(playerID string) game.Game
	fromExisting(playerID string, game game.Game) (game.Game, error)
	gameDefinition(gameName string) game.GameDefinition
}

type BattleshipGameCreator struct {
}

func (bgc BattleshipGameCreator) game(playerID string) game.Game {
	return battleship.NewBattleship(playerID)
}

func (bgc BattleshipGameCreator) gameDefinition(gameName string) game.GameDefinition {
	return battleship.NewGameDefinition(gameName)
}

func (bgc BattleshipGameCreator) fromExisting(playerID string, game game.Game) (game.Game, error) {
	bs, ok := game.(*battleship.Battleship)
	if ok {
		return bs.NewBattleshipFromExisting(playerID), nil
	}
	return nil, errors.New("unknown game")

}

var gameCreators = map[string]GameCreator{"battleships": BattleshipGameCreator{}}

func NewGame(gameName string, playerID string) (game.Game, error) {
	if val, ok := gameCreators[gameName]; ok {
		return val.game(playerID), nil
	}
	return nil, errors.New("unknown game")
}

func NewGameFromExistion(gameName string, game game.Game, playerID string) (game.Game, error) {
	if val, ok := gameCreators[gameName]; ok {
		return val.fromExisting(playerID, game)
	}
	return nil, errors.New("unknown game")

}

func NewGameDefinition(gameName string) (game.GameDefinition, error) {
	if val, ok := gameCreators[gameName]; ok {
		return val.gameDefinition(gameName), nil
	}
	return nil, errors.New("unknown game")
}
