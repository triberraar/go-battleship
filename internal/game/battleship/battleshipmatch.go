package battleship

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/triberraar/go-battleship/internal/client"
	"github.com/triberraar/go-battleship/internal/messages"
)

type BattleshipMatch struct {
	id uuid.UUID

	battleships map[string]*Battleship
	clients     map[string]*client.Client

	maxPlayers         int
	currentPlayerIndex int
	playersInOrder     []string
	waitTimer          *SecondsTimer

	removeMe chan bool
}

type SecondsTimer struct {
	timer *time.Timer
	end   time.Time
}

func (s *SecondsTimer) TimeRemaining() time.Duration {
	return s.end.Sub(time.Now())
}

func NewBattleshipMatch(maxPlayers int) *BattleshipMatch {
	return &BattleshipMatch{uuid.New(), make(map[string]*Battleship), make(map[string]*client.Client), maxPlayers, 0, []string{}, nil, make(chan bool)}

}

func (bm BattleshipMatch) ShouldRejoin(username string) bool {
	return bm.battleships[username] != nil
}

func (bm *BattleshipMatch) Join(client *client.Client) {
	log.Println("joining")
	bm.playersInOrder = append(bm.playersInOrder, client.Username)
	if len(bm.battleships) == 0 {
		game := NewBattleship(client.Username)
		bm.battleships[client.Username] = game
		bm.clients[client.Username] = client
	} else {
		game := bm.battleships[bm.currentPlayer()].NewBattleshipFromExisting(client.Username)
		bm.battleships[client.Username] = game
		bm.clients[client.Username] = client
	}

	if bm.allReady() {
		for _, c := range bm.clients {
			c.OutMessages <- messages.NewGameStartedMessage(c.Username, c.Username == bm.currentPlayer(), turnDuration, bm.playersInOrder)
		}
		bm.waitForAction(turnDuration)
	} else {
		client.OutMessages <- messages.NewAwaitingPlayersMessage(client.Username)
	}

	go func(c chan []byte) {
		for m := range c {
			if bm.currentPlayer() == client.Username {
				bm.battleships[client.Username].Process(m)
			}
		}
		log.Printf("reading inmessages ended for %s", client.Username)

	}(bm.clients[client.Username].InMessages2)
	go bm.processGameMessages(bm.battleships[client.Username].OutMessages2, client.Username)
}

func (bm *BattleshipMatch) processGameMessages(c chan interface{}, username string) {
	for m := range c {
		switch cm := m.(type) {
		case messages.TurnMessage:
			bm.nextConnection(cm.Duration)
		case messages.VictoryMessage:
			for _, c := range bm.clients {
				if c.Username == bm.currentPlayer() {
					log.Printf("forwarding victory")
					c.OutMessages <- m
					log.Printf("done forwarding victory")
				} else {
					log.Printf("sending loss")
					c.OutMessages <- messages.NewLossMessage(c.Username)
					log.Printf("done sending loss")
				}
			}
			for _, bs := range bm.battleships {
				close(bs.OutMessages2)
			}
		case messages.ShipDestroyedMessage:
			for _, c := range bm.clients {
				if c.Username == bm.currentPlayer() {
					c.OutMessages <- m
				} else {
					c.OutMessages <- messages.NewOpponentDestroyedShipMessage(cm.Username)
				}
			}
		default:
			for _, c := range bm.clients {
				c.OutMessages <- m
			}
		}
	}
	log.Printf("stopped processing game messages for %s")
	bm.removeMe <- true
}

func (bm BattleshipMatch) Rejoin(client *client.Client) {
	bm.clients[client.Username].Conn.Close()
	bm.clients[client.Username] = client
	bm.battleships[client.Username].Rejoin()

	client.OutMessages <- messages.NewTurnMessage(client.Username, bm.currentPlayer() == client.Username, int(bm.waitTimer.TimeRemaining().Seconds()))
}

func (bm *BattleshipMatch) currentPlayer() string {
	if len(bm.playersInOrder) == 0 {
		return ""
	}
	return bm.playersInOrder[bm.currentPlayerIndex]
}

func (bm *BattleshipMatch) allReady() bool {
	return len(bm.battleships) == bm.maxPlayers
}

func (bm *BattleshipMatch) waitForAction(duration int) {
	d := time.Duration(duration) * time.Second
	timer := time.AfterFunc(d, func() {
		bm.nextConnection(duration)
	})
	bm.waitTimer = &SecondsTimer{timer, time.Now().Add(d)}
}

func (bm *BattleshipMatch) nextConnection(duration int) {
	bm.waitTimer.timer.Stop()
	bm.currentPlayerIndex = (bm.currentPlayerIndex + 1) % len(bm.battleships)
	for _, c := range bm.clients {
		c.OutMessages <- messages.NewTurnMessage(c.Username, c.Username == bm.currentPlayer(), duration)
	}
	bm.waitForAction(duration)
}

func (bm BattleshipMatch) GetID() uuid.UUID {
	return bm.id
}

/*
	ShouldRejoin(username string) bool
	Join(client client.Client)
	Rejoin(client client.Client)
*/
