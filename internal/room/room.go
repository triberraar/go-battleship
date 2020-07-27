package room

import (
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
	playerID string
	game     game.Game
	client   *client.Client
}

type Room struct {
	maxPlayers              int
	players                 map[string]*Player
	currentPlayerIndex      int
	playersInOrder          []string
	aggregateGameMessages   chan messages.GameMessage
	aggregateClientMessages chan client.ClientMessage
	waitTimer               *time.Timer
	gameDefinition          game.GameDefinition
}

func NewRoom(maxPlayers int, gameName string) *Room {
	gd, _ := creator.NewGameDefinition(gameName)
	return &Room{maxPlayers: maxPlayers, players: make(map[string]*Player), playersInOrder: []string{}, currentPlayerIndex: 0, aggregateGameMessages: make(chan messages.GameMessage, 10), aggregateClientMessages: make(chan cl.ClientMessage, 10), gameDefinition: gd}
}

func (r *Room) joinPlayer(client *cl.Client) {
	playerID := client.PlayerID
	r.playersInOrder = append(r.playersInOrder, playerID)
	if len(r.players) == 0 {
		game, _ := creator.NewGame(r.gameDefinition.GameName(), playerID)
		r.players[playerID] = &Player{playerID: playerID, game: game, client: client}
	} else {
		game, _ := creator.NewGameFromExistion(r.gameDefinition.GameName(), r.players[r.currentPlayerID()].game, playerID)
		r.players[playerID] = &Player{playerID: playerID, game: game, client: client}
	}
	r.aggregateMessages(playerID)
	if r.isFull() {
		for _, pl := range r.players {
			pl.client.OutMessages <- messages.NewGameStartedMessage(pl.playerID == r.currentPlayerID(), r.gameDefinition.TurnDuration())
		}
		r.waitForAction(r.gameDefinition.TurnDuration())
	} else {
		r.players[playerID].client.OutMessages <- messages.NewAwaitingPlayersMessage()
	}
}

func (r *Room) aggregateMessages(playerID string) {
	go func(c chan messages.GameMessage) {
		for msg := range c {
			r.aggregateGameMessages <- msg
		}
	}(r.players[playerID].game.OutMessages())
	go func(c chan cl.ClientMessage) {
		for msg := range c {
			r.aggregateClientMessages <- msg
		}
	}(r.players[playerID].client.InMessages)
}

func (r *Room) currentPlayerID() string {
	return r.playersInOrder[r.currentPlayerIndex]
}

func (r Room) String() string {
	return fmt.Sprintf("Hej i am a room and can hold %d and have %d and it is this players turn: %s", r.maxPlayers, len(r.playersInOrder), r.currentPlayerID())
}

func (r *Room) isFull() bool {
	return len(r.playersInOrder) == r.maxPlayers
}

func (r *Room) Run() {
	for {
		select {
		case rm := <-r.aggregateClientMessages:
			if !r.isFull() {
				log.Println("room not full, skipping")
			} else if rm.PlayerID != r.currentPlayerID() {
				log.Println("Other player sends message, skip")
			} else {
				r.waitTimer.Stop()
				r.players[rm.PlayerID].game.InMessages() <- rm.Message
				r.waitForAction(r.gameDefinition.TurnDuration())
			}
		case m := <-r.aggregateGameMessages:
			switch cm := m.Message.(type) {
			case messages.TurnMessage:
				r.nextPlayer(cm.Duration)
			case messages.VictoryMessage:
				for _, pl := range r.players {
					if pl.playerID == r.currentPlayerID() {
						r.players[pl.playerID].client.OutMessages <- m.Message
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

func (r *Room) waitForAction(duration int) {
	r.waitTimer = time.AfterFunc(time.Duration(duration)*time.Second, func() {
		r.nextPlayer(duration)
	})
}

func (r *Room) nextPlayer(duration int) {
	r.waitTimer.Stop()
	r.currentPlayerIndex = (r.currentPlayerIndex + 1) % len(r.players)
	for _, pl := range r.players {
		r.players[pl.playerID].client.OutMessages <- messages.NewTurnMessage(pl.playerID == r.currentPlayerID(), duration)
	}
	r.waitForAction(duration)
}
