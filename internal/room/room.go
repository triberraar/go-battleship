package room

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/triberraar/go-battleship/internal/client"
	cl "github.com/triberraar/go-battleship/internal/client"
	"github.com/triberraar/go-battleship/internal/game"
	"github.com/triberraar/go-battleship/internal/game/creator"
	"github.com/triberraar/go-battleship/internal/messages"
)

type Player struct {
	connectionID string
	game         game.Game
	client       *client.Client
	username     string
}

type Room struct {
	maxPlayers              int
	players                 map[string]*Player
	currentConnectionIndex  int
	connectionsInOrder      []string
	aggregateGameMessages   chan messages.GameMessage
	aggregateClientMessages chan client.ClientMessage
	aggregateLeavers        chan string
	waitTimer               *time.Timer
	gameDefinition          game.GameDefinition
	leavers                 map[string]string
}

func NewRoom(maxPlayers int, gameName string) *Room {
	gd, _ := creator.NewGameDefinition(gameName)
	return &Room{maxPlayers: maxPlayers, players: make(map[string]*Player), connectionsInOrder: []string{}, currentConnectionIndex: 0, aggregateGameMessages: make(chan messages.GameMessage, 10), aggregateClientMessages: make(chan cl.ClientMessage, 10), gameDefinition: gd, aggregateLeavers: make(chan string, 2), leavers: make(map[string]string)}
}

func (r *Room) joinPlayer(client *cl.Client) {
	connectionID := client.ConnectionID
	r.connectionsInOrder = append(r.connectionsInOrder, connectionID)
	if len(r.players) == 0 {
		game, _ := creator.NewGame(r.gameDefinition.GameName(), connectionID)
		r.players[connectionID] = &Player{connectionID: connectionID, game: game, client: client}
	} else {
		game, _ := creator.NewGameFromExistion(r.gameDefinition.GameName(), r.players[r.currentConnectionID()].game, connectionID)
		r.players[connectionID] = &Player{connectionID: connectionID, game: game, client: client}
	}
	r.aggregateMessages(connectionID)
	if r.isFull() {
		for _, pl := range r.players {
			pl.client.OutMessages <- messages.NewGameStartedMessage(pl.connectionID == r.currentConnectionID(), r.gameDefinition.TurnDuration())
		}
		r.waitForAction(r.gameDefinition.TurnDuration())
	} else {
		r.players[connectionID].client.OutMessages <- messages.NewAwaitingPlayersMessage()
	}
}

func (r *Room) aggregateMessages(connectionID string) {
	go func(c chan messages.GameMessage) {
		for msg := range c {
			r.aggregateGameMessages <- msg
		}
	}(r.players[connectionID].game.OutMessages())
	go func(c chan cl.ClientMessage) {
		for msg := range c {
			r.aggregateClientMessages <- msg
		}
	}(r.players[connectionID].client.InMessages)
	go func(c chan string) {
		for msg := range c {
			r.aggregateLeavers <- msg
		}
	}(r.players[connectionID].client.Leavers)
}

func (r *Room) currentConnectionID() string {
	return r.connectionsInOrder[r.currentConnectionIndex]
}

func (r Room) String() string {
	return fmt.Sprintf("Hej i am a room and can hold %d and have %d and it is this players turn: %s", r.maxPlayers, len(r.connectionsInOrder), r.currentConnectionID())
}

func (r *Room) isFull() bool {
	return len(r.connectionsInOrder) == r.maxPlayers
}

func (r *Room) Run() {
	for {
		select {
		case rm := <-r.aggregateClientMessages:
			bm := messages.BaseMessage{}
			json.Unmarshal(rm.Message, &bm)
			if bm.Type == "PLAY" {
				pm := messages.PlayMessage{}
				json.Unmarshal(rm.Message, &pm)
				r.players[rm.ConnectionID].username = pm.Username
			} else if !r.isFull() {
				log.Println("room not full, skipping")
			} else if rm.ConnectionID != r.currentConnectionID() {
				log.Println("Other player sends message, skip")
			} else {
				r.waitTimer.Stop()
				r.players[rm.ConnectionID].game.InMessages() <- rm.Message
				r.waitForAction(r.gameDefinition.TurnDuration())
			}
		case m := <-r.aggregateGameMessages:
			switch cm := m.Message.(type) {
			case messages.TurnMessage:
				r.nextConnection(cm.Duration)
			case messages.VictoryMessage:
				for _, pl := range r.players {
					if pl.connectionID == r.currentConnectionID() {
						r.players[pl.connectionID].client.OutMessages <- m.Message
					} else {
						r.players[pl.connectionID].client.OutMessages <- messages.NewLossMessage()
					}

				}
			default:
				r.players[m.ConnectionID].client.OutMessages <- m.Message
			}
		case m := <-r.aggregateLeavers:
			log.Printf("player left %s", m)
			// add player id to leavers
			// remove from players
			// remove from players inorder
			// reset currentplayer if needed
		}

	}
}

func (r *Room) waitForAction(duration int) {
	r.waitTimer = time.AfterFunc(time.Duration(duration)*time.Second, func() {
		r.nextConnection(duration)
	})
}

func (r *Room) nextConnection(duration int) {
	r.waitTimer.Stop()
	r.currentConnectionIndex = (r.currentConnectionIndex + 1) % len(r.players)
	for _, pl := range r.players {
		r.players[pl.connectionID].client.OutMessages <- messages.NewTurnMessage(pl.connectionID == r.currentConnectionID(), duration)
	}
	r.waitForAction(duration)
}
