package main

import "fmt"

func detectType(i interface{}) {
	switch t := i.(type) {
	case int:
		fmt.Println("int ", t)
	case string:
		fmt.Println("string ", t)
	case bool:
		fmt.Println("bool ", t)
	case chan int:
		fmt.Println("chan int ", t)
	case chan string:
		fmt.Println("chan string ", t)
	case chan bool:
		fmt.Println("chan bool ", t)
	default:
		fmt.Println("unknown type ", t)
	}
}
func main() {
	detectType(10)
	detectType("WB")
	detectType(true)
	detectType(make(chan int))
	detectType(make(chan string))
	detectType(make(chan bool))
	detectType(9.5)
}
