package battleship

import (
	"log"

	"github.com/google/uuid"
	"github.com/triberraar/go-battleship/internal/client"
	"github.com/triberraar/go-battleship/internal/messages"
	"github.com/triberraar/go-battleship/internal/turndecider"
)

type BattleshipMatch struct {
	id uuid.UUID

	battleships map[string]*Battleship
	clients     map[string]*client.Client
	turnDecider *turndecider.TurnDecider

	RemoveMe chan bool
}

func NewBattleshipMatch(maxPlayers int) *BattleshipMatch {
	return &BattleshipMatch{uuid.New(), make(map[string]*Battleship), make(map[string]*client.Client), turndecider.NewTurnDecider(maxPlayers), make(chan bool)}

}

func (bm BattleshipMatch) GetRemoveChannel() chan bool {
	return bm.RemoveMe
}

func (bm BattleshipMatch) ShouldRejoin(username string) bool {
	return bm.battleships[username] != nil
}

func (bm *BattleshipMatch) Join(client *client.Client) {
	log.Println("joining")
	bm.turnDecider.AddPlayer(client)
	if len(bm.battleships) == 0 {
		game := NewBattleship(client.Username)
		bm.battleships[client.Username] = game
		bm.clients[client.Username] = client
	} else {
		game := bm.battleships[bm.turnDecider.CurrentPlayer()].NewBattleshipFromExisting(client.Username)
		bm.battleships[client.Username] = game
		bm.clients[client.Username] = client
	}

	if bm.turnDecider.IsFull() {
		for _, c := range bm.clients {
			c.OutMessages <- messages.NewGameStartedMessage(c.Username, bm.turnDecider.IsCurrentPlayer(c.Username), turnDuration, bm.turnDecider.Players())
		}
		bm.turnDecider.Start(turnDuration)
	} else {
		client.OutMessages <- messages.NewAwaitingPlayersMessage(client.Username)
	}

	go func(c chan []byte) {
		for m := range c {
			if bm.turnDecider.IsCurrentPlayer(client.Username) {
				bm.battleships[client.Username].Process(m)
			}
		}
		log.Printf("reading inmessages ended for %s", client.Username)

	}(bm.clients[client.Username].InMessages)
	go bm.processGameMessages(bm.battleships[client.Username].OutMessages, client.Username)
}

func (bm *BattleshipMatch) processGameMessages(c chan interface{}, username string) {
	for m := range c {
		switch cm := m.(type) {
		case messages.TurnMessage:
			bm.turnDecider.NextTurn(cm.Duration)
		case messages.VictoryMessage:
			for _, c := range bm.clients {
				if bm.turnDecider.IsCurrentPlayer(c.Username) {
					c.OutMessages <- m
				} else {
					c.OutMessages <- messages.NewLossMessage(c.Username)
				}
			}
			for _, bs := range bm.battleships {
				close(bs.OutMessages)
			}
			bm.RemoveMe <- true
		case messages.ShipDestroyedMessage:
			for _, c := range bm.clients {
				if bm.turnDecider.IsCurrentPlayer(c.Username) {
					c.OutMessages <- m
				} else {
					c.OutMessages <- messages.NewOpponentDestroyedShipMessage(cm.Username)
				}
			}
		case messages.TurnExtendedMessage:
			bm.turnDecider.ExtendTurn(turnDuration)
			bm.clients[cm.Username].OutMessages <- m
		default:
			for _, c := range bm.clients {
				c.OutMessages <- m
			}
		}
	}
	log.Printf("stopped processing game messages for %s", username)
}

func (bm BattleshipMatch) Rejoin(client *client.Client) {
	bm.clients[client.Username].Conn.Close()
	bm.clients[client.Username] = client
	bm.battleships[client.Username].Rejoin()
	bm.turnDecider.Rejoin(client)
	go func(c chan []byte) {
		for m := range c {
			if bm.turnDecider.IsCurrentPlayer(client.Username) {
				bm.battleships[client.Username].Process(m)
			}
		}
		log.Printf("reading inmessages ended for %s", client.Username)

	}(bm.clients[client.Username].InMessages)

	client.OutMessages <- messages.NewTurnMessage(client.Username, bm.turnDecider.IsCurrentPlayer(client.Username), bm.turnDecider.TimeRemaining())
}

func (bm BattleshipMatch) GetID() uuid.UUID {
	return bm.id
}

/*
	ShouldRejoin(username string) bool
	Join(client client.Client)
	Rejoin(client client.Client)
*/
