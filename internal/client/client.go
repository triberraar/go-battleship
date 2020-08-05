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
	OutMessages chan interface{}
	InMessages  chan ClientMessage
	Leaver      chan string
	Username    string
}

type ClientMessage struct {
	Username string
	Message  []byte
}

func (c *Client) ReadPump() {
	defer func() {
		if c.Username == "" {
			c.Leaver <- c.Username
		}
		c.Conn.Close()
	}()
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("readpump error: %v", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("unexcpected close")
			}
			break
		}
		bm := messages.BaseMessage{}
		json.Unmarshal(message, &bm)
		if bm.Type == "PING" {
		} else {
			c.InMessages <- ClientMessage{c.Username, message}
		}

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

func (c *Client) Close() {
	c.Conn.Close()
}
