package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"agones.dev/agones/pkg/sdk"
	agonesSDK "agones.dev/agones/sdks/go"
)

type myServer struct {
	http.Server
	shutdownReq  chan string
	agonesHealth *agonesHealth
}

type agonesHealth struct {
	agones *agonesSDK.SDK
	stop   chan bool
}

func (ah *agonesHealth) doHealth() {
	tick := time.Tick(2 * time.Second)
	for {
		<-tick
		if err := ah.agones.Health(); err != nil {
			log.Fatalf("Freaking failed the healtch %v", err)
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
	case sig := <-s.shutdownReq:
		log.Printf("Shutdown request thourh logic( %v)", sig)
	}
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

func (s *myServer) watch(gs *sdk.GameServer) {
	if gs.Status.State == "Allocated" {
		log.Println("server is allocated")
		log.Println("some metatdata")
		for k, v := range gs.ObjectMeta.Labels {
			log.Printf("%s:%s", k, v)
		}
	}
}
