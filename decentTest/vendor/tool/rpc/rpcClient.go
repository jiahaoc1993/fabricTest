package rpc

import (
	"time"
	//"math/rand"
	"fmt"
	"github.com/hyperledger/fabric/core/peer"
	pb "github.com/hyperledger/fabric/protos"
	context "golang.org/x/net/context"
	//grpc "google.golang.org/grpc"
	"log"
)

var vps  = []string {"172.22.28.134:7051", "172.22.28.178:7051","172.22.28.141:7051","172.22.28.144:7051"}


func RandomConnect(tx *pb.Transaction) ([]float64) {
	//source := rand.NewSource(time.Now().UnixNano())
	//r := rand.New(source)
	//target := r.Intn(4)
	var spents   []float64
	for i :=0 ;i< 4 ; i++ {
		start := time.Now().UnixNano()
		_ = Connect(tx, vps[i])
		end  := time.Now().UnixNano()
		spent := float64(end - start) / 1000000000
		spents = append(spents, spent)
	}
	return spents
}

func Connect(tx *pb.Transaction, address string) (*pb.Response) {
	conn, err := peer.NewPeerClientConnectionWithAddress(address)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//fmt.Println("connect sucessfully!")
	c := pb.NewPeerClient(conn)
	response, err := c.ProcessTransaction(context.Background(), tx)
	if err != nil {
		//log.Fatalf("Error: %v", err)
		return &pb.Response{Status: pb.Response_FAILURE, Msg: []byte(fmt.Sprintf("Error calling ProcessTransction on remote peer =%s : %v", address, err))}
	}
	return response
}
