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


func main() {
    arg := os.Args[1]
    addr := url + ":" + port
    resp, err := http.Get(addr+arg)
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


