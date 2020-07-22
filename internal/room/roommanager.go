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

func (rm *RoomManager) JoinRoom(client *client.Client) {
	rm.joinMutex.Lock()
	if len(rm.rooms) == 0 || rm.rooms[len(rm.rooms)-1].isFull() {
		log.Println("Creating new room")
		room := NewRoom(2)
		rm.rooms = append(rm.rooms, room)
		go room.Run()
	}
	var current = rm.rooms[len(rm.rooms)-1]
	current.joinPlayer(client)
	rm.joinMutex.Unlock()
}