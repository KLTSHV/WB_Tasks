package main

import (
	"fmt"
	"time"
)

func sleep(dur time.Duration) {
	ch := make(chan struct{})
	go func() {
		time.AfterFunc(dur, func() {
			close(ch)
		})
	}()
	<-ch
}

func main() {
	fmt.Println("Sleeping for 5 seconds")
	sleep(5 * time.Second)
	fmt.Println("Awake!")
}
