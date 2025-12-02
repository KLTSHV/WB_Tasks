package main

import (
	"fmt"

	"l29/internal/unpackstring"
)

func main() {
	var s string
	fmt.Scanln(&s)
	res, err := unpackstring.Unpack(s)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(res)

}
