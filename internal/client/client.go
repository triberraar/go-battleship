package client

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/triberraar/go-battleship/internal/messages"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}} // use default options

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Conn        *websocket.Conn
	PlayerID    string
	OutMessages chan interface{}
	InMessages  chan ClientMessage
}

type ClientMessage struct {
	PlayerID string
	Message  []byte
}

func (c *Client) ReadPump() {
	defer c.Conn.Close()
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("unexcpected close")
			}
			break
		}
		bm := messages.BaseMessage{}
		json.Unmarshal(message, &bm)
		c.InMessages <- ClientMessage{c.PlayerID, message}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	for {
		select {
		case message := <-c.OutMessages:
			c.Conn.WriteJSON(message)
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
