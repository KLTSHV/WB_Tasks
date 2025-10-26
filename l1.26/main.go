package main

import "strings"

func isUnique(s string) bool {
	charMap := make(map[string]bool) //Мэп для хранения уже встреченных символов, храним в нижнем регистре
	for _, c := range s {
		char := strings.ToLower(string(c)) //Преобразуем символ в нижний регистр
		if _, ok := charMap[char]; ok {
			return false
		} else {
			charMap[char] = true
		}

	}
	return true
}
func main() {
	s1 := "HElLoWOrlD"
	s2 := "ABCD"
	s3 := "aAaa"
	s4 := ""
	println(isUnique(s1)) // false
	println(isUnique(s2)) // true
	println(isUnique(s3)) // false
	println(isUnique(s4)) // true

}
