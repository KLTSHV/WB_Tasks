package main

import (
	"fmt"
	"sync"
)

func worker(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for x := range ch {
		fmt.Printf("worker %v: %v\n", id, x)
	}

}
func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		fmt.Println("ошибка сканирования")
		return
	}
	if n <= 0 {
		fmt.Println("n должно быть >= 1")
		return
	}
	var wg sync.WaitGroup

	ch := make(chan int, 3*n)

	for j := 1; j <= n; j++ {
		wg.Add(1)
		go worker(j, ch, &wg)
	}
	for i := 0; i < 100; i++ {
		ch <- i

	}
	close(ch)
	wg.Wait()

}
