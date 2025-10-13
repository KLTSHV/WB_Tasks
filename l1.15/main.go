package main

var justString string

// Функция createHugeString создаёт огромную строку заданного размера.
func createHugeString(size int) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = 'w'
	}
	return string(b)
}

// простой v[:100] НЕ РАБОТАЕТ т.к. строка остаётся в памяти целиком
func someFunc() {
	v := createHugeString(1 << 10)
	justString = string([]byte(v)[:100]) // Срезаем строку до первых 100 символов КОПИРУЯ
}

func main() {
	someFunc()
}
