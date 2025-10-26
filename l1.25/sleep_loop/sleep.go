package main

import (
	"fmt"
	"time"
)

func sleep(dur time.Duration) {
	start := time.Now()
	for time.Since(start) < dur {

	}
}

func main() {
	fmt.Println("Sleeping for 5 seconds")
	sleep(5 * time.Second)
	fmt.Println("Awake!")
}
