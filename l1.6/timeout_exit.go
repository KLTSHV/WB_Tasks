package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Горутина через дэдлайн
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("timeout завершает")
				wg.Done()
				return
			default:
				fmt.Println("Горутина работает")
				time.Sleep(time.Millisecond * 200)
			}
		}
	}(ctx)

	wg.Wait()

}
