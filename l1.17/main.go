package main

import "fmt"

func myFind(s []int, value int) int {
	right := len(s) - 1
	left := 0
	med := len(s) / 2
	for {
		if value == s[med] {
			return med
		} else if right == left || left-1 == right || left+1 == right { //если между левым и правым элементом нет ничего или они равны
			if value == s[left] && left >= 0 && left <= len(s) {
				return left
			} else if value == s[right] && right < len(s) && right >= 0 {
				return right
			} else {
				return -1
			}
			// если искомое меньше (больше) середины, то новая правая (левая) граница будет середина
		} else if value < s[med] {
			right = med
		} else if value > s[med] {
			left = med
		}
		med = (right-left)/2 + left
	}
}
func main() {
	A := []int{0, 1, 2, 3, 5, 6, 7, 8, 9} //нет 4
	fmt.Println("pos", "element")
	for i := -1; i <= len(A)+1; i++ { //итерируюсь от -1 до 10 (-1, 4, 10 нет в массиве)
		s := myFind(A, i)
		fmt.Println(s, i)
	}
	fmt.Println("------")
	//тут попробую найди в срезе из одного элемента
	s := []int{7}
	p := myFind(s, 1)
	fmt.Println(p, 1)

}
