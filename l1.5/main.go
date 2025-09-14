package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	done := make(chan struct{})
	N := 5
	go func() {
		i := 0
		for {
			select {
			case <-done:
				close(ch)
				return
			case ch <- i:
				i++
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()
	timeout := time.After(time.Duration(N) * time.Second)

	for {
		select {
		case val, ok := <-ch:
			if !ok {
				fmt.Println("Канал закрыт, завершение")
				return
			}
			fmt.Println("Прочитано из канала:", val)
		case <-timeout:
			fmt.Println("Время истекло")
			close(done)
			return
		}
	}
}
