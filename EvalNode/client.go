package main

import (
	"fmt"
	//"errors"
	//"golang.org/x/net/context"
	"github.com/golang/protobuf/jsonpb"
	pb "github.com/hyperledger/fabric/protos"


type Transaction struct {
	Jsonrpc *string				`json:"jsonrpc,omitempty"`
	Method  *string				`json:"method,omitempty"`
	Params  *pb.ChaincodeSpec   `json:"params,omitempty"`
	ID      *rpcID
	}

type rpcID struct {
	StringValue *string
	IntValue    *int64
	}

func RandomId() string

func MakeATransaction() *Transaction {
	t := &Transaction{
		"2.0",
		"deploy",
		Params{
			1,
			pb.ChaincodeID{"github.com/hyperledger/fabric/examples/chaincode/go/Hello_World"},
			pb.ChainInput{"Hello", []string{"abc"}},
			"diego"},
		rpcID{"id":RandomId()},
	}
	return t
}


/*
func Deploy() {
	
	}
*/

func main() {
	var spec pb.ChaincodeSpec
	t := MakeATransaction()
	err := jsonpb.Unmarshal(t, &spec)
	if err != nil {
		fmt.Println("error raised: %v", err)
	}
}
