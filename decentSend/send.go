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
    port string = "7050"
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

func Post(){
    t := &transmit{
        "2.0",
        "deploy",
        map[string]interface{}{"abc" : "123","cba" : 123},
        123,
    }
    b, err := json.Marshal(t)
    if err != nil {
        fmt.Printf("error raise: %v", err)
    }
    res, err := http.Post("http://172.22.28.130:7050/registrar", "application/json", bytes.NewBuffer(b))
    if err != nil {
        fmt.Println("error raise; %v", err)
    }
    fmt.Println(res)
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
    Post()

}

