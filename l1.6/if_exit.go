package main

import (
	"fmt"
	"time"
)

// Выход по условию.
// просто считаем в горутине до 10 и когда досчитали выходим
// В мэйне считаем подольше
func main() {
	go func() {
		i := 0
		for {
			fmt.Println("func: ", i)
			i++
			time.Sleep(300 * time.Millisecond)
			if i == 10 {
				fmt.Println("func: Counted to %d, exit", i)
				return

			}
		}

	}()
	i := 0
	for {
		fmt.Println("main: ", i)
		i++
		time.Sleep(200 * time.Millisecond)
		if i == 20 {
			fmt.Println("main: Counted to %d, exit", i)
			return

		}
	}

}
