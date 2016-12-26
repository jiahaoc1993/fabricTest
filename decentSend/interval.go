package main

import (
	"./method"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	//c := make(chan int)
	count := 0
	for time.Since(start).Seconds() <= 1 {
		go method.Query()
		count++
	}
	fmt.Println(count)
	time.Sleep(10)
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
