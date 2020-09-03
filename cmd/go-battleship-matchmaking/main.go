package main

import (
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"open-match.dev/open-match/pkg/matchfunction"
	"open-match.dev/open-match/pkg/pb"

	"google.golang.org/grpc"
)

const (
	openMatchQuery = "om-query.open-match.svc.cluster.local:50503"
	// openMatchQuery = "localhost:50503"
	port = 50502
)

type mmf struct {
	qsc pb.QueryServiceClient
}

func main() {
	log.Println("Stargin match function")

	conn, err := grpc.Dial(openMatchQuery, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	mmf := &mmf{
		qsc: pb.NewQueryServiceClient(conn),
	}

	server := grpc.NewServer()
	pb.RegisterMatchFunctionServer(server, mmf)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen %v, %s", port, err.Error())
	}
	log.Println("pre serve")
	err = server.Serve(ln)
	if err != nil {
		log.Fatalf("failed serving %s", err.Error())
	}
	log.Println("post serve")
}

func (m *mmf) Run(req *pb.RunRequest, stream pb.MatchFunction_RunServer) error {
	log.Printf("Generating proposals for function %v", req.GetProfile().GetName())
	poolTickets, err := matchfunction.QueryPools(stream.Context(), m.qsc, req.GetProfile().GetPools())
	if err != nil {
		log.Printf("Failed to query tickets for the given pools, got %s", err.Error())
		return err
	}

	totalMatches := 0
	for pool, tickets := range poolTickets {
		log.Printf("making proposals for %s with %d number of tickets", pool, len(tickets))
		numberOfMatches := 0
		for i := 0; i+1 < len(tickets); i += 2 {
			proposal := &pb.Match{
				MatchId:       uuid.New().String(),
				MatchProfile:  req.Profile.Name,
				MatchFunction: "my-first-match-function",
				Tickets: []*pb.Ticket{
					tickets[i], tickets[i+1],
				},
			}
			numberOfMatches++
			totalMatches++

			err := stream.Send(&pb.RunResponse{Proposal: proposal})
			if err != nil {
				return err
			}
			log.Printf("Made %d matches for pool %s", numberOfMatches, pool)
		}
	}
	log.Printf("Done creating matches, made %d matches", totalMatches)
	return nil
}
