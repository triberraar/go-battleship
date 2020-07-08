package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/triberraar/go-battleship/internal/messages"
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
	for {
		bm := messages.BaseMessage{}
		_, message, _ := c.ReadMessage()
		json.Unmarshal(message, &bm)
		if bm.Type == "FIRE" {
			fm := messages.FireMessage{}
			json.Unmarshal(message, &fm)
			fmt.Printf("%+v\n", fm)
			if fm.Coordinate.X%2 == 0 {
				c.WriteJSON(newHitMessage(fm.Coordinate))
			} else {
				c.WriteJSON(newMissMessage(fm.Coordinate))
			}
		}
	}
}

func newHitMessage(coordinate messages.Coordinate) messages.HitMessage {
	return messages.HitMessage{
		Type:       "HIT",
		Coordinate: coordinate,
	}
}

func newMissMessage(coordinate messages.Coordinate) messages.MissMessage {
	return messages.MissMessage{
		Type:       "MISS",
		Coordinate: coordinate,
	}
}
