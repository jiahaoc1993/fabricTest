package main

import (
	//"bytes"
	//	"encoding/binary"
	"fmt"
)

func main() {
	b := []byte{255, 255, 255, 0}
	for _, h := range b {
		fmt.Printf("%x\n", h)
	}
}
