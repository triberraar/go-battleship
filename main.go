package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type baseMessage struct {
	Type string `json:"type"`
}

type coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type hitMessage struct {
	Type       string     `json:"type"`
	Coordinate coordinate `json:"coordinate"`
}

type missMessage struct {
	Type       string     `json:"type"`
	Coordinate coordinate `json:"coordinate"`
}

type fireMessage struct {
	Type       string     `json:"type"`
	Coordinate coordinate `json:"coordinate"`
}

func newHitMessage(coordinate coordinate) hitMessage {
	return hitMessage{
		Type:       "HIT",
		Coordinate: coordinate,
	}
}

func newMissMessage(coordinate coordinate) missMessage {
	return missMessage{
		Type:       "MISS",
		Coordinate: coordinate,
	}
}

var addr = flag.String("addr", "localhost:10002", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}} // use default options

func battleship(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		bm := baseMessage{}
		_, message, _ := c.ReadMessage()
		json.Unmarshal(message, &bm)
		if bm.Type == "FIRE" {
			fm := fireMessage{}
			json.Unmarshal(message, &fm)
			fmt.Printf("%+v\n", fm)
			if fm.Coordinate.X%2 == 0 {
				c.WriteJSON(newHitMessage(fm.Coordinate))
			} else {
				c.WriteJSON(newMissMessage(fm.Coordinate))
			}
		}
		// mt, message, _ := c.ReadMessage()
		// c.ReadJSON()
		// if mt == websocket.TextMessage {
		// 	log.Print("Text message")
		// 	log.Printf("%s", message)
		// 	var bm baseMessage
		// 	json.Unmarshal(message, &bm)
		// 	log.Print(bm.Data)
		// 	log.Print(bm.Type)
		// 	log.Print(bm.FS)
		// }
		// if err != nil {
		// 	log.Println("read:", err)
		// 	break
		// }
		// log.Printf("recv: %s", message)
		// err = c.WriteMessage(mt, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10002"
	}
	flag.Parse()
	log.SetFlags(0)
	router := mux.NewRouter()
	router.HandleFunc("/battleship", battleship)
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))
	log.Printf("Server listening on %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
