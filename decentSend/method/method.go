package method

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
	info string = "aa5552158725fda5bc0764d0aaa7ccb31e0887a6e10b1c773f586f57deb5e05d8f5400c26ed9e401d865b52159e203a987c2b0b99311445830f49713c3cf080b"
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

func Get(s string) error {
	resp, err := http.Get(s)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("Error: %v", err)
		return err
	}
	fmt.Printf("Message : %s", b)
	return nil
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

func Deploy() {
	t := &transmit{
		"2.0",
		"deploy",
		params{
			1,
			map[string]string{"path": "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02"},
			ctorMsg{"init", []string{"a", "10000000", "b", "0"}},
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

func Invoke() {
	t := &transmit{
		"2.0",
		"invoke",
		params{
			1,
			map[string]string{"name": info},
			ctorMsg{"invoke", []string{"a", "b", "1"}},
			"diego"},
		RandomId(),
	}
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println("error raised: %v", err)
	}
	_, err = http.Post(addr+"/chaincode", "application/json", bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Error raised: %v", err)
		os.Exit(0)
	}
	/*
		res, err := http.Post(addr+"/chaincode", "application/json", bytes.NewBuffer(b))
		if err != nil {
			fmt.Println("error raised: %v", err)
		}
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
	*/
}

func Query() {
	t := &transmit{
		"2.0",
		"query",
		params{
			1,
			map[string]string{"name": info},
			ctorMsg{"query", []string{"a"}},
			"diego"},
		RandomId(),
	}
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println("error raised: %v", err)
	}
	//resp, err := http.Post(addr+"/chaincode", "application/json", bytes.NewBuffer(b))
	_, _ = http.Post(addr+"/chaincode", "application/json", bytes.NewBuffer(b))
	//if err != nil {
	//fmt.Println("error raised: %v", err)
	//}
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
}

func main() {
	arg := os.Args[1]
	switch arg {
	case "register":
		Register()
	case "deploy":
		Deploy()
	case "invoke":
		Invoke()
	case "query":
		Query()
	default:
		fmt.Println("use deploy/register")
	}
}
