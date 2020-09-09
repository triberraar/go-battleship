package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	agones "agones.dev/agones/sdks/go"
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
	in     chan rpsMove
	out    map[string]chan result
	cancel map[string]chan bool
}

type myServer struct {
	http.Server
	shutdownReq     chan string
	rpsGameInstance *rpsGame
	agonesHealth    *agonesHealth
}

type agonesHealth struct {
	agones *agones.SDK
	stop   chan bool
}

func (ah *agonesHealth) doHealth() {
	tick := time.Tick(2 * time.Second)
	for {
		select {
		case <-ah.stop:
			log.Println("stopping healthcheck, going dowm")
			return
		case <-tick:
			if err := ah.agones.Health(); err != nil {
				log.Fatalf("Freaking failed the healtch %v", err)
			}
		}

	}
}

func (s *myServer) WaitShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	//Wait interrupt or shutdown request through /shutdown
	select {
	case sig := <-irqSig:
		log.Printf("Shutdown request through ctrl c (signal: %v)", sig)
		for k, c := range s.rpsGameInstance.cancel {
			log.Println("sending to somebody ", k)
			c <- true
		}
	case sig := <-s.shutdownReq:
		log.Printf("Shutdown request thourh logic( %v)", sig)
		if sig == "time" {
			for k, c := range s.rpsGameInstance.cancel {
				log.Println("sending to somebody ", k)
				c <- true
			}
		}
	}
	s.agonesHealth.stop <- true
	s.agonesHealth.agones.Shutdown()

	log.Printf("Stoping http server ...")

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//shutdown the server
	err := s.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutdown request error: %v", err)
	}
}

func (s *myServer) run() {
	log.Println("running")
	var u1 string
	var m1 string
	var u2 string
	var m2 string
	for {
		m := <-s.rpsGameInstance.in
		log.Println("got a message")
		if len(s.rpsGameInstance.out) == 1 {
			u1 = m.Username
			m1 = m.Move
			log.Println(u1)
			go s.cancelGame()
		} else {
			u2 = m.Username
			m2 = m.Move
			log.Println("2 messages")
			if m1 == m2 {
				s.rpsGameInstance.out[u1] <- result{"d"}
				s.rpsGameInstance.out[u2] <- result{"d"}
			} else if m1 == "r" && m2 == "p" {
				s.rpsGameInstance.out[u1] <- result{"l"}
				s.rpsGameInstance.out[u2] <- result{"w"}
			} else if m1 == "r" && m2 == "s" {
				s.rpsGameInstance.out[u1] <- result{"w"}
				s.rpsGameInstance.out[u2] <- result{"l"}
			} else if m1 == "p" && m2 == "r" {
				s.rpsGameInstance.out[u1] <- result{"w"}
				s.rpsGameInstance.out[u2] <- result{"l"}
			} else if m1 == "p" && m2 == "s" {
				s.rpsGameInstance.out[u1] <- result{"l"}
				s.rpsGameInstance.out[u2] <- result{"w"}
			} else if m1 == "s" && m2 == "r" {
				s.rpsGameInstance.out[u1] <- result{"l"}
				s.rpsGameInstance.out[u2] <- result{"w"}
			} else if m1 == "s" && m2 == "p" {
				s.rpsGameInstance.out[u1] <- result{"w"}
				s.rpsGameInstance.out[u2] <- result{"l"}
			} else {
				s.rpsGameInstance.out[u1] <- result{"d"}
				s.rpsGameInstance.out[u2] <- result{"d"}
			}
			delete(s.rpsGameInstance.out, u1)
			delete(s.rpsGameInstance.out, u2)
			u1 = ""
			m1 = ""
			u2 = ""
			m2 = ""
			s.shutdownReq <- "done"

		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10012"
	}

	flag.Parse()
	log.SetFlags(0)

	log.Printf("Server listening on %s for shure", port)

	agones, err := agones.NewSDK()
	if err != nil {
		log.Fatalf("agones sdk creation failed %v", err)
	}
	agonesHealth := &agonesHealth{
		agones: agones,
		stop:   make(chan bool),
	}

	addr := ":" + port
	server := &myServer{
		http.Server{Addr: addr},
		make(chan string),
		&rpsGame{
			make(chan rpsMove),
			make(map[string]chan result),
			make(map[string]chan bool),
		},
		agonesHealth,
	}

	router := mux.NewRouter()
	router.HandleFunc("/play", server.playrps).Methods("POST")
	handler := handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)
	server.Server.Handler = handler
	go server.run()
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

func (s *myServer) playrps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rpsMove rpsMove
	_ = json.NewDecoder(r.Body).Decode(&rpsMove)
	wc := make(chan result)
	cc := make(chan bool)
	s.rpsGameInstance.out[rpsMove.Username] = wc
	s.rpsGameInstance.cancel[rpsMove.Username] = cc
	s.rpsGameInstance.in <- rpsMove
	time.NewTicker(30000 * time.Millisecond)

	select {
	case res := <-wc:
		json.NewEncoder(w).Encode(res)
	case <-cc:
		log.Println("cancelling")
		json.NewEncoder(w).Encode(&result{"c"})
	}
}

func (s *myServer) cancelGame() {
	ticker := time.NewTicker(30000 * time.Millisecond)
	<-ticker.C

	log.Println("going byebye")
	go func() {
		s.shutdownReq <- "time"
	}()
}
