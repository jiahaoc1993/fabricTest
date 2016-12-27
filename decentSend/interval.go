package main

import (
	"./method"
	"fmt"
	"time"
)

func test() {
	start := time.Now()
	//c := make(chan int)
	count := 0
	for time.Since(start).Seconds() <= 2 {
		method.Invoke()
		count++
	}
	fmt.Println(count)
	/*
		for i := 0; i < 10; i++ {
			go func() {
				method.Query()
				c <- 1
			}()
		}
		for {
			count += <-c
			if count == 10 {
				fmt.Println(time.Since(start).Seconds())
				break
			}
		}
	*/
}

func main() {
	for i := 0; i < 50; i++ {
		go test()
	}
	time.Sleep(10 * time.Second)
}
