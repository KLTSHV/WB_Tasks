package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"l2_10/internal/compare"
)

func main() {
	var input *os.File

	// Расширяем короткие флаги
	os.Args = compare.ExpandShortFlags(os.Args)

	// Определяем флаги
	column := flag.Int("k", 1, "column")
	number := flag.Bool("n", false, "number")
	reverse := flag.Bool("r", false, "reverse")
	unique := flag.Bool("u", false, "unique")
	month := flag.Bool("M", false, "month")
	blanks := flag.Bool("b", false, "blanks")
	checksort := flag.Bool("c", false, "checksort")
	humannumeric := flag.Bool("h", false, "humannumeric")

	flag.Parse()

	// Открываем файл или используем stdin
	if flag.NArg() > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		input = f
		defer f.Close()
	} else {
		input = os.Stdin
	}

	// Чтение строк
	scanner := bufio.NewScanner(input)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Сортировка строк
	sorted := compare.SortLines(lines, *column, *number, *humannumeric, *month, *blanks, *reverse, *unique, *checksort)

	// Вывод результата
	for _, line := range sorted {
		fmt.Println(line)
	}
}
