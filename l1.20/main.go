package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	r := []rune(strings.TrimSpace(s))
	// полностью реверсим строку
	for i := 0; i < len(r)/2; i++ {
		r[i], r[len(r)-i-1] = r[len(r)-i-1], r[i]
	}
	// реверсим каждое слово
	start := 0
	for i := 0; i < len(r); i++ {
		if r[i] == ' ' || i == len(r)-1 {
			if i == len(r)-1 {
				i++ // чтобы включить последний символ в реверс
			}
			for j := start; j < ((i-start)/2)+start; j++ {
				r[j], r[i-(j-start)-1] = r[i-(j-start)-1], r[j]
			}
			start = i + 1 // следующий старт после пробела
		}
	}
	fmt.Println(string(r))
}
