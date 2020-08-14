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

// func (rm *RoomManager) JoinRoom(client *client.Client, gameName string) {
// 	rm.joinMutex.Lock()
// 	for _, room := range rm.rooms {
// 		if room.HasPlayer(client.Username) {
// 			room.rejoinPlayer(client)
// 			rm.joinMutex.Unlock()
// 			return
// 		}
// 	}
// 	if rm.waitingRoom == nil || rm.waitingRoom.isFinished() {
// 		log.Println("Creating new room")
// 		room := NewRoom(2, gameName)
// 		rm.rooms[room.id] = room
// 		rm.waitingRoom = room
// 		go room.Run()
// 		go func(c chan bool) {
// 			<-c
// 			log.Println("removing finished room")
// 			rm.joinMutex.Lock()
// 			delete(rm.rooms, room.id)
// 			rm.joinMutex.Unlock()
// 		}(room.removeMe)
// 	}
// 	rm.waitingRoom.joinPlayer(client)
// 	rm.joinMutex.Unlock()
// }
