package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"open-match.dev/open-match/pkg/pb"

	"google.golang.org/grpc"
)

const (
	openMatchBackEnd = "om-backend.open-match.svc.cluster.local:50505"
	// openMatchBackEnd = "localhost:50505"
	openMatchMatchmakingHost = "go-battleship-matchmaking.triberraar-mm.svc.cluster.local"
	// openMatchMatchmakingHost       = "192.168.1.171"
	openMatchMatchmakingPort int32 = 50502
)

func main() {
	log.Println("Director doing directoring")
	for range time.Tick(time.Second * 10) {
		log.Println("tick")
		if err := run(); err != nil {
			log.Println("I ran ;) into error", err.Error())
		}
	}
}

func run() error {
	conn, err := grpc.Dial(openMatchBackEnd, grpc.WithInsecure())
	if err != nil {
		log.Println("failing to dail backend ", err)
		return err
	}
	defer conn.Close()
	bc := pb.NewBackendServiceClient(conn)

	// bc.FetchMatches(context.Background(), fetchMatchesRequest())
	fetchMatches(bc)

	return nil
}

func fetchMatches(bc pb.BackendServiceClient) {
	var profiles []*pb.MatchProfile
	// games := []string{"battleships", "rps"}
	xps := []string{"noob", "master"}
	// for _, game := range games {
	for _, xp := range xps {
		var pools []*pb.Pool
		pools = append(pools, &pb.Pool{
			Name:              fmt.Sprintf("pool_%s_%s", "battleships", xp),
			TagPresentFilters: []*pb.TagPresentFilter{{Tag: "battleships"}},
			StringEqualsFilters: []*pb.StringEqualsFilter{
				{
					StringArg: "xp",
					Value:     xp,
				},
			},
		})
		profiles = append(profiles, &pb.MatchProfile{
			Name:  fmt.Sprintf("Profile_%s_%s", "battleships", xp),
			Pools: pools,
		})
		// }
	}

	var pools []*pb.Pool
	pools = append(pools, &pb.Pool{
		Name:              fmt.Sprintf("pool_%s_%s", "rps", "blaat"),
		TagPresentFilters: []*pb.TagPresentFilter{{Tag: "rps"}},
	})
	profiles = append(profiles, &pb.MatchProfile{
		Name:  fmt.Sprintf("Profile_%s_%s", "rps", "blaat"),
		Pools: pools,
	})

	// subroutine this stuff
	for _, p := range profiles {

		req := &pb.FetchMatchesRequest{
			Config: &pb.FunctionConfig{
				Host: openMatchMatchmakingHost,
				Port: openMatchMatchmakingPort,
				Type: pb.FunctionConfig_GRPC,
			},
			Profile: p,
		}

		stream, err := bc.FetchMatches(context.Background(), req)
		if err != nil {
			log.Println("error fetching matches ", err)
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("error receiving from stream", err)
				break
			}
			bc.AssignTickets(context.Background(), createAssignTicketRequest(resp.GetMatch()))

		}

	}
}

func createAssignTicketRequest(match *pb.Match) *pb.AssignTicketsRequest {
	tids := []string{}
	for _, t := range match.GetTickets() {
		tids = append(tids, t.GetId())
	}

	host := getServerFromProfile(match.GetMatchProfile())
	return &pb.AssignTicketsRequest{
		Assignments: []*pb.AssignmentGroup{
			{
				TicketIds: tids,
				Assignment: &pb.Assignment{
					Connection: host,
				},
			},
		},
	}
}

// my super agones stub
func getServerFromProfile(profile string) string {
	splitted := strings.Split(profile, "_")
	if len(splitted) != 3 {
		return ""
	}
	if splitted[1] == "battleships" && splitted[2] == "noob" {
		return "localhost:10003"
	} else if splitted[1] == "battleships" && splitted[2] == "master" {
		return "localhost:10004"
	} else if splitted[1] == "rps" {
		return "localhost:10012"
	} else {
		return ""
	}
}
