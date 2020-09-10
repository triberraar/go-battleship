package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	agonesSDK "agones.dev/agones/sdks/go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/triberraar/battleship/internal/game/battleship"
	bsHandlers "github.com/triberraar/battleship/internal/handlers"
)

var addr = flag.String("addr", "localhost:10002", "http service address")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10003"
	}
	flag.Parse()
	log.SetFlags(0)

	agones, err := agonesSDK.NewSDK()
	if err != nil {
		log.Fatalf("agones sdk creation failed %v", err)
	}
	agonesHealth := &agonesHealth{
		agones: agones,
		stop:   make(chan bool),
	}

	server := &myServer{
		http.Server{Addr: ":" + port},
		make(chan string),
		agonesHealth,
	}

	router := mux.NewRouter()
	bs := battleship.NewBattleshipMatch(2, server.shutdownReq)
	router.HandleFunc("/battleship", func(w http.ResponseWriter, r *http.Request) {
		bsHandlers.Battleship(bs, w, r)
	})
	log.Printf("Server listening on %s for shure", port)
	handler := handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)
	server.Server.Handler = handler

	done := make(chan bool)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Listen and serve: %v", err)
		}
		done <- true
	}()
	agones.Ready()
	go agonesHealth.doHealth()

	//wait shutdown
	server.WaitShutdown()
	<-done
	log.Printf("DONE!")
}
