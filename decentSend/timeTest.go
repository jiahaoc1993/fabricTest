package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	ticker := time.NewTicker(1)
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		fmt.Printf("Sleeping 1 second start from %v ,", (<-ticker.C).Unix())
		fmt.Printf("Time past since programing started: %0.0f second(s)\n", time.Since(start).Seconds())
	}
	ticker.Stop()
}
