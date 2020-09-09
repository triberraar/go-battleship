package handlers

import (
	"log"
	"net/http"

	"github.com/triberraar/battleship/internal/game/battleship"

	"github.com/gorilla/websocket"
	"github.com/triberraar/battleship/internal/client"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}} // use default options

func Battleship(bs *battleship.BattleshipMatch, w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	client := &client.Client{Conn: c, OutMessages: make(chan interface{}, 10), InMessages: make(chan []byte, 10), Username: r.URL.Query()["username"][0]}

	// mm.Play(client, "battleships")
	if bs.ShouldRejoin(client.Username) {
		bs.Rejoin(client)
	} else {
		bs.Join(client)
	}

	go client.ReadPump()
	go client.WritePump()
}
