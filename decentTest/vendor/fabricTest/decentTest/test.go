package main

import(
	"fmt"
	"time"
)

func main() {
	fmt.Println(100/1000)
	t := time.NewTicker(time.Second)
	for i:=0 ; i< 10 ; i++{
		select {
		case <-t.C :
			fmt.Println("Hello World")
		default :
			fmt.Println("Not time yet")
		}
		time.Sleep(time.Millisecond * 500)
	}
	t.Stop()

}
