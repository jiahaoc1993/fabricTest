package main

import (
	"fmt"
	"tool/loadKey"
)

func main() {
	_, _, err := loadKey.LoadEnrollment()
	if err != nil {
		fmt.Println(err)
	}
}
