package game_logic

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/triberraar/go-battleship/internal/messages"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type WriteClient struct {
	Conn *websocket.Conn
	Send chan interface{}
}

func (c *WriteClient) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	for {
		select {
		case message := <-c.Send:
			c.Conn.WriteJSON(message)
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type OutMessage struct {
	PlayerID string
	Message  interface{}
}

type ship struct {
	x        int
	y        int
	size     int
	vertical bool
	hits     int
}

func newShip(size int) ship {
	return ship{
		size: size,
	}
}

func (s *ship) hit() {
	s.hits++
}

func (s *ship) isDestroyed() bool {
	return s.hits == s.size
}

type tile struct {
	status string
	ship   *ship
}

func (t *tile) hasShip() bool {
	return t.ship != nil
}

type Battleship struct {
	InMessages  chan []byte
	OutMessages chan OutMessage
	playerID    string

	board     [][]tile
	dimension int
	ships     [6]ship
	victory   bool
}

func (bs *Battleship) SendMessage(message interface{}) {
	bs.OutMessages <- OutMessage{bs.playerID, message}
}

func (bs *Battleship) Run() {

	for {
		message := <-bs.InMessages
		bm := messages.BaseMessage{}
		json.Unmarshal(message, &bm)
		if bm.Type == "FIRE" && !bs.victory {
			fm := messages.FireMessage{}
			json.Unmarshal(message, &fm)
			if bs.board[fm.Coordinate.X][fm.Coordinate.Y].status == "fired" {
				continue
			}
			if bs.board[fm.Coordinate.X][fm.Coordinate.Y].hasShip() {

				bs.board[fm.Coordinate.X][fm.Coordinate.Y].status = "fired"
				bs.board[fm.Coordinate.X][fm.Coordinate.Y].ship.hit()
				if bs.board[fm.Coordinate.X][fm.Coordinate.Y].ship.isDestroyed() {
					coordinate := messages.Coordinate{X: bs.board[fm.Coordinate.X][fm.Coordinate.Y].ship.x, Y: bs.board[fm.Coordinate.X][fm.Coordinate.Y].ship.y}
					// bs.client.Send <- messages.NewShipDestroyedMessage(coordinate, bs.board[fm.Coordinate.X][fm.Coordinate.Y].ship.size, bs.board[fm.Coordinate.X][fm.Coordinate.Y].ship.vertical)
					bs.SendMessage(messages.NewShipDestroyedMessage(coordinate, bs.board[fm.Coordinate.X][fm.Coordinate.Y].ship.size, bs.board[fm.Coordinate.X][fm.Coordinate.Y].ship.vertical))
				} else {
					bs.SendMessage(messages.NewHitMessage(fm.Coordinate))
				}
				if bs.hasVictory() {
					bs.victory = true
					bs.SendMessage(messages.NewVictoryMessage())
				}
			} else {
				bs.SendMessage(messages.NewMissMessage(fm.Coordinate))
				log.Printf("player %s missed", bs.playerID)
				bs.board[fm.Coordinate.X][fm.Coordinate.Y].status = "fired"
			}
		}
	}
}

func NewBattleship(playerID string) *Battleship {
	bs := Battleship{
		dimension:   10,
		victory:     false,
		InMessages:  make(chan []byte, 10),
		OutMessages: make(chan OutMessage, 10),
		playerID:    playerID,
	}
	bs.newBoard()
	var shipSizes [len(bs.ships)]int
	for i := 0; i < len(bs.ships); i++ {
		shipSizes[i] = bs.ships[i].size
	}
	go bs.Run()
	log.Println("1")
	bs.SendMessage(messages.NewBoardMessage(shipSizes[:]))
	log.Println("2")
	return &bs
}

func NewBattleshipFromExisting(bs *Battleship, playerID string) *Battleship {
	nbs := Battleship{
		dimension:   bs.dimension,
		victory:     false,
		InMessages:  make(chan []byte, 10),
		OutMessages: make(chan OutMessage, 10),
		playerID:    playerID,
	}

	board := make([][]tile, len(bs.board))
	for i := 0; i < len(bs.board); i++ {
		board[i] = make([]tile, len(bs.board))
	}
	for i := 0; i < len(bs.board); i++ {
		for j := 0; j < len(bs.board[i]); j++ {
			board[i][j] = bs.board[i][j]
			if board[i][j].status == "fired" {
				board[i][j].status = "ship"
			}
		}
	}
	nbs.board = board

	for i := 0; i < len(bs.ships); i++ {
		nbs.ships[i] = ship{x: bs.ships[i].x, y: bs.ships[i].y, vertical: bs.ships[i].vertical, size: bs.ships[i].size}
		nbs.placeShip(&nbs.ships[i])
	}

	var shipSizes [len(bs.ships)]int
	for i := 0; i < len(bs.ships); i++ {
		shipSizes[i] = bs.ships[i].size
	}
	go nbs.Run()
	nbs.SendMessage(messages.NewBoardMessage(shipSizes[:]))
	return &nbs
}

func (b *Battleship) newBoard() {
	b.board = make([][]tile, b.dimension)
	for i := 0; i < b.dimension; i++ {
		b.board[i] = make([]tile, b.dimension)
	}
	for i := 0; i < b.dimension; i++ {
		for j := 0; j < b.dimension; j++ {
			b.board[i][j].status = "sea"
		}
	}

	for i := 0; i < len(b.ships); i++ {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		size := r1.Intn(3)
		b.ships[i] = newShip(size + 2)
	}

	b.victory = false

	for i := 0; i < len(b.ships); i++ {
		b.generateShip(&b.ships[i])
	}
	log.Print("Generated ships \n")
	for i := 0; i < len(b.ships); i++ {
		log.Printf("ship %+v\n", b.ships[i])
	}
}

func (b *Battleship) hasVictory() bool {
	result := true
	for _, s := range b.ships {
		result = result && s.isDestroyed()
	}
	return result
}

func (b *Battleship) generateShip(s *ship) {
	vertical, x, y := randomPlace(b.dimension, s.size)

	for !b.isValidPlacement(s.size, vertical, x, y) {
		vertical, x, y = randomPlace(b.dimension, s.size)
	}
	s.x = x
	s.y = y
	s.vertical = vertical
	b.placeShip(s)
}

func (b *Battleship) placeShip(s *ship) {
	if s.vertical {
		for i := 0; i < s.size; i++ {
			b.board[s.x][s.y+i].status = "ship"
			b.board[s.x][s.y+i].ship = s
		}
	} else {
		for i := 0; i < s.size; i++ {
			b.board[s.x+i][s.y].status = "ship"
			b.board[s.x+i][s.y].ship = s
		}

	}
}

func randomPlace(dimension int, size int) (bool, int, int) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	vertical := r1.Intn(100) < 50
	var x, y int
	if vertical {
		x = r1.Intn(dimension)
		y = r1.Intn(dimension - size)
	} else {
		x = r1.Intn(dimension - size)
		y = r1.Intn(dimension)
	}

	return vertical, x, y
}

func (b *Battleship) isValidPlacement(size int, vertical bool, x int, y int) bool {
	if vertical {
		if y-1 >= 0 && b.board[x][y-1].status != "sea" {
			return false
		}
		for i := 0; i < size; i++ {
			if b.board[x][y+i].status != "sea" {
				return false
			}
			if x-1 >= 0 && b.board[x-1][y+i].status != "sea" {
				return false
			}
			if x+1 < b.dimension && b.board[x+1][y+i].status != "sea" {
				return false
			}
		}
		if y+size < b.dimension && b.board[x][y+size].status != "sea" {
			return false
		}
	} else {
		if x-1 >= 0 && b.board[x-1][y].status != "sea" {
			return false
		}
		for i := 0; i < size; i++ {
			if b.board[x+i][y].status != "sea" {
				return false
			}
			if y-1 > 0 && b.board[x+i][y-1].status != "sea" {
				return false
			}
			if y+1 < b.dimension && b.board[x+i][y+1].status != "sea" {
				return false
			}
		}
		if x+size < b.dimension && b.board[x+size][y].status != "sea" {
			return false
		}

	}
	return true
}
