package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Выход по контексту.
// считаем в горутине и в это время в мэйне спим 3 секунды
func main() {
	ctx, cancel := context.WithCancel(context.Background()) // создаем контекст с возможностью отмены
	var wg sync.WaitGroup                                   // для ожидания завершения горутины
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		i := 0
		for {
			select {
			case <-ctx.Done():
				fmt.Println("func exiting")
				return
			default:
				fmt.Println("func:", i)
				i++
				time.Sleep(300 * time.Millisecond)
			}
		}
	}(ctx)
	time.Sleep(3 * time.Second)
	cancel()  // вызываем отмену контекста, чтобы уведомить горутину о выходе
	wg.Wait() // ждем завершения горутины

}
