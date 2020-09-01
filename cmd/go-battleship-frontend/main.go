package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	router.HandleFunc("/play", play).Methods("GET")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("dist/"))))
	log.Printf("Server listening on %s for shure", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}

type playMessage struct {
	URL string
}

func play(w http.ResponseWriter, r *http.Request) {
	// do the frontend bit of matchmaking, return hard coded now
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&playMessage{URL: "localhost:10003"})
}
