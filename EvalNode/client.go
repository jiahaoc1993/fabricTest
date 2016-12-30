package main

import (
	"fmt"
	"errors"
	"golang.org/x/net/context"
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

func MakeATransaction() {
	t := &Transaction{
		"2.0",
		"deploy",
		Params{
			1,

		}
	}

	}



func Deploy() {

	}

func main() {
	fmt.Println("vim-go")
}
