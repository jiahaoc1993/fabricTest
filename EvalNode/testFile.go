package main

import (
	"fmt"
	"os"
)

func main() {

	f, err := os.OpenFile("./test.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString("Fuck you!!\n")
	if err != nil {
		panic(err)
	}

	fmt.Println("write to test.txt successs!")
	fmt.Println(500/50)
	fmt.Println(2/50)
	fmt.Println(3/50)
	fmt.Println(1/100)
}
