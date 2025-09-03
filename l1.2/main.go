package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	a := [5]int{2, 4, 6, 8, 10}
	for _, n := range a {
		wg.Add(1)
		go func(x int) {
			b := x * x
			fmt.Println(b)
			defer wg.Done()
		}(n)
	}
	wg.Wait()
}
