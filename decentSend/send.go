package main

import (
    "fmt"
    "net/http"
    "os"
    "io/ioutil"
)

const (
    url string = "http://172.22.28.118"
    port string = "7050"
)

func GetdMessage() {
    args := os.Args[1:]
    if len(args) == 1 {
        
    }
}


func main() {
    SendMessage()
    addr := url + ":" + port
    resp, err := http.Get(addr+"/chain/blocks/0")
    if err != nil {
        fmt.Printf("Error: %v", err)
        os.Exit(1)
    }
    b, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
        fmt.Printf("Error: %v", err)
        os.Exit(1)
    }
    fmt.Printf("Message : %s", b)
}


