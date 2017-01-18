package main

import (
	"fmt"
	//"tool/loadKey"
	pb "github.com/hyperledger/fabric/protos"
	"tool/rpc"
)

func main() {
	response := rpc.Connect(&pb.Transaction{})
	fmt.Println(response)
}
