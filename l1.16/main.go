package main

import "fmt"

func quickSort(s []int) []int {
	if len(s) < 2 {
		return s
	}
	pivot := s[len(s)/2] //случайно выбираем опорный элемент от которого будет сравнивать другие элементы
	var less, equal, greater []int

	for _, x := range s {
		if x < pivot {
			less = append(less, x)
		} else if x == pivot {
			equal = append(equal, x)
		} else {
			greater = append(greater, x)
		}
	}
	//рекурсивно возвращаем все что меньше опорного элемента, опорный элемент и ему равные, и все что больше опорного элемента для каждого куска.
	return append(append(quickSort(less), equal...), quickSort(greater)...)
}
func main() {
	a1 := []int{1, 0, 3, 0, 2, 1, 10, 1, 2, -1, -9, 80, -100}
	a2 := []int{-1}
	a3 := []int{0, 0, 0, 0}

	s1 := quickSort(a1)
	s2 := quickSort(a2)
	s3 := quickSort(a3)

	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s3)

}
