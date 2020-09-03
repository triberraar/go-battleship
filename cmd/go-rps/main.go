package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type rpsMove struct {
	Move     string
	Username string
}

type result struct {
	Result string
}

type rpsGame struct {
	in  chan rpsMove
	out map[string]chan result
}

var rpsGameInstance = &rpsGame{
	make(chan rpsMove),
	make(map[string]chan result),
}

func (rps *rpsGame) run() {
	log.Println("running")
	var u1 string
	var m1 string
	var u2 string
	var m2 string
	for {
		m := <-rps.in
		log.Println("got a message")
		if len(rps.out) == 1 {
			u1 = m.Username
			m1 = m.Move
			log.Println(u1)
		} else {
			u2 = m.Username
			m2 = m.Move
			log.Println("2 messages")
			if m1 == m2 {
				rps.out[u1] <- result{"d"}
				rps.out[u2] <- result{"d"}
			} else if m1 == "r" && m2 == "p" {
				rps.out[u1] <- result{"l"}
				rps.out[u2] <- result{"w"}
			} else if m1 == "r" && m2 == "s" {
				rps.out[u1] <- result{"w"}
				rps.out[u2] <- result{"l"}
			} else if m1 == "p" && m2 == "r" {
				rps.out[u1] <- result{"w"}
				rps.out[u2] <- result{"l"}
			} else if m1 == "p" && m2 == "s" {
				rps.out[u1] <- result{"l"}
				rps.out[u2] <- result{"w"}
			} else if m1 == "s" && m2 == "r" {
				rps.out[u1] <- result{"l"}
				rps.out[u2] <- result{"w"}
			} else if m1 == "s" && m2 == "p" {
				rps.out[u1] <- result{"w"}
				rps.out[u2] <- result{"l"}
			} else {
				rps.out[u1] <- result{"d"}
				rps.out[u2] <- result{"d"}
			}
			delete(rps.out, u1)
			delete(rps.out, u2)
			u1 = ""
			m1 = ""
			u2 = ""
			m2 = ""

		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10012"
	}
	go rpsGameInstance.run()
	flag.Parse()
	log.SetFlags(0)
	router := mux.NewRouter()
	router.HandleFunc("/play", playrps).Methods("POST")
	log.Printf("Server listening on %s for shure", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))

}

func playrps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rpsMove rpsMove
	_ = json.NewDecoder(r.Body).Decode(&rpsMove)
	wc := make(chan result)
	rpsGameInstance.out[rpsMove.Username] = wc
	rpsGameInstance.in <- rpsMove
	time.NewTicker(30000 * time.Millisecond)

	res := <-wc
	json.NewEncoder(w).Encode(res)
}
