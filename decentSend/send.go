package main

import (
    "fmt"
    "net/http"
    "os"
    "io/ioutil"
    "encoding/json"
    "bytes"
)

const (
    url string = "http://172.22.28.118"
    port string = "7051"
)

type transmit struct {
    Jsonrpc string      `json:"jsonrpc"`
    Method  string      `json:"method"`
    Params  map[string]interface{}  `json:"params"`
    Id      int         `json:"id"`
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

func Register(){
    loginInfo := []byte(`{
        "enrollId" : "jim",
        "enrollSecret" : "abcdefg"
    }`)
    res, err := http.Post("http://172.22.28.130:7050/registrar", "application/json", bytes.NewBuffer(loginInfo))
    if err != nil {
        fmt.Printf("Error raised: %v", err)
    }

}



func Deploy(){
    t := &transmit{
        "2.0",
        "deploy",
        map[string]interface{}{
            "chaincodeID" : map[string]interface{}{"path": "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02"},
            "ctorMsg" : map[string]interface{}{"function" : "deploy","args": []string{"a", "10000", "b", "1 00"}},
            "secureContext" : "lukas"},
        1,
    }
    b, err := json.Marshal(t)
    if err != nil {
        fmt.Printf("error raise: %v", err)
    }
    res, err := http.Post("http://172.22.28.130:7050/chaincode", "application/json", bytes.NewBuffer(b))
    if err != nil {
        fmt.Println("error raise; %v", err)
    }
    body, _ := ioutil.ReadAll(res.Body)
    fmt.Println(string(body))
}


func main() {
    arg := os.Args[1]
    if arg != "post" {
        addr := url + ":" + port
        err := Get(addr+arg)
        if err != nil {
            fmt.Printf("error raised")
        }
    }

}

