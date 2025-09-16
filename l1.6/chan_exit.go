package main

import (
	"fmt"
	"sync"
	"time"
)

// Выход по каналу.
// просто считаем в горутине и в это время в мэйне спим 3 секунды
func main() {
	done := make(chan struct{})
	var wg sync.WaitGroup // для ожидания завершения горутины
	wg.Add(1)             // добавляем одну горутину в группу ожидания

	go func() {
		defer wg.Done()
		i := 0
		for {
			select {
			case <-done: // если канал закрыт, выходим
				fmt.Println("Notification received")
				return
			default:
				fmt.Println("func: ", i)
				i++
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	time.Sleep(3 * time.Second)
	close(done) // закрываем канал, чтобы уведомить горутину о выходе
	wg.Wait()   // ждем завершения горутины
}
