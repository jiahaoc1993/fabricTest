package main

import (
	//"fmt"
	pb "github.com/hyperledger/fabric/protos"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"log"
)

const (
	address = "172.22.28.123:7051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//fmt.Println("connect sucessfully!")
	c := pb.NewPeerClient(conn)
	_, err = c.ProcessTransaction(context.Background(), &pb.Transaction{})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
