package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
	"github.com/triberraar/go-battleship/internal/client"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}} // use default options

// Battleship the handlers for the battleship socket stuff
func Battleship(rm *client.RoomManager, w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	pid := uuid.New()
	client := &client.Client{Conn: c, PlayerID: pid.String(), OutMessages: make(chan interface{}, 10)}

	go client.ReadPump()
	go client.WritePump()
}
