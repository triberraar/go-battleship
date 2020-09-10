package turndecider

import (
	"log"
	"time"

	"github.com/triberraar/battleship/internal/client"
	"github.com/triberraar/battleship/internal/messages"
)

type TurnDecider struct {
	maxPlayers         int
	currentPlayerIndex int
	playersInOrder     []string
	waitTimer          *SecondsTimer
	clients            map[string]*client.Client
	endTimer           *SecondsTimer
	duration           int
	stop               chan string
}

type SecondsTimer struct {
	timer *time.Timer
	end   time.Time
}

func (s *SecondsTimer) TimeRemaining() time.Duration {
	return s.end.Sub(time.Now())
}

func NewTurnDecider(maxPlayers int, duration int, stop chan string) *TurnDecider {
	td := TurnDecider{
		maxPlayers:         maxPlayers,
		currentPlayerIndex: 0,
		playersInOrder:     []string{},
		waitTimer:          nil,
		clients:            make(map[string]*client.Client),
		duration:           duration,
		stop:               stop,
	}
	return &td
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

func (td *TurnDecider) Start() {
	d := time.Duration(td.duration) * time.Second
	timer := time.AfterFunc(d, func() {
		td.NextTurn(false)
	})
	td.waitTimer = &SecondsTimer{timer, time.Now().Add(d)}

	td.resetEndTimer()
}

func (td *TurnDecider) ExtendTurn() {
	td.waitTimer.timer.Stop()
	d := time.Duration(td.duration) * time.Second
	timer := time.AfterFunc(d, func() {
		td.NextTurn(false)
	})
	td.waitTimer = &SecondsTimer{timer, time.Now().Add(d)}

	log.Println("extending end time")
	td.resetEndTimer()
}

func (td *TurnDecider) NextTurn(action bool) {
	td.waitTimer.timer.Stop()
	td.currentPlayerIndex = (td.currentPlayerIndex + 1) % len(td.playersInOrder)
	for _, c := range td.clients {
		c.OutMessages <- messages.NewTurnMessage(c.Username, td.IsCurrentPlayer(c.Username), td.duration)
	}
	td.wait(action)
}

func (td *TurnDecider) wait(action bool) {
	d := time.Duration(td.duration) * time.Second
	timer := time.AfterFunc(d, func() {
		td.NextTurn(false)
	})
	td.waitTimer = &SecondsTimer{timer, time.Now().Add(d)}

	if action {
		log.Println("extending end time in wait")
		td.resetEndTimer()
	}
}

func (td *TurnDecider) resetEndTimer() {
	if td.endTimer != nil {
		td.endTimer.timer.Stop()
	}
	d := time.Duration(td.duration*4) * time.Second
	endTimer := time.AfterFunc(d, func() {
		log.Println("ended due to timeout")
		if td.waitTimer != nil {
			td.waitTimer.timer.Stop()
		}
		for _, c := range td.clients {
			c.OutMessages <- messages.NewCancelledMessage()
		}
		log.Println("sending timeout sig")
		td.stop <- "timeout"
		log.Println("sent timeout sig")
	})
	td.endTimer = &SecondsTimer{endTimer, time.Now().Add(d)}
}

func (td *TurnDecider) IsFull() bool {
	return len(td.playersInOrder) == td.maxPlayers
}
