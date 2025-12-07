package compare

import (
	"sort"
	"strconv"
	"strings"
)

// SortLines сортирует строки с учетом переданных флагов
func SortLines(
	lines []string,
	column int,
	number, humannumeric, month, blanks, reverse, unique, checksort bool,
) []string {
	// Убираем дубликаты
	if unique {
		seen := map[string]int{}
		var uniqueLines []string
		for _, line := range lines {
			key := Getcolumn(line, column)
			if seen[key] == 0 {
				seen[key]++
				uniqueLines = append(uniqueLines, line)
			} else {
				seen[key]++
				if checksort {
					panic("Duplicate key found: " + key)
				}
			}
		}
		lines = uniqueLines
	}

	// Проверка сортировки
	if checksort {
		if !sort.SliceIsSorted(lines, func(i, j int) bool {
			return Compare(lines[i], lines[j], column, number, humannumeric, month, blanks, reverse)
		}) {
			panic("Данные не отсортированы")
		}
		return lines
	}

	// Сортировка строк
	sort.Slice(lines, func(i, j int) bool {
		return Compare(lines[i], lines[j], column, number, humannumeric, month, blanks, reverse)
	})

	return lines
}

// Compare сравнивает две строки с учетом всех флагов
func Compare(a, b string, column int, number, humannumeric, month, blanks, reverse bool) bool {
	av := Getcolumn(a, column)
	bv := Getcolumn(b, column)

	if blanks {
		trimAv := strings.TrimLeft(av, " \t")
		trimBv := strings.TrimLeft(bv, " \t")

		// полностью пустые строки считаем меньше любых других
		if trimAv == "" && trimBv != "" {
			return reverse // если reverse=true, пустые идут в конец
		}
		if trimBv == "" && trimAv != "" {
			return !reverse
		}

		if trimAv != "" {
			av = trimAv
		}
		if trimBv != "" {
			bv = trimBv
		}
	}

	if month {
		av = MonthValue(av)
		bv = MonthValue(bv)
		av += Getcolumn(a, column+1)
		bv += Getcolumn(b, column+1)

	}
	if humannumeric {
		an := HumannumericValue(av)
		bn := HumannumericValue(bv)
		if reverse {
			return an > bn
		}
		return an < bn
	}
	if number {
		an, err1 := strconv.ParseFloat(av, 64)
		bn, err2 := strconv.ParseFloat(bv, 64)
		if err1 == nil && err2 == nil {
			if reverse {
				return an > bn
			}
			return an < bn
		}
		if err1 != nil {
			return false
		}
		if err2 != nil {
			return true
		}
	}
	if reverse {
		return av > bv
	}
	return av < bv
}

// Getcolumn возвращает столбец строки по индекс
func Getcolumn(s string, c int) string {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '\t'
	})
	if c-1 >= 0 && c-1 < len(fields) {
		return fields[c-1]
	}
	// Если колонка нет то возвращаем исходную строку
	return s
}

// MonthValue преобразует название месяца в числовое значение в формате стороки
func MonthValue(s string) string {
	months := map[string]string{
		"Jan": "01", "Feb": "02", "Mar": "03", "Apr": "04",
		"May": "05", "Jun": "06", "Jul": "07", "Aug": "08",
		"Sep": "09", "Oct": "10", "Nov": "11", "Dec": "12",
	}
	val, ok := months[s]
	if !ok {
		return "00"
	}
	return val
}

// HumannumericValue преобразует значения с K/M/G/T/P/E в число
func HumannumericValue(s string) float64 {
	if len(s) == 0 {
		return 0
	}
	multipliers := map[string]float64{
		"K": 1e3, "M": 1e6, "G": 1e9, "T": 1e12,
		"P": 1e15, "E": 1e18,
	}
	lastChar := s[len(s)-1:]
	multiplier, ok := multipliers[lastChar]
	if ok {
		num, err := strconv.ParseFloat(s[:len(s)-1], 64)
		if err != nil {
			return 0
		}
		return num * multiplier
	}
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return num
}

// ExpandShortFlags расширяет короткие флаги типа -nr в -n -r
func ExpandShortFlags(args []string) []string {
	var expanded []string
	for _, arg := range args {
		if arg[0] == '-' && len(arg) > 2 {
			for _, ch := range arg[1:] {
				expanded = append(expanded, "-"+string(ch))
			}
		} else {
			expanded = append(expanded, arg)
		}
	}
	return expanded
}
