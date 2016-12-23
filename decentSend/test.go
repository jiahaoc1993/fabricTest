package main

import (
    "fmt"
    "time"
)

func main(){
    fmt.Println(time.Now().Unix())
    time.Sleep(10* time.Second)
    fmt.Println(time.Now().Unix())
}
