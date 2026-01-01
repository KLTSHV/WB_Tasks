package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
)

type Options struct {
	Prompt bool
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type Shell struct {
	opt    Options
	reader *bufio.Reader
	sigint chan os.Signal
}

func New(opt Options) *Shell {
	if opt.Stdin == nil {
		opt.Stdin = os.Stdin
	}
	if opt.Stdout == nil {
		opt.Stdout = os.Stdout
	}
	if opt.Stderr == nil {
		opt.Stderr = os.Stderr
	}

	sh := &Shell{
		opt:    opt,
		reader: bufio.NewReader(opt.Stdin),
		sigint: make(chan os.Signal, 16),
	}
	signal.Notify(sh.sigint, os.Interrupt)
	return sh
}

func (sh *Shell) Run() {
	for {
		if sh.opt.Prompt {
			sh.printPrompt()
		}
		line, err := sh.reader.ReadString('\n')
		if err != nil {
			// Ctrl+D
			if errors.Is(err, io.EOF) {
				fmt.Fprintln(sh.opt.Stdout)
				return
			}
			fmt.Fprintf(sh.opt.Stderr, "read: %v\n", err)
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" {
			return
		}

		toks, err := Tokenize(line)
		if err != nil {
			fmt.Fprintf(sh.opt.Stderr, "parse: %v\n", err)
			continue
		}

		chain, err := ParseChain(toks)
		if err != nil {
			fmt.Fprintf(sh.opt.Stderr, "parse: %v\n", err)
			continue
		}

		// Выполнение
		_ = RunChain(chain, sh.sigint, sh.opt)
	}
}

func (sh *Shell) printPrompt() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprint(sh.opt.Stdout, "minishell> ")
		return
	}
	fmt.Fprintf(sh.opt.Stdout, "%s minishell> ", wd)
}
