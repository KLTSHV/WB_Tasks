package main

import (
	"fmt"
)

func main() {
	var n int64
	fmt.Scan(&n)
	for {
		var i int
		fmt.Scan(&i)
		if i > 1 && i < 64 {
			for {
				var bit int
				fmt.Scan(&bit)
				// установка i-го бита числа n в значение bit
				if bit == 1 {
					n |= 1 << (i - 1)
					fmt.Println(n)
					return
				} else if bit == 0 {
					n &^= 1 << (i - 1)
					fmt.Println(n)
					return
				} else {
					fmt.Println("bit must be 0 or 1")
				}
			}
		} else {
			fmt.Println("i must be in range 1..64")
		}

	}
}
