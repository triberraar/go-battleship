package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	gl "github.com/triberraar/go-battleship/internal/game_logic"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}} // use default options

// Battleship the handlers for the battleship socket stuff
func Battleship(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	gl.RunBattleship(c)
}
