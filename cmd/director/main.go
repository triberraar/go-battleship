package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"open-match.dev/open-match/pkg/pb"

	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	allocationv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"agones.dev/agones/pkg/client/clientset/versioned"
	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

const (
	openMatchBackEnd = "om-backend.open-match.svc.cluster.local:50505"
	// openMatchBackEnd = "localhost:50505"
	openMatchMatchmakingHost = "matchmaking.triberraar-mm.svc.cluster.local"
	// openMatchMatchmakingHost       = "192.168.1.171"
	openMatchMatchmakingPort int32 = 50502
)

type Prof struct {
	name         string
	matchProfile *pb.MatchProfile
	fleet        string
}

func initBattleshipNoobProfile(name string) Prof {
	pool := pb.Pool{
		Name:              name,
		TagPresentFilters: []*pb.TagPresentFilter{{Tag: "battleship"}},
		StringEqualsFilters: []*pb.StringEqualsFilter{
			{
				StringArg: "xp",
				Value:     "noob",
			},
		},
	}
	pools := []*pb.Pool{&pool}
	profile := &pb.MatchProfile{
		Name:  name,
		Pools: pools,
	}
	return Prof{
		name,
		profile,
		"battleship",
	}
}

func initBattleshipMasterProfile(name string) Prof {
	pool := pb.Pool{
		Name:              name,
		TagPresentFilters: []*pb.TagPresentFilter{{Tag: "battleship"}},
		StringEqualsFilters: []*pb.StringEqualsFilter{
			{
				StringArg: "xp",
				Value:     "master",
			},
		},
	}
	pools := []*pb.Pool{&pool}
	profile := &pb.MatchProfile{
		Name:  name,
		Pools: pools,
	}
	return Prof{
		name,
		profile,
		"battleship",
	}
}

func initrpsProfile(name string) Prof {
	pool := pb.Pool{
		Name:              name,
		TagPresentFilters: []*pb.TagPresentFilter{{Tag: "rps"}},
	}
	pools := []*pb.Pool{&pool}
	profile := &pb.MatchProfile{
		Name:  name,
		Pools: pools,
	}
	return Prof{
		name,
		profile,
		"rps",
	}
}

func initProfiles() map[string]Prof {
	result := make(map[string]Prof)
	result["battleship_noob"] = initBattleshipNoobProfile("battleship_noob")
	result["battleship_master"] = initBattleshipMasterProfile("battleship_master")
	result["rps"] = initrpsProfile("rps")
	return result
}

var profiles map[string]Prof

func main() {
	log.Println("Director doing directoring")
	profiles = initProfiles()

	agonesClient := createAgonesClient()

	for range time.Tick(time.Second * 10) {
		log.Println("tick")
		if err := run(agonesClient); err != nil {
			log.Println("I ran ;) into error", err.Error())
		}
	}
}

func createAgonesClient() *versioned.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	agonesClient, err := versioned.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return agonesClient
}

func run(agonesClient *versioned.Clientset) error {
	conn, err := grpc.Dial(openMatchBackEnd, grpc.WithInsecure())
	if err != nil {
		log.Println("failing to dail backend ", err)
		return err
	}
	defer conn.Close()
	bc := pb.NewBackendServiceClient(conn)

	fetchMatches(bc, agonesClient)

	return nil
}

func fetchMatches(bc pb.BackendServiceClient, agonesClient *versioned.Clientset) {
	// subroutine this stuff
	for _, p := range profiles {

		req := &pb.FetchMatchesRequest{
			Config: &pb.FunctionConfig{
				Host: openMatchMatchmakingHost,
				Port: openMatchMatchmakingPort,
				Type: pb.FunctionConfig_GRPC,
			},
			Profile: p.matchProfile,
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
			adr, err := allocateGameServer(agonesClient, resp.Match.MatchProfile)
			if err == nil {

				bc.AssignTickets(context.Background(), createAssignTicketRequest(resp.GetMatch(), adr))
			}

		}

	}
}

func allocateGameServer(agonesClient *versioned.Clientset, profile string) (string, error) {
	gsa, err := agonesClient.AllocationV1().GameServerAllocations("default").Create(
		&allocationv1.GameServerAllocation{
			Spec: allocationv1.GameServerAllocationSpec{
				Required: metav1.LabelSelector{
					MatchLabels: map[string]string{agonesv1.FleetNameLabel: profiles[profile].fleet},
				},
			},
		},
	)

	if err != nil {
		log.Printf("couldnt get me a server allocated %v", err)
		return "", fmt.Errorf("Failed to allocate server")
	}
	if gsa.Status.State != allocationv1.GameServerAllocationAllocated {
		log.Printf("server is not in allocated state %v", gsa.Status.State)
		return "", fmt.Errorf("Failed to allocate server")
	}
	log.Printf("gonna connect to server %s", gsa.Status.GameServerName)
	return fmt.Sprintf("%s:%d", gsa.Status.Address, gsa.Status.Ports[0].Port), nil
}

func createAssignTicketRequest(match *pb.Match, adr string) *pb.AssignTicketsRequest {
	tids := []string{}
	for _, t := range match.GetTickets() {
		tids = append(tids, t.GetId())
	}

	return &pb.AssignTicketsRequest{
		Assignments: []*pb.AssignmentGroup{
			{
				TicketIds: tids,
				Assignment: &pb.Assignment{
					Connection: adr,
				},
			},
		},
	}
}
