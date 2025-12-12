package main

import (
	"fmt"
	"sort"
	"strings"
)

// сортирует строку в алфавитном порядке c помощью слайса рун
func sortStr(s string) string {
	runes := []rune(s)
	sort.Slice(runes, func(a, b int) bool {
		return runes[a] < runes[b]
	})
	return string(runes)
}

// Как работает алгоритм
// Создаем мапу, где ключ - отсортированная строка, значение - слайс слов, которые являются анаграммами
// Проходим по каждому слову, приводим его к нижнему регистру, сортируем буквы и добавляем в мапу
// После этого проходим по мапе и отбираем только те слайсы, где больше одного слова
func anagrams(words []string) map[string][]string {
	anagramMap := make(map[string][]string)
	for _, w := range words {
		w = strings.ToLower(w)
		sorted := sortStr(w)
		anagramMap[sorted] = append(anagramMap[sorted], w)
	}
	for ind := range anagramMap {
		sort.Strings(anagramMap[ind])
	}
	result := make(map[string][]string)
	for _, s := range anagramMap {
		if len(s) < 2 {
			continue
		}
		result[s[0]] = s
	}
	return result
}

func main() {
	ex1 := []string{"стол", "тОлс", "лист", "Кот", "олтс", "ток", "ситл"}
	r1 := anagrams(ex1)
	fmt.Println(r1)
	ex2 := []string{"пятак", "пятка", "тЯпка", "листок", "сЛиТок", "столик", "кот"} //кот - 1, не должен выводиться
	r2 := anagrams(ex2)
	fmt.Println(r2)

}
