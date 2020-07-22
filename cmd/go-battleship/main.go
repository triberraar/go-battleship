package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	bsHandlers "github.com/triberraar/go-battleship/internal/handlers"
	"github.com/triberraar/go-battleship/internal/room"
)

var addr = flag.String("addr", "localhost:10002", "http service address")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10002"
	}
	flag.Parse()
	log.SetFlags(0)
	router := mux.NewRouter()
	rm := room.NewRoomManager()
	router.HandleFunc("/battleship", func(w http.ResponseWriter, r *http.Request) {
		bsHandlers.Battleship(rm, w, r)
	})
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("dist/"))))
	log.Printf("Server listening on %s for shure", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
