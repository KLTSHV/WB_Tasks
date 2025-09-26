package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

// счетчик для того чтобы отслеживать прогресс генерации символов в массив продюсером (value_produce)
// и отслеживать прогресс считывания символов воркером (value_worker)
type Counter struct {
	value_produce int
	value_worker  int
	mux           sync.Mutex
}

func producer(ctx context.Context, arr *[50]int, c *Counter, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 50; i++ {

		select {
		case <-ctx.Done():
			fmt.Println("Producer was interrupted!")
			return
		default:

			c.mux.Lock()
			c.value_produce++
			arr[i] = rand.Intn(100)
			c.mux.Unlock()
			time.Sleep(100 * time.Millisecond)

		}

	}
	fmt.Println("Producer Done!")
}

// достает числа из массива в порядке их появления (arr[0] -> arr[199])
// пока продюсер не сгенерирует следующий символ, будет стопориться на одном и том же индексе и ждать
func worker1(ctx context.Context, chx chan<- int, arr *[50]int, c *Counter, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(chx)

	for i := 0; i < 50; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Worker1 was interrupted!")
			return
		default:
			c.mux.Lock()
			if c.value_produce > c.value_worker {
				chx <- arr[i]
				c.value_worker++
			} else {
				i--
			}
			c.mux.Unlock()

		}
	}
	fmt.Println("Worker1 Done!")
}

func worker2(ctx context.Context, chx <-chan int, chx2 chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(chx2)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker2 Done!")
			return
		case n, ok := <-chx:
			if !ok {
				fmt.Println("Worker2 Done! No numbers left")
				return
			}
			select {
			case <-ctx.Done():
				return
			case chx2 <- n * 2:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
func printer(ctx context.Context, chx2 <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Printer Done!")
			return
		case n, ok := <-chx2:
			if !ok {
				return
			}
			fmt.Println(n)
		}
	}
}
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	var wg sync.WaitGroup
	c := Counter{
		value_produce: 0,
		value_worker:  0,
		mux:           sync.Mutex{},
	}

	chx := make(chan int)
	chx2 := make(chan int)
	var arr [50]int
	wg.Add(4)
	go producer(ctx, &arr, &c, &wg)
	go worker1(ctx, chx, &arr, &c, &wg)
	go worker2(ctx, chx, chx2, &wg)
	go printer(ctx, chx2, &wg)
	wg.Wait()
	fmt.Println(c.value_produce, c.value_worker)

}
