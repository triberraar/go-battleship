package room

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/triberraar/go-battleship/internal/client"
)

type RoomManager struct {
	rooms       map[uuid.UUID]*Room
	waitingRoom *Room
	joinMutex   sync.Mutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{rooms: make(map[uuid.UUID]*Room), waitingRoom: nil}
}

func (rm RoomManager) String() string {
	return fmt.Sprintf("I be the room manager, here be my rooms %v", rm.rooms)
}

func (rm *RoomManager) JoinRoom(client *client.Client, gameName string) {
	rm.joinMutex.Lock()
	for _, room := range rm.rooms {
		if room.HasPlayer(client.Username) {
			room.rejoinPlayer(client)
			rm.joinMutex.Unlock()
			return
		}
	}
	if rm.waitingRoom == nil || rm.waitingRoom.isFinished() {
		log.Println("Creating new room")
		room := NewRoom(2, gameName)
		rm.rooms[room.id] = room
		rm.waitingRoom = room
		go room.Run()
		go func(c chan bool) {
			<-c
			log.Println("removing finished room")
			rm.joinMutex.Lock()
			delete(rm.rooms, room.id)
			rm.joinMutex.Unlock()
		}(room.removeMe)
	}
	rm.waitingRoom.joinPlayer(client)
	rm.joinMutex.Unlock()
}
