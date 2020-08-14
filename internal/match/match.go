package match

import (
	"errors"

	"github.com/google/uuid"
	"github.com/triberraar/go-battleship/internal/client"
	"github.com/triberraar/go-battleship/internal/game/battleship"
)

type Match interface {
	ShouldRejoin(username string) bool
	Join(client *client.Client)
	Rejoin(client *client.Client)
	GetID() uuid.UUID
}

func NewMatch(gameName string) (Match, error) {
	switch gameName {
	case "battleships":
		return battleship.NewBattleshipMatch(2), nil
	}
	return nil, errors.New("unknown game")
}
