package room

import (
	"fmt"
	"log"

	"github.com/triberraar/go-battleship/internal/client"
	cl "github.com/triberraar/go-battleship/internal/client"
	"github.com/triberraar/go-battleship/internal/game/battleship"
	"github.com/triberraar/go-battleship/internal/messages"
)

type Player struct {
	playerID string
	game     *battleship.Battleship
	client   *client.Client
}

type Room struct {
	maxPlayers              int
	players                 map[string]*Player
	currentPlayerIndex      int
	playersInOrder          []string
	aggregateGameMessages   chan messages.GameMessage
	aggregateClientMessages chan client.ClientMessage
}

func NewRoom(maxPlayers int) *Room {
	return &Room{maxPlayers: maxPlayers, players: make(map[string]*Player), playersInOrder: []string{}, currentPlayerIndex: 0, aggregateGameMessages: make(chan messages.GameMessage, 10), aggregateClientMessages: make(chan cl.ClientMessage, 10)}
}

func (r *Room) joinPlayer(client *cl.Client) {
	playerID := client.PlayerID
	r.playersInOrder = append(r.playersInOrder, playerID)
	if len(r.players) == 0 {
		r.players[playerID] = &Player{playerID: playerID, game: battleship.NewBattleship(playerID), client: client}
	} else {
		r.players[playerID] = &Player{playerID: playerID, game: r.players[r.currentPlayerID()].game.NewBattleshipFromExisting(playerID), client: client}
	}
	r.aggregateMessages(playerID)
	if r.isFull() {
		for _, pl := range r.players {
			pl.client.OutMessages <- messages.NewGameStartedMessage(pl.playerID == r.currentPlayerID())
		}
	} else {
		r.players[playerID].client.OutMessages <- messages.NewAwaitingPlayersMessage()
	}
}

func (r *Room) aggregateMessages(playerID string) {
	go func(c chan messages.GameMessage) {
		for msg := range c {
			r.aggregateGameMessages <- msg
		}
	}(r.players[playerID].game.OutMessages)
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
				r.players[rm.PlayerID].game.InMessages <- rm.Message
			}
		case m := <-r.aggregateGameMessages:
			switch m.Message.(type) {
			case messages.TurnMessage:
				r.currentPlayerIndex = (r.currentPlayerIndex + 1) % len(r.players)
				for _, pl := range r.players {
					r.players[pl.playerID].client.OutMessages <- messages.NewTurnMessage(pl.playerID == r.currentPlayerID())
				}
			case messages.VictoryMessage:
				for _, pl := range r.players {
					if pl.playerID == r.currentPlayerID() {
						r.players[pl.playerID].client.OutMessages <- m
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
