package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	addr string = "http://172.22.28.123:7050"
	info string = "823d33a2dce8bcfffffb102ffea6f37215426beba38c984721ff0c015625405b4fa99fafe829bc48a8d367d3b56147dd88909da02d97f590b3609e165a57a314"
)

type transmit struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  params `json:"params"`
	Id      string `json:"id"`
}

type params struct {
	Type          int               `json:"type"`
	ChaincodeID   map[string]string `json:"chaincodeID"`
	CtorMsg       ctorMsg           `json:"ctorMsg"`
	SecureContext string            `json:"secureContext"`
}

type ctorMsg struct {
	Function string   `json:"function"`
	Args     []string `json:"args"`
}

func RandomId() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return strconv.Itoa(r1.Intn(100000))
}

func Register() {
	loginInfo := []byte(`{
        "enrollId" : "diego",
        "enrollSecret" : "DRJ23pEQl16a"
    }`)
	res, err := http.Post(addr+"/registrar", "application/json", bytes.NewBuffer(loginInfo))
	if err != nil {
		fmt.Printf("Error raised: %v", err)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error raised: %v", err)
	}
	fmt.Println(string(b))
}

func Deploy(args []string) {
	t := &transmit{
		"2.0",
		"deploy",
		params{
			1,
			map[string]string{"path": "github.com/hyperledger/fabric/examples/chaincode/go/HelloWorld"},
			ctorMsg{"init", args},
			"diego"},
		RandomId(),
	}
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("error raise: %v", err)
	}
	fmt.Println(string(b))
	res, err := http.Post(addr+"/chaincode", "application/json", bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("error raise; %v", err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

func Invoke(args []string) {
	t := &transmit{
		"2.0",
		"invoke",
		params{
			1,
			map[string]string{"name": info},
			ctorMsg{"write", args},
			"diego"},
		RandomId(),
	}
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println("error raised: %v", err)
	}
	/*
		_, err = http.Post(addr+"/chaincode", "application/json", bytes.NewBuffer(b))
		if err != nil {
			fmt.Println("Error raised: %v", err)
			os.Exit(0)
		}*/

	res, err := http.Post(addr+"/chaincode", "application/json", bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("error raised: %v", err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

}

func Query(args []string) {
	t := &transmit{
		"2.0",
		"query",
		params{
			1,
			map[string]string{"name": info},
			ctorMsg{"read", args},
			"diego"},
		RandomId(),
	}
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println("error raised: %v", err)
	}
	resp, err := http.Post(addr+"/chaincode", "application/json", bytes.NewBuffer(b))
	//_, _ = http.Post(addr+"/chaincode", "application/json", bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("error raised: %v", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func main() {
	args := os.Args[1:]
	switch args[0] {
	case "register":
		Register()
	case "deploy":
		Deploy(args[1:])
	case "invoke":
		Invoke(args[1:])
	case "query":
		Query(args[1:])
	default:
		fmt.Println("use deploy/register")
	}
}
