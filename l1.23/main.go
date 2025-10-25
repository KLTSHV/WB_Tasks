package main

import "fmt"

func main() {
	s := []int{1, 2, 3, 4, 5, 6, 7}
	i := 4
	fmt.Printf("Before deleting %dth element: %d\n", i, s)

	copy(s[i:], s[i+1:])
	s = s[:len(s)-1]

	fmt.Printf("After deleting %dth element: %d\n", i, s)
}
