package room

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/triberraar/go-battleship/internal/client"
	"github.com/triberraar/go-battleship/internal/messages"
)

type RoomManager struct {
	rooms     []*Room
	joinMutex sync.Mutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{}
}

func (rm RoomManager) String() string {
	return fmt.Sprintf("I be the room manager, here be my rooms %v", rm.rooms)
}

func (rm *RoomManager) JoinRoom(client *client.Client) {
	rm.joinMutex.Lock()
	if len(rm.rooms) == 0 {
		log.Println("making first room")
		room := NewRoom(2, client)
		rm.rooms = append(rm.rooms, room)
		go room.Run()
	} else {
		if rm.rooms[len(rm.rooms)-1].isFull() {
			log.Println("room is full, making new room")
			room := NewRoom(2, client)
			rm.rooms = append(rm.rooms, room)
			go room.Run()
		} else {
			log.Println("joining room")
			rm.rooms[len(rm.rooms)-1].joinPlayer(client)
		}
	}
	if rm.rooms[len(rm.rooms)-1].isFull() {
		for _, pl := range rm.rooms[len(rm.rooms)-1].players {
			rm.rooms[len(rm.rooms)-1].players[pl.playerID].game.SendMessage(messages.NewGameStartedMessage(pl.playerID == rm.rooms[len(rm.rooms)-1].currentPlayer))
		}
	} else {
		for _, pl := range rm.rooms[len(rm.rooms)-1].players {
			rm.rooms[len(rm.rooms)-1].players[pl.playerID].game.SendMessage(messages.NewAwaitingPlayersMessage())
		}
	}
	rm.joinMutex.Unlock()
}

func (rm *RoomManager) Run() {
	ticker := time.NewTicker(60 * time.Second)
	for {
		<-ticker.C
		log.Println(rm)
	}
}
