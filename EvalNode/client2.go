package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
	//"errors"
	//"golang.org/x/net/context"
	"github.com/golang/protobuf/jsonpb"
	pb "github.com/hyperledger/fabric/protos"
)

type Transaction struct {
	Jsonrpc string `json:"jsonrpc,omitempty"`
	Method  string `json:"method,omitempty"`
	Params  params `json:"params,omitempty"`
	ID      int    `json:"id,omitempty"`
}

type params struct {
	Type          int               `json:"type,omitempty"`
	ChaincodeID   map[string]string `json:"chaincodeID,omitempty"`
	CtorMsg       ctorMsg           `json:"ctorMsg"`
	SecureContext string            `josn:"secureContext,omitempty"`
}

type ctorMsg struct {
	Function string   `json:"function,omitempty"`
	Args     []string `json:"args,omitempty"`
}

func RandomId() int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	return r.Intn(1000000)
}

// this is for normal resp trasnaction upon http
func MakeATransaction() (*bytes.Buffer, error) {
	t := &Transaction{
		"2.0",
		"deploy",
		params{
			1,
			map[string]string{"path": "github.com/hyperledger/fabric/examples/chaincode/go/Hello_World"},
			ctorMsg{"init", []string{"Hello", "World"}},
			"diego"},
		RandomId(),
	}
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("error raised: %v\n", err)
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

// this is only for pb.chaincodSpec
func MakeAChaincodeSpec() (*bytes.Buffer, error) {
	t := &params{
		1,
		map[string]string{"path": "github.com/hyperledger/fabric/example/chaincode/go/Hello_World"},
		ctorMsg{"init", []string{"Hello", "World"}},
		"diego",
	}
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("Error raised: %v", err)
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

/*
func Deploy() {

	}
*/

func main() {
	var spec pb.ChaincodeSpec
	//	t, err := MakeATransaction()
	t, err := MakeAChaincodeSpec()
	if err != nil {
		os.Exit(0)
	}
	err = jsonpb.Unmarshal(t, &spec)
	if err != nil {
		fmt.Printf("f error raised: %v\n", err)
	}
}
