package main

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
)

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	//	fmt.Println(devices)
	fmt.Println("Devices found:")
	for _, device := range devices {
		fmt.Println("\nName:", device.Name)
		fmt.Println("Description", device.Description)
		fmt.Println("Device address", device.Description)
		for _, address := range device.Addresses {
			fmt.Println("-- IP ADDRESS: ", address.IP)
			fmt.Printf("-- Netmask: %d.%d.%d.%d\n", address.Netmask[0], address.Netmask[1], address.Netmask[2], address.Netmask[3])
		}
	}

}
