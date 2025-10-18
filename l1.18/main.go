package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup

// Создам два счетчика: с мьютексом и атомиком,
// и напишу для них функции увеличения значение и получение этого значения
type CounterMutex struct {
	mu    sync.Mutex
	count int
}

func (c *CounterMutex) Inc() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}
func (c *CounterMutex) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

type CounterAtomic struct {
	ai atomic.Int64
}

func (c *CounterAtomic) Inc() {
	c.ai.Add(1)
}
func (c *CounterAtomic) Value() int {
	return int(c.ai.Load())
}
func main() {
	var count_a CounterAtomic
	var count_m CounterMutex
	//Атомик счетчик
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(count_a *CounterAtomic) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				count_a.Inc()
			}
		}(&count_a)
	}
	//Мьютекс счетчик
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(count_m *CounterMutex) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				count_m.Inc()
			}
		}(&count_m)
	}
	wg.Wait()
	//Оба должны показать 100
	fmt.Println("CountAtomic: ", count_a.Value())
	fmt.Println("CountMutex: ", count_m.Value())
}
