package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"open-match.dev/open-match/pkg/pb"
)

var addr = flag.String("addr", "localhost:10002", "http service address")

const (
	openMatchFrontEnd = "om-frontend.open-match.svc.cluster.local:50504"
	// openMatchFrontEnd = "localhost:50504"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10002"
	}
	flag.Parse()
	log.SetFlags(0)
	router := mux.NewRouter()
	router.HandleFunc("/battleships/play", playBattleships).Methods("GET")
	router.HandleFunc("/rps/play", playRps).Methods("GET")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("dist/"))))
	log.Printf("Server listening on %s for shure", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}

type playMessage struct {
	URL string
}

type errorMessage struct {
	Error string
}

func playRps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("connecting")
	conn, err := grpc.Dial(openMatchFrontEnd, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to contact open match %v", err)
	}
	log.Println("connected")
	defer conn.Close()
	fe := pb.NewFrontendServiceClient(conn)
	req := &pb.CreateTicketRequest{
		Ticket: &pb.Ticket{
			SearchFields: &pb.SearchFields{
				Tags: []string{
					"rps",
				},
			},
		},
	}
	resp, err := fe.CreateTicket(context.Background(), req)
	if err != nil {
		log.Printf("failed to create ticket, got %s", err.Error())
	}
	ticketId := resp.Id
	assignment := make(chan *pb.Assignment)
	errch := make(chan string)

	defer deleteTicket(fe, ticketId)

	ctx, cancel := context.WithCancel(r.Context())
	go streamAssignment(ctx, fe, ticketId, assignment, errch)
	ticker := time.NewTicker(30000 * time.Millisecond)
	for {
		select {
		case err := <-errch:
			log.Println("something on the errorchannel ", err)
			json.NewEncoder(w).Encode(&errorMessage{err})
			return
		case assignment := <-assignment:
			log.Println("got an assignement ", assignment)
			// json.NewEncoder(w).Encode(&playMessage{URL: "localhost:10003"})
			json.NewEncoder(w).Encode(&playMessage{URL: assignment.Connection})
			return
		case <-ticker.C:
			log.Println("no assignment")
			cancel()
			json.NewEncoder(w).Encode(&errorMessage{"no assignemnt"})
			ticker.Stop()
			return
		}
	}
}

func playBattleships(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("connecting")
	conn, err := grpc.Dial(openMatchFrontEnd, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to contact open match %v", err)
	}
	log.Println("connected")
	defer conn.Close()
	fe := pb.NewFrontendServiceClient(conn)
	req := &pb.CreateTicketRequest{
		Ticket: &pb.Ticket{
			SearchFields: &pb.SearchFields{
				Tags: []string{
					"battleship",
				},
				StringArgs: map[string]string{
					"xp": xp(r.URL.Query()["username"][0]),
				},
			},
		},
	}
	resp, err := fe.CreateTicket(context.Background(), req)
	if err != nil {
		log.Printf("failed to create ticket, got %s", err.Error())
	}
	ticketId := resp.Id
	assignment := make(chan *pb.Assignment)
	errch := make(chan string)

	defer deleteTicket(fe, ticketId)

	ctx, cancel := context.WithCancel(r.Context())
	go streamAssignment(ctx, fe, ticketId, assignment, errch)
	ticker := time.NewTicker(30000 * time.Millisecond)
	for {
		select {
		case err := <-errch:
			log.Println("something on the errorchannel ", err)
			json.NewEncoder(w).Encode(&errorMessage{err})
			return
		case assignment := <-assignment:
			log.Println("got an assignement ", assignment)
			// json.NewEncoder(w).Encode(&playMessage{URL: "localhost:10003"})
			json.NewEncoder(w).Encode(&playMessage{URL: assignment.Connection})
			return
		case <-ticker.C:
			log.Println("no assignment")
			cancel()
			json.NewEncoder(w).Encode(&errorMessage{"no assignemnt"})
			ticker.Stop()
			return
		}
	}
	// json.NewEncoder(w).Encode(&playMessage{URL: "localhost:10003"})
}

func deleteTicket(fe pb.FrontendServiceClient, ticketId string) {
	_, err := fe.DeleteTicket(context.Background(), &pb.DeleteTicketRequest{TicketId: ticketId})
	if err != nil {
		log.Printf("error deleteing ticket %s %v", ticketId, err)
	}
}

func streamAssignment(ctx context.Context, fe pb.FrontendServiceClient, ticketId string, assingment chan *pb.Assignment, errch chan string) {
	stream, err := fe.WatchAssignments(ctx, &pb.WatchAssignmentsRequest{TicketId: ticketId})
	if err != nil {
		log.Printf("error streaming assingment for %s: %v", ticketId, err)
		errch <- "error streaming assignemnt"
		return
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("error getting from stream for %s, %v", ticketId, err)
			errch <- "error getting from stream"
			return
		}
		log.Println("assigned %v", resp.Assignment)
		assingment <- resp.Assignment
	}
	log.Println("Done")
}

func xp(username string) string {
	if len(username) <= 3 {
		return "noob"
	} else {
		return "master"
	}
}
