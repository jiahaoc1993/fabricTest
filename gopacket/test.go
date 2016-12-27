package main

import (
	//"bytes"
	//	"encoding/binary"
	"fmt"
	"net"
)

type testIN struct {
	B []byte
}

func main() {
	var mask net.IPMask
	//b := []byte{255, 255, 255, 0}
	b := make([]byte, 4)
	b[0] = 255
	b[1] = 255
	b[2] = 255
	b[3] = 0
	mask = b
	fmt.Println(mask[0])
}

//I finally found out why address.net return bytes ffffff00 rather than 255.255.255.255 , the type of address.net is net.IPMask
