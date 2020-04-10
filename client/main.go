package main

import (
	//"io"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "kkagitala/go-rest-api/transport/pb"
)

const (
	address = "localhost:50051"
)

// createCustomer calls the RPC method CreateCustomer of CustomerServer
func createCustomer(client pb.SubscriptionClient, customer *pb.CreateRequest) {
	resp, err := client.Create(context.Background(), customer)
	if err != nil {
		log.Fatalf("Could not create Bill: %v", err)
	}
	if resp.Id > 0 {
		log.Printf("A new Bill has been added with id: %d", resp.Id)
	}
}

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// Creates a new CustomerClient
	client := pb.NewSubscriptionClient(conn)

	customer := &pb.CreateRequest{
		BillId:      200,
		PledgeId:    100,
		UserId:      101,
		CampaignId:  102,
		AmountCents: 103,
	}

	// Create a new customer
	createCustomer(client, customer)

}
