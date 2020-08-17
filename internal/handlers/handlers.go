package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/triberraar/go-battleship/internal/client"
	"github.com/triberraar/go-battleship/internal/match"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}} // use default options

func Battleship(mm *match.Matchmaker, w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	client := &client.Client{Conn: c, OutMessages: make(chan interface{}, 10), InMessages: make(chan []byte, 10), Username: r.URL.Query()["username"][0]}

	mm.Play(client, "battleships")

	go client.ReadPump()
	go client.WritePump()
}
