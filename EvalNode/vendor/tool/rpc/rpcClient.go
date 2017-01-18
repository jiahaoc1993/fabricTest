package rpc

import (
	"fmt"
	"github.com/hyperledger/fabric/core/peer"
	pb "github.com/hyperledger/fabric/protos"
	context "golang.org/x/net/context"
	//grpc "google.golang.org/grpc"
	"log"
)

const (
	peerAddress = "172.22.28.123:7051"
)

func Connect(tx *pb.Transaction) (response *pb.Response) {
	conn, err := peer.NewPeerClientConnectionWithAddress(peerAddress)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//fmt.Println("connect sucessfully!")
	c := pb.NewPeerClient(conn)
	response, err = c.ProcessTransaction(context.Background(), tx)
	if err != nil {
		//log.Fatalf("Error: %v", err)
		return &pb.Response{Status: pb.Response_FAILURE, Msg: []byte(fmt.Sprintf("Error calling ProcessTransction on remote peer =%s : %v", peerAddress, err))}
	}
	return response
}
