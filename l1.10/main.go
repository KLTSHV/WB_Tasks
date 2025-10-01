package main

import "fmt"

func sortTempatures(arr []float32) map[int][]float32 {
	result := make(map[int][]float32)
	for i := range arr {
		key := int(arr[i]) / 10 * 10
		if arr[i] >= 0 || arr[i] == float32(key) {
			result[key] = append(result[key], arr[i])
		} else if arr[i] != float32(key) {
			result[key-10] = append(result[key-10], arr[i])
		}
	}
	return result
}

/*
В примере задания я не понимаю,
что включает в себя диапазон 0: {0, 1, ... 9, 10} или {-10, -9, ..., -1, 0}
поэтому пусть каждый ключ мапы обозначачает самый низкий порог диапазона.
То есть -20: {-20, -19.1, ..., -10.1}
-10: {-10, -9.1, ..., -0.1}
*/
func main() {
	a := []float32{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	r := sortTempatures(a)
	fmt.Println(r)
}
