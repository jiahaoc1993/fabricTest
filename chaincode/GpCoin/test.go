package main

import(
	"fmt"
	"flag"
)

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println(args[0], args[1])

}
