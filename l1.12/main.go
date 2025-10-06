package main

import "fmt"

func main() {

	animals := []string{"cat", "cat", "dog", "cat", "bird", "fish", "dog"}
	set := make(map[string]bool) //Буду использовать мэп чтобы быстро добавлять значения по ключу названия животного
	//Это поможет избежать дубликатов

	for _, anim := range animals {
		set[anim] = true
	}

	for anim := range set {
		fmt.Print(anim, ", ")
	}
}
