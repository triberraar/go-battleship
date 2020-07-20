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
			log.Println("sending to room manager")
			c.RoomManager.Joins <- NewJoin(c)
		} else if c.Room != nil {
			c.Room.playerMessages <- roomMessage{c.PlayerID, message}
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
	maxPlayers         int
	clients            map[string]*gl.Battleship
	currentPlayer      string
	currentPlayerIndex int
	players            []string
	playerMessages     chan roomMessage
}

func NewRoom(maxPlayers int, client *Client) *Room {
	player := client.PlayerID
	pl := []string{}
	pl = append(pl, player)
	wc := gl.WriteClient{client.Conn, make(chan interface{})}
	r := Room{maxPlayers: maxPlayers, clients: make(map[string]*gl.Battleship), currentPlayer: player, playerMessages: make(chan roomMessage), players: pl, currentPlayerIndex: 0}
	r.clients[player] = gl.NewBattleship(&wc)
	client.Room = &r
	r.currentPlayerIndex = 20
	return &r
}

func (r *Room) joinPlayer(client *Client) {
	wc := gl.WriteClient{client.Conn, make(chan interface{})}
	player := client.PlayerID
	r.players = append(r.players, player)
	r.clients[player] = gl.NewBattleshipFromExisting(r.clients[r.currentPlayer], &wc)
	client.Room = r
	log.Printf("len %d", len(r.players))
	log.Printf("stuff %v", r.players)
	log.Printf("clients %d", len(r.clients))
	r.currentPlayerIndex = 50
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
	// log.Printf("maxPlayers %d, current players %d", r.maxPlayers, len(r.clients))
	return len(r.clients) == r.maxPlayers
}

type RoomManager struct {
	rooms []*Room
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
	if len(rm.rooms) == 0 {
		log.Println("making first room")
		room := NewRoom(2, client)
		rm.rooms = append(rm.rooms, room)
		go room.Run()
	} else {
		if rm.rooms[len(rm.rooms)-1].isFull() {
			log.Println("room is full, making new room")
			log.Printf("nother one %v", rm.rooms[len(rm.rooms)-1].players)
			room := NewRoom(2, client)
			rm.rooms = append(rm.rooms, room)
			go room.Run()
		} else {
			log.Println("joining room")
			rm.rooms[len(rm.rooms)-1].joinPlayer(client)
			log.Printf("nother one %v", rm.rooms[len(rm.rooms)-1].players)
		}
	}
	if rm.rooms[len(rm.rooms)-1].isFull() {
		for _, pl := range rm.rooms[len(rm.rooms)-1].players {
			rm.rooms[len(rm.rooms)-1].clients[pl].SendMessage(messages.NewGameStartedMessage(pl == rm.rooms[len(rm.rooms)-1].currentPlayer))
		}
	} else {
		for _, pl := range rm.rooms[len(rm.rooms)-1].players {
			rm.rooms[len(rm.rooms)-1].clients[pl].SendMessage(messages.NewAwaitingPlayersMessage())
		}
		// send waiting for players to all players
	}
}

func (rm *RoomManager) Run() {
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case jm := <-rm.Joins:
			log.Println("room manager got a join message")
			rm.joinRoom(jm.client)
		case <-ticker.C:
			log.Println(rm)
		}
	}
}

func (r *Room) Run() {
	for {
		select {
		case rm := <-r.playerMessages:
			if !r.isFull() {
			} else if rm.playerID != r.currentPlayer {
				log.Println("Other player sends message, skip")
			} else {
				r.currentPlayerIndex = (r.currentPlayerIndex + 1) % len(r.players)
				r.currentPlayer = r.players[r.currentPlayerIndex]
				r.clients[rm.playerID].Messages <- rm.message
				for _, pl := range r.players {
					r.clients[pl].SendMessage(messages.NewTurnMessage(pl == r.currentPlayer))
				}
			}
		}
	}
}
