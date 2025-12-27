package main

import (
	"fmt"
	"time"
)

// or канал, который закрывается, когда закрывается любой из переданных каналов
// Реализация использует рекурсию для объединения нескольких каналов
func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	default:
		r := make(chan interface{})
		go func() {
			defer close(r)
			select {
			case <-channels[0]:
			case <-or(channels[1:]...):
			}
		}()
		return r
	}
}
func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
