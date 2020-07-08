package game_logic

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/websocket"
	"github.com/triberraar/go-battleship/internal/messages"
)

type battleship struct {
	c     *websocket.Conn
	board [][]byte
}

func RunBattleship(c *websocket.Conn) {
	bs := battleship{c: c}
	bs.newBoard()
	for {
		bm := messages.BaseMessage{}
		_, message, _ := c.ReadMessage()
		json.Unmarshal(message, &bm)
		if bm.Type == "FIRE" {
			fm := messages.FireMessage{}
			json.Unmarshal(message, &fm)
			if bs.board[fm.Coordinate.X][fm.Coordinate.Y] == 'b' {
				c.WriteJSON(messages.NewHitMessage(fm.Coordinate))
				bs.board[fm.Coordinate.X][fm.Coordinate.Y] = 'f'
			} else if bs.board[fm.Coordinate.X][fm.Coordinate.Y] == 's' {
				c.WriteJSON(messages.NewMissMessage(fm.Coordinate))
				bs.board[fm.Coordinate.X][fm.Coordinate.Y] = 'f'
			}
		}
	}
}

func (b *battleship) newBoard() {
	b.board = make([][]byte, 10)
	for i := 0; i < 10; i++ {
		b.board[i] = make([]byte, 10)
	}
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if rand.Intn(100) < 50 {
				b.board[i][j] = 's'
			} else {
				b.board[i][j] = 'b'
			}
		}
	}

}
