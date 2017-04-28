package main

import(
	"fmt"
	"tool/loadKey"
)

func main(){
	enrollmentCert, privKey, err := loadKey.LoadFakeEnrollment()
	if err != nil {
		fmt.Println("Failed loading enrollment metieral")
		return 
	}

	fmt.Println("Public Key: ", enrollmentCert.Raw)
	fmt.Println("Private Key:", privKey)

}
