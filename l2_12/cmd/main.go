package main

import (
	"flag"
	"l2_12/internal/grep"
	"log"
	"os"
)

// Разворачивание флагов (-iF в -i -F)
func expandShortFlags(args []string) []string {
	var expanded []string
	for _, arg := range args {
		if len(arg) > 2 && arg[0] == '-' && arg[1] != '-' {
			for _, ch := range arg[1:] {
				expanded = append(expanded, "-"+string(ch))
			}
		} else {
			expanded = append(expanded, arg)
		}
	}
	return expanded
}

func parseFlags() grep.Options {
	var opt grep.Options
	flag.IntVar(&opt.After, "A", 0, "lines after match")
	flag.IntVar(&opt.Before, "B", 0, "lines before match")
	c := flag.Int("C", 0, "lines of context")
	flag.BoolVar(&opt.CountOnly, "c", false, "count matches only")
	flag.BoolVar(&opt.IgnoreCase, "i", false, "ignore case")
	flag.BoolVar(&opt.Invert, "v", false, "invert match")
	flag.BoolVar(&opt.Fixed, "F", false, "fixed string search")
	flag.BoolVar(&opt.LineNumbers, "n", false, "print line numbers")

	os.Args = append([]string{os.Args[0]}, expandShortFlags(os.Args[1:])...)
	flag.Parse()

	if *c > 0 {
		opt.Before = *c
		opt.After = *c
	}

	return opt
}

func main() {
	opt := parseFlags()
	args := flag.Args()

	if len(args) < 1 {
		log.Fatal("pattern required")
	}

	pattern := args[0]
	files := args[1:]

	if len(files) > 0 {
		for _, f := range files {
			grep.ProcessFile(f, pattern, opt)
		}
	} else {
		grep.ProcessStdin(pattern, opt)
	}
}
