package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

var mtx sync.Mutex             // мьютекс для синхронизации доступа к карте
var m = make(map[int]int)      // общая карта, с которой работают горутины
var done = make(chan struct{}) // канал для завершения горутин
var wg sync.WaitGroup          // чтобы дождаться завершения всех горутин

// функция записи значения по ключу в map с блокировкой
func Set(sm map[int]int, key int, value int) {
	mtx.Lock()
	sm[key] = value
	mtx.Unlock()
}

// функция чтения значения по ключу в map с блокировкой
func Read(rm map[int]int, key int) int {
	mtx.Lock()
	value := rm[key]
	mtx.Unlock()
	return value
}

// воркер, который ищет максимальное значение в карте и уменьшает его
func workerDecreaser() {
	for {
		select {
		case <-done:
			wg.Done()
			return
		default:
			max := math.MinInt64
			var max_id int
			for i := 0; i < 10; i++ {
				value := Read(m, i)
				if value >= max {
					max_id = i
					max = value
				}
			}
			Set(m, max_id, max-20) // уменьшаем максимальное значение
			time.Sleep(300 * time.Millisecond)
		}
	}
}

// воркер, который ограничивает значения в карте в пределах (-50, 50)
func workerLimiter() {
	for {
		select {
		case <-done:
			wg.Done()
			return
		default:
			for i := 0; i < 10; i++ {
				value := Read(m, i)
				if value >= 50 || value <= -50 {
					Set(m, i, 0) // сбрасываем значение
				}
			}
			time.Sleep(700 * time.Millisecond)
		}
	}
}

// воркер, который случайно увеличивает одно из значений
func workerIncreaser() {
	for {
		select {
		case <-done:
			wg.Done()
			return
		default:
			var r = rand.New(rand.NewSource(time.Now().UnixNano()))
			key := r.Intn(10) // случайный индекс
			Set(m, key, Read(m, key)+1)
			time.Sleep(15 * time.Millisecond)
		}
	}
}

func main() {
	wg.Add(3) // три горутины

	// инициализация карты
	for i := 0; i < 10; i++ {
		m[i] = 0
	}

	// запускаем горутины
	go workerIncreaser()
	go workerDecreaser()
	go workerLimiter()

	// выводим состояние карты 50 раз
	for i := 0; i < 50; i++ {
		fmt.Print("Map: ")
		for j := 0; j < 10; j++ {
			fmt.Print(Read(m, j), " ")
		}
		fmt.Print("\n")
		time.Sleep(500 * time.Millisecond)
	}

	// останавливаем горутины
	close(done)
	wg.Wait()
}
