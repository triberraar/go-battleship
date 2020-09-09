package main

import (
	"flag"
	"log"
	"net/http"
	"os"

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
	router := mux.NewRouter()
	bs := battleship.NewBattleshipMatch(2)
	router.HandleFunc("/battleship", func(w http.ResponseWriter, r *http.Request) {
		bsHandlers.Battleship(bs, w, r)
	})
	log.Printf("Server listening on %s for shure", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
