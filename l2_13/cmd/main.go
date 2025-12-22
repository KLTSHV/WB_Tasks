package main

import (
	"bufio"
	"l2_13/internal/cut"
	"os"
)

func main() {
	opt := cut.ParseFlags()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		output := cut.ProcessLine(line, opt)
		if output != "" {
			println(output)
		}

	}

}
