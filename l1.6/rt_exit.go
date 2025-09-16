package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer fmt.Println("defer сработает перед завершением")
		fmt.Println("Горутина запущена")
		runtime.Goexit()
		fmt.Println("Это не будет напечатано")
	}()
	fmt.Scanln() //чтобы программа не завершилась сразу

}
