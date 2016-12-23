package main

import (
    "./method"
    "fmt"
    //"io"
    "time"
)

func test(c chan int){
    method.Register()
    method.Deploy()
    method.Invoke()
    method.Query()
    c <- 1
}


func main(){
    numOfReq := 100
    c := make(chan int)
    start := time.Now().Unix()
    for i := 0 ; i < numOfReq ; i++ {
       go test(c)
    }
    for s :=0; s < numOfReq ; {
        s += <-c
    }
    end := time.Now().Unix()
    spent := end - start
    fmt.Printf("All Done! Spent %d seconds\n", spent)
}



