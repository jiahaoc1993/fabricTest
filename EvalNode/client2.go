package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/container"
	//"strings"
	//"github.com/hyperledger/fabric/core/chaincode"
	//"github.com/hyperledger/fabric/core/peer"
	//"github.com/spf13/viper"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
	//"errors"
	"golang.org/x/net/context"
	//"github.com/golang/protobuf/jsonpb"
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
			map[string]string{"path": "github.com/hyperledger/fabric/examples/chaincode/go/HelloWorld"},
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
func MakeAChaincodeSpec() (*pb.ChaincodeSpec, error) {
	var spec pb.ChaincodeSpec

	t := &params{
		1,
		map[string]string{"path": "github.com/hyperledger/fabric/examples/chaincode/go/HelloWorld"},
		ctorMsg{"init", []string{"Hello, World"}},
		"diego"}

	b, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("Error raised: %v", err)
		return nil, err
	}

	tmp, err := ioutil.ReadAll(bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Read error: %v", err)
		return nil, err
	}
	//fmt.Println(b, bytes.NewBuffer(b))
	err = json.Unmarshal(tmp, &spec)
	if err != nil {
		fmt.Printf("pb unmarshal error: %v", err)
		os.Exit(0)
	}
	return &spec, nil
}

func getChaincodeBytes(context context.Context, spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error) {
	var codePackageBytes []byte
	var err error
	codePackageBytes, err = container.GetChaincodePackageBytes(spec)
	if err != nil {
		return nil, err
	}
	chaincodeDeploymentSpec := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackageBytes}
	return chaincodeDeploymentSpec, nil
}

func main() {
	//for viper testing
	/*
		viper.SetEnvPrefix("core")
		viper.AutomaticEnv()
		replacer := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(replacer)

		fmt.Println(viper.GetBool("security.enabled"))
		fmt.Println(string(viper.GetString("peer.validator.consensus.plugin")))
	*/
	//fmt.Println(viper.GetString("chaincode.mode") == chaincode.DevModeUserRunsChaincode)
	//fmt.Println(viper.GetBool("security.privacy"))
	//fmt.Println(viper.GetBool("security.enabled"))

	//define the devop server
	//var serverDevops pb.DevopsServer
	//serverDevops = //use underlying *core.Devops
	//var spec pb.ChaincodeSpec
	//	t, err := MakeATransaction()
	spec, err := MakeAChaincodeSpec()
	if err != nil {
		os.Exit(0)
	}
	chaincodeDeploymentSpec, err := getChaincodeBytes(context.Background(), spec)
	if err != nil {
		fmt.Printf("error raised: %v", err)
	}

	fmt.Println(chaincodeDeploymentSpec)
}
