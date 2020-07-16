package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	gl "github.com/triberraar/go-battleship/internal/game_logic"
	"github.com/triberraar/go-battleship/internal/messages"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}} // use default options

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Conn        *websocket.Conn
	RoomManager *RoomManager
	Room        *Room
	PlayerID    string
}

type roomMessage struct {
	playerID string
	message  []byte
}

func (c *Client) ReadPump() {
	defer c.Conn.Close()
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Print("error: %v", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("unexcpected close")
			}
			break
		}
		bm := messages.BaseMessage{}
		json.Unmarshal(message, &bm)
		if bm.Type == "PLAY" {
			c.RoomManager.Joins <- NewJoin(c)
		}
		if c.Room != nil {
			c.Room.messages <- roomMessage{c.PlayerID, message}
		} else {
			log.Println("no room to send message")
		}
	}
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	client := &Client{Conn: c}
	// battleship := NewBattleship(client)

	// go client.WritePump()
	go client.ReadPump()
	// go battleship.Run()
}

type Room struct {
	maxPlayers    int
	clients       map[string]*gl.Battleship
	currentPlayer string
	messages      chan roomMessage
}

type Join struct {
	client *Client
}

func NewJoin(client *Client) Join {
	return Join{
		client,
	}
}

func (r Room) String() string {
	return fmt.Sprintf("Hej i am a room and can hold %d and have %d and it is this players turn: %s", r.maxPlayers, len(r.clients), r.currentPlayer)
}

func (r *Room) isFull() bool {
	log.Printf("maxPlayers %d, current players %d", r.maxPlayers, len(r.clients))
	return len(r.clients) == r.maxPlayers
}

type RoomManager struct {
	rooms []Room
	Joins chan Join
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		Joins: make(chan Join),
	}
}

func (rm RoomManager) String() string {
	return fmt.Sprintf("I be the room manager, here be my rooms %v", rm.rooms)
}

func (rm *RoomManager) joinRoom(client *Client) {
	player := client.PlayerID
	if len(rm.rooms) == 0 {
		log.Println("making first room")
		room := Room{maxPlayers: 2, clients: make(map[string]*gl.Battleship), currentPlayer: player, messages: make(chan roomMessage)}
		wc := gl.WriteClient{client.Conn, make(chan interface{})}
		go wc.WritePump()
		room.clients[player] = gl.NewBattleship(&wc)
		rm.rooms = append(rm.rooms, room)
		client.Room = &room
		go room.Run()
		go room.clients[player].Run()
	} else {
		if rm.rooms[len(rm.rooms)-1].isFull() {
			log.Println("room is full, making new room")
			room := Room{maxPlayers: 2, clients: make(map[string]*gl.Battleship), currentPlayer: player, messages: make(chan roomMessage)}
			wc := gl.WriteClient{client.Conn, make(chan interface{})}
			go wc.WritePump()
			room.clients[player] = gl.NewBattleship(&wc)
			rm.rooms = append(rm.rooms, room)
			client.Room = &room
			go room.Run()
			go room.clients[player].Run()
		} else {
			log.Println("joining room")
			wc := gl.WriteClient{client.Conn, make(chan interface{})}
			go wc.WritePump()
			rm.rooms[len(rm.rooms)-1].clients[player] = gl.NewBattleshipFromExisting(rm.rooms[len(rm.rooms)-1].clients[rm.rooms[len(rm.rooms)-1].currentPlayer], &wc)
			client.Room = &rm.rooms[len(rm.rooms)-1]
			go rm.rooms[len(rm.rooms)-1].clients[player].Run()
		}
	}
}

func (rm *RoomManager) Run() {
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case jm := <-rm.Joins:
			rm.joinRoom(jm.client)

		case <-ticker.C:
			log.Println(rm)
		}
	}
}

func (r *Room) Run() {
	for {
		select {
		case rm := <-r.messages:
			r.clients[rm.playerID].Messages <- rm.message
		}
	}
}
