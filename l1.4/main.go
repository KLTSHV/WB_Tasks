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

type Task struct {
	ID int
	A  int
	B  int
}

// считает сумму двух рандомных чисел, которые воркер передаст,
// seed вставляется в main
// мне кажется так правильней чем генерировать seed в продюсере)
func doWork(ctx context.Context, t Task) {
	// имитация работы
	select {
	case <-ctx.Done():
		return
	case <-time.After(time.Duration(50+rand.Intn(100)) * time.Millisecond):
	}
	result := t.A + t.B
	fmt.Printf("task #%d: %d + %d = %d\n", t.ID, t.A, t.B, result)
}

func worker(id int, ctx context.Context, jobs <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			// мягкая остановка по Ctrl+C
			return
		case t, ok := <-jobs:
			if !ok {
				// канал закрыт продьюсером — задач больше нет
				return
			}
			doWork(ctx, t)
		}
	}
}

func producer(ctx context.Context, jobs chan<- Task, r *rand.Rand) {
	defer close(jobs)
	id := 0
	for {
		task := Task{
			ID: id,
			A:  r.Intn(10) + 5,
			B:  r.Intn(10) + 5,
		}
		id++
		select {
		case <-ctx.Done():
			// отмена по Ctrl+C — выходим, defer закрое jobs
			return
		case jobs <- task:
			// задача поставлена в очередь
		}
	}
}

func main() {
	// Контекст, привязанный к SIGINT (Ctrl+C)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// генератор случайных чисел
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	const workerCount = 4
	jobs := make(chan Task, 3*workerCount)

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker(i, ctx, jobs, &wg)
	}

	go producer(ctx, jobs, r)

	// Ждём Ctrl+C
	<-ctx.Done()

	wg.Wait()
	fmt.Println("Завершение выполнено!")
}
