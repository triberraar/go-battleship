package room

import (
	"fmt"
	"log"
	"sync"

	"github.com/triberraar/go-battleship/internal/client"
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

func (rm *RoomManager) JoinRoom(client *client.Client, gameName string) {
	rm.joinMutex.Lock()
	for _, room := range rm.rooms {
		log.Println(client.Username)
		if room.HasPlayer(client.Username) {
			room.rejoinPlayer(client)
			rm.joinMutex.Unlock()
			return
		}
	}
	if len(rm.rooms) == 0 || rm.rooms[len(rm.rooms)-1].isFull() {
		log.Println("Creating new room")
		room := NewRoom(2, gameName)
		rm.rooms = append(rm.rooms, room)
		go room.Run()
	}
	var current = rm.rooms[len(rm.rooms)-1]
	current.joinPlayer(client)
	rm.joinMutex.Unlock()
}
