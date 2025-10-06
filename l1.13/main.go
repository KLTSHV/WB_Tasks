package main

import "fmt"

func main() {
	//Сложение вычитание
	a := 7
	b := 10
	a = a + b //17
	b = a - b //7
	a = a - b //10
	fmt.Println("a = ", a)
	fmt.Println("b = ", b)
	//XOR
	a = a ^ b
	b = a ^ b
	a = a ^ b
	fmt.Println("a = ", a)
	fmt.Println("b = ", b)

}
