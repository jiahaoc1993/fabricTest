package main

import (
	//"bytes"
	//	"encoding/binary"
	"fmt"
	//"net"
)

type testIN struct {
	B []byte
}

const (
	hexDigit = "0123456789abcdef"
)

func hexstring(b []byte) string { //all right , this is the secret, seprate one byte(8 bits) into two grop of 4 bits, each group represent a hex digit
	s := make([]byte, 2*len(b))
	for i, tn := range b {
		s[2*i], s[2*i+1] = hexDigit[tn>>4], hexDigit[tn&0xf]
	}
	return string(s)
}

func main() {

	//b := []byte{255, 255, 255, 0}
	b := make([]byte, 4)
	b[0] = 255
	b[1] = 255
	b[2] = 255
	b[3] = 0
	fmt.Println(b, hexstring(b))
}

//I finally found out why address.net return bytes ffffff00 rather than 255.255.255.255 , the type of address.net is net.IPMask
