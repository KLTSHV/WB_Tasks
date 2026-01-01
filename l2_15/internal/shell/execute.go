package shell

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

type pipePair struct {
	r *os.File
	w *os.File
}

func RunChain(items []ChainItem, sigint <-chan os.Signal, opt Options) int {
	status := 0
	for idx, it := range items {
		if idx == 0 {
			status = RunPipeline(it.P, sigint, opt)
			continue
		}
		switch it.Op {
		case OpAnd:
			if status == 0 {
				status = RunPipeline(it.P, sigint, opt)
			}
		case OpOr:
			if status != 0 {
				status = RunPipeline(it.P, sigint, opt)
			}
		default:
			status = RunPipeline(it.P, sigint, opt)
		}
	}
	return status
}

func RunPipeline(p Pipeline, sigint <-chan os.Signal, opt Options) int {
	// Одиночный builtin в родителе
	if len(p.Cmds) == 1 && IsBuiltin(p.Cmds[0].Args[0]) {
		return runSingleBuiltin(p.Cmds[0], opt)
	}
	return runExternalPipeline(p, sigint, opt)
}

func runSingleBuiltin(c CmdSpec, opt Options) int {
	inF, outF, err := openRedirectsSingle(c)
	if err != nil {
		fmt.Fprintf(opt.Stderr, "redirect: %v\n", err)
		return 1
	}
	defer func() {
		if inF != nil {
			_ = inF.Close()
		}
		if outF != nil {
			_ = outF.Close()
		}
	}()

	var in io.Reader = opt.Stdin
	var out io.Writer = opt.Stdout
	if inF != nil {
		in = inF
	}
	if outF != nil {
		out = outF
	}
	return RunBuiltinParent(c, in, out)
}

func openRedirectsSingle(c CmdSpec) (*os.File, *os.File, error) {
	var inF *os.File
	var outF *os.File
	if c.InFile != "" {
		f, err := os.Open(c.InFile)
		if err != nil {
			return nil, nil, err
		}
		inF = f
	}
	if c.OutFile != "" {
		f, err := os.OpenFile(c.OutFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			if inF != nil {
				_ = inF.Close()
			}
			return nil, nil, err
		}
		outF = f
	}
	return inF, outF, nil
}

// RunExternal: запускает одну внешнюю команду,
// pgidOverride=0 — создать новую группу, иначе присоединиться к группе
func RunExternal(argv []string, stdin io.Reader, stdout, stderr io.Writer, pgidOverride int) int {
	if len(argv) == 0 {
		return 127
	}
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if pgidOverride != 0 {
		cmd.SysProcAttr.Pgid = pgidOverride
	}
	if err := cmd.Run(); err != nil {
		return exitCode(err)
	}
	return 0
}

func runExternalPipeline(p Pipeline, sigint <-chan os.Signal, opt Options) int {
	n := len(p.Cmds)

	// Запретим cd в pipeline (иначе ожидания пользователя будут неверными)
	for _, c := range p.Cmds {
		if len(c.Args) > 0 && c.Args[0] == "cd" {
			fmt.Fprintln(opt.Stderr, "cd: cannot be used in pipeline")
			return 2
		}
	}

	// Создаём пайпы между командами
	pipes := make([]pipePair, 0, max(0, n-1))
	for i := 0; i < n-1; i++ {
		r, w, err := os.Pipe()
		if err != nil {
			fmt.Fprintf(opt.Stderr, "pipe: %v\n", err)
			return 1
		}
		pipes = append(pipes, pipePair{r: r, w: w})
	}

	// Редиректы: < только у первой, > только у последней (MVP)
	inF, outF, err := openPipelineEndRedirects(p)
	if err != nil {
		closeAllPipes(pipes)
		fmt.Fprintf(opt.Stderr, "redirect: %v\n", err)
		return 1
	}
	defer func() {
		if inF != nil {
			_ = inF.Close()
		}
		if outF != nil {
			_ = outF.Close()
		}
	}()

	cmds := make([]*exec.Cmd, 0, n)

	// Собираем exec.Cmd
	for i := 0; i < n; i++ {
		c := p.Cmds[i]
		name := c.Args[0]
		args := c.Args[1:]

		cmd := exec.Command(name, args...)
		cmd.Stderr = opt.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		// stdin
		if i == 0 {
			if inF != nil {
				cmd.Stdin = inF
			} else {
				cmd.Stdin = opt.Stdin
			}
		} else {
			cmd.Stdin = pipes[i-1].r
		}
		// stdout
		if i == n-1 {
			if outF != nil {
				cmd.Stdout = outF
			} else {
				cmd.Stdout = opt.Stdout
			}
		} else {
			cmd.Stdout = pipes[i].w
		}
		cmds = append(cmds, cmd)
	}

	// Стартуем процессы, формируем одну process group для всего pipeline
	var pgid int
	for i, cmd := range cmds {
		if i == 0 {
			if err := cmd.Start(); err != nil {
				fmt.Fprintf(opt.Stderr, "%s: %v\n", cmd.Path, err)
				closeAllPipes(pipes)
				return 127
			}
			pgid = cmd.Process.Pid
		} else {
			cmd.SysProcAttr.Pgid = pgid
			if err := cmd.Start(); err != nil {
				fmt.Fprintf(opt.Stderr, "%s: %v\n", cmd.Path, err)
				_ = syscall.Kill(-pgid, syscall.SIGTERM)
				closeAllPipes(pipes)
				return 127
			}
		}
	}

	// В родителе закрываем все pipe FD, иначе можно повиснуть на EOF
	closeAllPipes(pipes)

	// SIGINT всей группе при Ctrl+C
	go func() {
		for range sigint {
			if pgid != 0 {
				_ = syscall.Kill(-pgid, syscall.SIGINT)
			}
		}
	}()

	// Wait всех, статус берём по последней команде
	lastStatus := 0
	for i, cmd := range cmds {
		err := cmd.Wait()
		code := exitCode(err)
		if i == len(cmds)-1 {
			lastStatus = code
		}
	}
	return lastStatus
}

func openPipelineEndRedirects(p Pipeline) (*os.File, *os.File, error) {
	var inF *os.File
	var outF *os.File

	first := p.Cmds[0]
	last := p.Cmds[len(p.Cmds)-1]

	if first.InFile != "" {
		f, err := os.Open(first.InFile)
		if err != nil {
			return nil, nil, err
		}
		inF = f
	}
	if last.OutFile != "" {
		f, err := os.OpenFile(last.OutFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			if inF != nil {
				_ = inF.Close()
			}
			return nil, nil, err
		}
		outF = f
	}
	return inF, outF, nil
}

func closeAllPipes(pipes []pipePair) {
	for _, p := range pipes {
		_ = p.r.Close()
		_ = p.w.Close()
	}
}

// exitCode: 0 если nil; 128+signal если завершено сигналом; иначе ExitStatus; иначе 127.
func exitCode(err error) int {
	if err == nil {
		return 0
	}

	var ee *exec.ExitError
	if errors.As(err, &ee) {
		if ws, ok := ee.Sys().(syscall.WaitStatus); ok {
			if ws.Signaled() {
				return 128 + int(ws.Signal())
			}
			return ws.ExitStatus()
		}
		return 1
	}
	return 127
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
