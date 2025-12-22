package cut

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type Options struct {
	Fields        map[int]struct{}
	Delimiter     string
	OnlySeparated bool
}

func ParseFlags() Options {
	var opt Options
	flag.StringVar(&opt.Delimiter, "d", "\t", "delimiter")
	flag.BoolVar(&opt.OnlySeparated, "s", false, "only separated")
	flag.Func("f", "fields", func(s string) error { //записываем номера колонок которые нужно вывести
		parts := strings.Split(s, ",")
		opt.Fields = make(map[int]struct{}) //с мапой поиск будет быстрее и гарантируется уникальность номеров
		for _, part := range parts {
			if strings.Contains(part, "-") { //если есть дефис, значит это диапазон
				rangeParts := strings.SplitN(part, "-", 2)
				if rangeParts[0] == "" || rangeParts[1] == "" { //проверка на корректность диапазона
					return fmt.Errorf("invalid range: %s", part)
				}
				start, err := strconv.Atoi(rangeParts[0])
				if err != nil || start < 1 { //номера колонок начинаются с 1
					return fmt.Errorf("invalid start: %s", rangeParts[0])
				}
				end, err := strconv.Atoi(rangeParts[1])
				if err != nil || end < start { //проверка на корректность конца диапазона
					return fmt.Errorf("invalid end: %s", rangeParts[1])
				}
				for i := start; i <= end; i++ { //добавляем все номера из диапазона в мапу
					opt.Fields[i] = struct{}{}
				}
			} else { //если нет дефиса, значит это одиночный номер
				num, err := strconv.Atoi(part)
				if err != nil || num < 1 {
					return fmt.Errorf("invalid field: %s", part)
				}
				opt.Fields[num] = struct{}{}
			}
		}
		return nil
	})

	flag.Parse()
	return opt
}

func ProcessLine(line string, opt Options) string {
	fields := strings.Split(line, opt.Delimiter)
	if len(fields) == 1 && opt.OnlySeparated { //если нет разделителя и установлен флаг -s то пропускаем строку
		return ""
	}
	var selected []string
	for i := 1; i <= len(fields); i++ {
		if _, ok := opt.Fields[i]; ok { //проверяем есть ли номер колонки в мапе
			selected = append(selected, fields[i-1])
		}
	}
	output := strings.Join(selected, opt.Delimiter)
	return output
}
