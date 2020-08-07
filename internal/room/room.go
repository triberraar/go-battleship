package room

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/triberraar/go-battleship/internal/client"
	cl "github.com/triberraar/go-battleship/internal/client"
	"github.com/triberraar/go-battleship/internal/game"
	"github.com/triberraar/go-battleship/internal/game/creator"
	"github.com/triberraar/go-battleship/internal/messages"
)

type Player struct {
	game   game.Game
	client *client.Client
}

type Room struct {
	maxPlayers              int
	players                 map[string]*Player
	currentPlayerIndex      int
	playersInOrder          []string
	aggregateGameMessages   chan messages.GameMessage
	aggregateClientMessages chan client.ClientMessage
	waitTimer               *SecondsTimer
	gameDefinition          game.GameDefinition
	finished                bool
	removeMe                chan bool
	id                      uuid.UUID
}

type SecondsTimer struct {
	timer *time.Timer
	end   time.Time
}

func (s *SecondsTimer) TimeRemaining() time.Duration {
	return s.end.Sub(time.Now())
}

func NewRoom(maxPlayers int, gameName string) *Room {
	gd, _ := creator.NewGameDefinition(gameName)
	return &Room{id: uuid.New(), maxPlayers: maxPlayers, players: make(map[string]*Player), playersInOrder: []string{}, currentPlayerIndex: 0, aggregateGameMessages: make(chan messages.GameMessage, 10), aggregateClientMessages: make(chan cl.ClientMessage, 10), gameDefinition: gd, finished: false, removeMe: make(chan bool)}
}

func (r *Room) joinPlayer(client *cl.Client) {
	log.Printf("Joining %s", client.Username)
	r.playersInOrder = append(r.playersInOrder, client.Username)
	if len(r.players) == 0 {
		game, _ := creator.NewGame(r.gameDefinition.GameName(), client.Username)
		r.players[client.Username] = &Player{game: game, client: client}
	} else {
		game, _ := creator.NewGameFromExistion(r.gameDefinition.GameName(), r.players[r.currentPlayer()].game, client.Username)
		r.players[client.Username] = &Player{game: game, client: client}
	}
	r.aggregateMessages(client.Username)
	if r.isFull() {
		for _, pl := range r.players {
			pl.client.OutMessages <- messages.NewGameStartedMessage(pl.client.Username, pl.client.Username == r.currentPlayer(), r.gameDefinition.TurnDuration())
		}
		r.waitForAction(r.gameDefinition.TurnDuration())
	} else {
		client.OutMessages <- messages.NewAwaitingPlayersMessage(client.Username)
	}
}

func (r *Room) rejoinPlayer(client *cl.Client) {
	log.Printf("Rejoining %s", client.Username)
	r.players[client.Username].client.Close()

	r.players[client.Username].client = client
	r.players[client.Username].game.Rejoin()
	r.aggregateMessages(client.Username)
	client.OutMessages <- messages.NewTurnMessage(client.Username, r.currentPlayer() == client.Username, int(r.waitTimer.TimeRemaining().Seconds()))
	// should also send all the state of all other players :D
}

func (r *Room) HasPlayer(username string) bool {
	_, ok := r.players[username]
	return ok && !r.finished
}

func (r *Room) aggregateMessages(username string) {
	go func(c chan messages.GameMessage) {
		for msg := range c {
			r.aggregateGameMessages <- msg
		}
		log.Println("stopped aggregating")
	}(r.players[username].game.OutMessages())
	go func(c chan cl.ClientMessage) {
		for msg := range c {
			r.aggregateClientMessages <- msg
		}
		log.Println("stopped aggregating")
	}(r.players[username].client.InMessages)
}

func (r *Room) currentPlayer() string {
	return r.playersInOrder[r.currentPlayerIndex]
}

func (r Room) String() string {
	return fmt.Sprintf("Hej i am a room and can hold %d and have %d and it is this players turn: %s", r.maxPlayers, len(r.playersInOrder), r.currentPlayer())
}

func (r *Room) isFull() bool {
	return len(r.playersInOrder) == r.maxPlayers
}

func (r *Room) Run() {
	endChannel := make(chan bool, 5)
	r.listenForMessages(endChannel)
	log.Println("room ended")
	for _, pl := range r.players {
		close(pl.game.InMessages())
		close(pl.game.OutMessages())
		close(pl.client.InMessages)
	}
	r.removeMe <- true
}

func (r *Room) listenForMessages(endChannel chan bool) {
	for {
		select {
		case rm := <-r.aggregateClientMessages:
			bm := messages.BaseMessage{}
			json.Unmarshal(rm.Message, &bm)
			if !r.isFull() {
				log.Println("room not full, skipping")
			} else if rm.Username != r.currentPlayer() {
				log.Println("Other player sends message, skip")
			} else {
				r.waitTimer.timer.Stop()
				r.players[rm.Username].game.InMessages() <- rm.Message
				r.waitForAction(r.gameDefinition.TurnDuration())
			}
		case m := <-r.aggregateGameMessages:
			switch cm := m.Message.(type) {
			case messages.TurnMessage:
				r.nextConnection(cm.Duration)
			case messages.VictoryMessage:
				r.finished = true
				for _, pl := range r.players {
					if pl.client.Username == r.currentPlayer() {
						pl.client.OutMessages <- m.Message
					} else {
						pl.client.OutMessages <- messages.NewLossMessage(pl.client.Username)
					}
				}
				endChannel <- true
			default:
				r.players[m.Username].client.OutMessages <- m.Message
			}
		case <-endChannel:
			return
		}
	}
}

func (r *Room) waitForAction(duration int) {
	d := time.Duration(duration) * time.Second
	timer := time.AfterFunc(d, func() {
		r.nextConnection(duration)
	})
	r.waitTimer = &SecondsTimer{timer, time.Now().Add(d)}
}

func (r *Room) nextConnection(duration int) {
	r.waitTimer.timer.Stop()
	r.currentPlayerIndex = (r.currentPlayerIndex + 1) % len(r.players)
	for _, pl := range r.players {
		pl.client.OutMessages <- messages.NewTurnMessage(pl.client.Username, pl.client.Username == r.currentPlayer(), duration)
	}
	r.waitForAction(duration)
}

func (r *Room) isFinished() bool {
	return r.finished
}
