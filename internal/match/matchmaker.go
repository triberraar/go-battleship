package match

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/triberraar/go-battleship/internal/client"
)

type Matchmaker struct {
	matches map[string]*matches
}

type matches struct {
	matches       map[uuid.UUID]Match
	awaitingMatch Match
	mutex         *sync.Mutex
}

func NewMatchmaker(gameNames []string) *Matchmaker {
	matchesm := make(map[string]*matches)
	for _, gameName := range gameNames {
		matchesm[gameName] = &matches{make(map[uuid.UUID]Match), nil, &sync.Mutex{}}
	}
	return &Matchmaker{matchesm}
}

func (m *Matchmaker) Play(client *client.Client, gameName string) error {
	if _, ok := m.matches[gameName]; !ok {
		return errors.New("unknown game")
	}

	gameMatch := m.matches[gameName]
	gameMatch.mutex.Lock()

	for _, match := range gameMatch.matches {
		if match.ShouldRejoin(client.Username) {
			match.Rejoin(client)
			gameMatch.mutex.Unlock()
			return nil
		}
	}
	if gameMatch.awaitingMatch == nil {
		newMatch, _ := NewMatch(gameName)
		gameMatch.awaitingMatch = newMatch
		gameMatch.matches[newMatch.GetID()] = newMatch
	}
	gameMatch.awaitingMatch.Join(client)

	gameMatch.mutex.Unlock()

	return nil
}
