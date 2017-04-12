package main

import(
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	time.Sleep(1 * time.Second)
	spent := time.Now().Sub(now)
	fmt.Println(spent.Seconds())
}
