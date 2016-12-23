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
    addr string = "http://172.22.28.118:7050"
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

func Deploy(){
    t := &transmit{
        "2.0",
        "deploy",
        map[string]interface{}{
            "chaincodeID" : map[string]interface{}{"path": "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02"},
            "ctorMsg" : map[string]interface{}{"function" : "deploy","args": []string{"a", "10000000", "b", "0"}},
            "secureContext" : "jim"},
        1,
    }
    b, err := json.Marshal(t)
    if err != nil {
        fmt.Printf("error raise: %v", err)
    }
    res, err := http.Post(addr+"/chaincode", "application/json", bytes.NewBuffer(b))
    if err != nil {
        fmt.Println("error raise; %v", err)
    }
    body, _ := ioutil.ReadAll(res.Body)
    fmt.Println(string(body))
}

func Invoke(){
    t := &transmit{
        "2.0",
        "invoke",
        map[string]interface{}{"chaincodeID" : map[string]interface{}{"name":"123"},"ctorMsg" : map[string]interface{}{"function" : "invoke","args" : []string{"a", "b", "1"}},"secureContext":"jim"},3}
    b, err := json.Marshal(t)
    if err != nil {
        fmt.Println("error raised: %v", err)
    }
    res, err := http.Post(addr+"/chaincode", "application/json", bytes.NewBuffer(b))
    if err != nil {
        fmt.Println("error raised: %v", err)
    }
    body, _ := ioutil.ReadAll(res.Body)
    fmt.Println(string(body))
}

func Query(){
    t := &transmit{
        "2.0",
        "query",
        map[string]interface{}{"chaincodeId" : map[string]interface{}{"name":"123"},"ctorMsg":map[string]interface{}{"function" : "query", "args" : []string{"a"}},"secureContext":"jim"},3}
    b, err := json.Marshal(t)
    if err != nil {
        fmt.Println("error raised: %v", err)
    }
    resp, err := http.Post(addr+"/chaincode", "application/json", bytes.NewBuffer(b))
    if err != nil {
        fmt.Println("error raised: %v", err)
    }
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}

func main() {
    arg := os.Args[1]
    switch arg {
    case "register": Register()
    case "deploy"  : Deploy()
    case "invoke"  : Invoke()
    case "query"   : Query()
    default : fmt.Println("use deploy/register")
    }
}

