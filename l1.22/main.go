package main

import (
	"fmt"
	"math/big"
)

func main() {
	a := big.NewInt(123000123011)
	b := big.NewInt(56700657000799)
	sum := new(big.Int).Add(a, b)
	fmt.Println("sum: ", sum)
	subs := new(big.Int).Sub(a, b)
	fmt.Println("substraction: ", subs)
	mult := new(big.Int).Mul(a, b)
	fmt.Println("multiplication: ", mult)

	//Для деления использую big.Float
	fa := new(big.Float).SetInt(a)
	fb := new(big.Float).SetInt(b)
	div := new(big.Float).Quo(fb, fa)
	fmt.Println("division: ", div)

}
