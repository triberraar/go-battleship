package game_logic

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/triberraar/go-battleship/internal/messages"
)

type ship struct {
	x        int
	y        int
	size     int
	vertical bool
	hits     int
}

func newShip4() ship {
	return ship{
		size: 4,
	}
}

func newShip3() ship {
	return ship{
		size: 3,
	}
}

func newShip2() ship {
	return ship{
		size: 2,
	}
}

func newShip1() ship {
	return ship{
		size: 1,
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

type battleship struct {
	c         *websocket.Conn
	board     [][]tile
	dimension int
	ships     [6]ship
	victory   bool
}

func RunBattleship(c *websocket.Conn) {
	bs := battleship{c: c, dimension: 10, victory: false}
	bs.newBoard()
	for {
		bm := messages.BaseMessage{}
		_, message, _ := c.ReadMessage()
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
					c.WriteJSON(messages.NewShipDestroyedMessage(coordinate, bs.board[fm.Coordinate.X][fm.Coordinate.Y].ship.size, bs.board[fm.Coordinate.X][fm.Coordinate.Y].ship.vertical))
				} else {
					c.WriteJSON(messages.NewHitMessage(fm.Coordinate))
				}
				if bs.hasVictory() {
					bs.victory = true
					c.WriteJSON(messages.NewVictoryMessage())
					select {
					case <-time.After(10 * time.Second):
						bs.newBoard()
						c.WriteJSON(messages.NewRestartMessage())
					}
				}
			} else {
				c.WriteJSON(messages.NewMissMessage(fm.Coordinate))
				bs.board[fm.Coordinate.X][fm.Coordinate.Y].status = "fired"
			}
		}
	}
}

func (b *battleship) newBoard() {
	b.board = make([][]tile, b.dimension)
	for i := 0; i < b.dimension; i++ {
		b.board[i] = make([]tile, b.dimension)
	}
	for i := 0; i < b.dimension; i++ {
		for j := 0; j < b.dimension; j++ {
			b.board[i][j].status = "sea"
		}
	}

	b.ships[0] = newShip4()
	b.ships[1] = newShip3()
	b.ships[2] = newShip2()
	b.ships[3] = newShip2()
	b.ships[4] = newShip1()
	b.ships[5] = newShip1()

	for i := 0; i < len(b.ships); i++ {
		b.generateShip(&b.ships[i])
	}
	log.Print("Generated ships \n")
	for i := 0; i < len(b.ships); i++ {
		log.Printf("ship %+v\n", b.ships[i])
	}
}

func (b *battleship) hasVictory() bool {
	result := true
	for _, s := range b.ships {
		result = result && s.isDestroyed()
	}
	return result
}

func (b *battleship) generateShip(s *ship) {
	vertical, x, y := randomPlace(b.dimension, s.size)

	for !b.isValidPlacement(s.size, vertical, x, y) {
		vertical, x, y = randomPlace(b.dimension, s.size)
	}
	s.x = x
	s.y = y
	s.vertical = vertical
	if vertical {
		for i := 0; i < s.size; i++ {
			b.board[x][y+i].status = "ship"
			b.board[x][y+i].ship = s
		}
	} else {
		for i := 0; i < s.size; i++ {
			b.board[x+i][y].status = "ship"
			b.board[x+i][y].ship = s
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

func (b *battleship) isValidPlacement(size int, vertical bool, x int, y int) bool {
	if vertical {
		for i := 0; i < size; i++ {
			if b.board[x][y+i].status != "sea" {
				return false
			}
		}
	} else {
		for i := 0; i < size; i++ {
			if b.board[x+i][y].status != "sea" {
				return false
			}

		}

	}
	return true
}
