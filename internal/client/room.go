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

type RoomManager struct {
	rooms []*Room
	Joins chan Join
}

var RM *RoomManager

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
	Room        *Room
	PlayerID    string
	OutMessages chan interface{}
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
			log.Printf("error: %v", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("unexcpected close")
			}
			break
		}
		bm := messages.BaseMessage{}
		json.Unmarshal(message, &bm)
		if bm.Type == "PLAY" {
			RM.Joins <- NewJoin(c)
		} else if c.Room != nil {
			c.Room.playerMessages <- roomMessage{c.PlayerID, message}
		} else {
			log.Println("no room to send message")
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	for {
		select {
		case message := <-c.OutMessages:
			c.Conn.WriteJSON(message)
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type Room struct {
	maxPlayers            int
	players               map[string]*Player
	currentPlayer         string
	currentPlayerIndex    int
	playersInOrder        []string
	playerMessages        chan roomMessage // change this
	aggregateGameMessages chan gl.OutMessage
}

type Player struct {
	playerID string
	game     *gl.Battleship
	client   *Client
}

func NewRoom(maxPlayers int, client *Client) *Room {
	player := client.PlayerID
	pl := []string{}
	pl = append(pl, player)
	r := Room{maxPlayers: maxPlayers, players: make(map[string]*Player), currentPlayer: player, playerMessages: make(chan roomMessage, 10), playersInOrder: pl, currentPlayerIndex: 0, aggregateGameMessages: make(chan gl.OutMessage, 10)}
	r.players[player] = &Player{playerID: player, game: gl.NewBattleship(player), client: client}
	client.Room = &r
	go func(c chan gl.OutMessage) {
		for msg := range c {
			r.aggregateGameMessages <- msg
		}
	}(r.players[player].game.OutMessages)
	return &r
}

func (r *Room) joinPlayer(client *Client) {
	player := client.PlayerID
	r.playersInOrder = append(r.playersInOrder, player)
	r.players[player] = &Player{playerID: player, game: gl.NewBattleshipFromExisting(r.players[r.currentPlayer].game, player), client: client}
	go func(c chan gl.OutMessage) {
		for msg := range c {
			r.aggregateGameMessages <- msg
		}
	}(r.players[player].game.OutMessages)
	client.Room = r
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
	return fmt.Sprintf("Hej i am a room and can hold %d and have %d and it is this players turn: %s", r.maxPlayers, len(r.playersInOrder), r.currentPlayer)
}

func (r *Room) isFull() bool {
	return len(r.playersInOrder) == r.maxPlayers
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		Joins: make(chan Join, 50),
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
		case rm := <-r.playerMessages:
			if !r.isFull() {
				log.Println("room not full, skipping")
			} else if rm.playerID != r.currentPlayer {
				log.Println("Other player sends message, skip")
			} else {
				// 	r.currentPlayerIndex = (r.currentPlayerIndex + 1) % len(r.players)
				// 	r.currentPlayer = r.playersInOrder[r.currentPlayerIndex]
				r.players[rm.playerID].game.InMessages <- rm.message
				// 	for _, pl := range r.players {
				// 		r.players[pl.playerID].client.OutMessages <- messages.NewTurnMessage(pl.playerID == r.currentPlayer)
				// }
			}
		case m := <-r.aggregateGameMessages:
			switch m.Message.(type) {
			case messages.TurnMessage:
				r.currentPlayerIndex = (r.currentPlayerIndex + 1) % len(r.players)
				r.currentPlayer = r.playersInOrder[r.currentPlayerIndex]
				for _, pl := range r.players {
					r.players[pl.playerID].client.OutMessages <- messages.NewTurnMessage(pl.playerID == r.currentPlayer)
				}
			case messages.VictoryMessage:
				for _, pl := range r.players {
					if pl.playerID == r.currentPlayer {
						r.players[pl.playerID].client.OutMessages <- messages.NewVictoryMessage()
					} else {
						r.players[pl.playerID].client.OutMessages <- messages.NewLossMessage()
					}

				}
			default:
				r.players[m.PlayerID].client.OutMessages <- m.Message
			}

		}

	}
}
