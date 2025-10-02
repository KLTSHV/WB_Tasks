package main

func contains(arr []int, v int) bool {
	for _, val := range arr {
		if val == v {
			return true
		}
	}
	return false
}
func intersection(arr1 []int, arr2 []int) []int {
	result := []int{}
	for _, v := range arr1 {
		if contains(result, v) {
			continue
		}
		if contains(arr2, v) {
			result = append(result, v)
		}

	}
	return result
}
func main() {
	a := []int{1, 3, 2, 5, 4, 3}
	b := []int{3, 4, 9, 10, 4, 5}
	r := intersection(a, b)
	print("{")
	for i, v := range r {
		if i != 0 {
			print(", ")
		}
		print(v)
	}
	println("}")
}
