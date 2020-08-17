package turndecider

import (
	"time"

	"github.com/triberraar/go-battleship/internal/client"
	"github.com/triberraar/go-battleship/internal/messages"
)

type TurnDecider struct {
	maxPlayers         int
	currentPlayerIndex int
	playersInOrder     []string
	waitTimer          *SecondsTimer
	clients            map[string]*client.Client
}

type SecondsTimer struct {
	timer *time.Timer
	end   time.Time
}

func (s *SecondsTimer) TimeRemaining() time.Duration {
	return s.end.Sub(time.Now())
}

func NewTurnDecider(maxPlayers int) *TurnDecider {
	return &TurnDecider{
		maxPlayers:         maxPlayers,
		currentPlayerIndex: 0,
		playersInOrder:     []string{},
		waitTimer:          nil,
		clients:            make(map[string]*client.Client),
	}
}

func (td *TurnDecider) AddPlayer(client *client.Client) {
	td.playersInOrder = append(td.playersInOrder, client.Username)
	td.clients[client.Username] = client
}

func (td *TurnDecider) Rejoin(client *client.Client) {
	td.clients[client.Username] = client
}

func (td *TurnDecider) Players() []string {
	return td.playersInOrder
}

func (td *TurnDecider) TimeRemaining() int {
	return int(td.waitTimer.TimeRemaining().Seconds())
}

func (td *TurnDecider) CurrentPlayer() string {
	if len(td.playersInOrder) == 0 {
		return ""
	}
	return td.playersInOrder[td.currentPlayerIndex]
}

func (td *TurnDecider) IsCurrentPlayer(username string) bool {
	return td.CurrentPlayer() == username
}

func (td *TurnDecider) Start(duration int) {
	d := time.Duration(duration) * time.Second
	timer := time.AfterFunc(d, func() {
		td.NextTurn(duration)
	})
	td.waitTimer = &SecondsTimer{timer, time.Now().Add(d)}
}

func (td *TurnDecider) ExtendTurn(duration int) {
	td.waitTimer.timer.Stop()
	d := time.Duration(duration) * time.Second
	timer := time.AfterFunc(d, func() {
		td.NextTurn(duration)
	})
	td.waitTimer = &SecondsTimer{timer, time.Now().Add(d)}
}

func (td *TurnDecider) NextTurn(duration int) {
	td.waitTimer.timer.Stop()
	td.currentPlayerIndex = (td.currentPlayerIndex + 1) % len(td.playersInOrder)
	for _, c := range td.clients {
		c.OutMessages <- messages.NewTurnMessage(c.Username, td.IsCurrentPlayer(c.Username), duration)
	}
	td.wait(duration)
}

func (td *TurnDecider) wait(duration int) {
	d := time.Duration(duration) * time.Second
	timer := time.AfterFunc(d, func() {
		td.NextTurn(duration)
	})
	td.waitTimer = &SecondsTimer{timer, time.Now().Add(d)}
}

func (td *TurnDecider) IsFull() bool {
	return len(td.playersInOrder) == td.maxPlayers
}
