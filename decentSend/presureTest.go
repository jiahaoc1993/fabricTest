package main

import (
	"./method"
	"fmt"
	//"io"
	"os"
	"strconv"
	"time"
)

func mutilFunc(f string, n int) {
	if f == "invoke" {
		c := make(chan int)
		start := time.Now().Unix()
		for i := 0; i < n; i++ {
			method.Invoke()
			c <- 1
		}
		for s := 0; s < n; {
			s += <-c
		}
		end := time.Now().Unix()
		spent := end - start
		fmt.Println("All Done! Spent %d Seconds\n", spent)
	}
}

func main() {
	switch os.Args[1] {
	case "register":
		method.Register()
	case "deploy":
		method.Deploy()
	case "invoke":
		i, _ := strconv.Atoi(os.Args[2])
		mutilFunc("invoke", i)
	case "query":
		method.Query()
	}

}
