package grep

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Options struct {
	After       int
	Before      int
	IgnoreCase  bool
	Invert      bool
	Fixed       bool
	CountOnly   bool
	LineNumbers bool
}

func match(line string, pattern string, opt Options) bool {
	str := line
	pat := pattern

	if opt.IgnoreCase {
		str = strings.ToLower(str)
		pat = strings.ToLower(pat)
	}

	var ok bool
	if opt.Fixed {
		ok = strings.Contains(str, pat)
	} else {
		re, err := regexp.Compile(pat)
		if err != nil {
			log.Fatalf("invalid regexp pattern: %v", err)
		}
		ok = re.MatchString(str)
	}

	if opt.Invert {
		return !ok
	}
	return ok
}

// обрабатывает слайс строк, применяя контекст и флаги
func ProcessLines(lines []string, pattern string, opt Options) {
	matched := make([]bool, len(lines))

	for i, line := range lines {
		if match(line, pattern, opt) {
			matched[i] = true
		}
	}

	// если флаг -c, считаем совпадения и выходим
	if opt.CountOnly {
		count := 0
		for _, m := range matched {
			if m {
				count++
			}
		}
		fmt.Println(count)
		return
	}

	// формируем map строк для вывода с контекстом
	toPrint := make(map[int]struct{})
	for i, m := range matched {
		if m {
			from := max(0, i-opt.Before)
			to := min(len(lines)-1, i+opt.After)
			for j := from; j <= to; j++ {
				toPrint[j] = struct{}{}
			}
		}
	}

	// выводим строки в порядке файла
	for i := 0; i < len(lines); i++ {
		if _, ok := toPrint[i]; ok {
			if opt.LineNumbers {
				fmt.Printf("%d:%s\n", i+1, lines[i])
			} else {
				fmt.Println(lines[i])
			}
		}
	}
}

// ProcessFile читает файл и вызывает processLines
func ProcessFile(filename string, pattern string, opt Options) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("cannot open file %s: %v", filename, err)
		return
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Printf("error reading file %s: %v", filename, err)
		return
	}

	ProcessLines(lines, pattern, opt)
}

// processStdin читает stdin и вызывает processLines
func ProcessStdin(pattern string, opt Options) {
	lines := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Printf("error reading stdin: %v", err)
		return
	}

	ProcessLines(lines, pattern, opt)
}

// вспомогательные функции
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
