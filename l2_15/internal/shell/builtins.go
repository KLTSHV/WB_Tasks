package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func IsBuiltin(name string) bool {
	switch name {
	case "cd", "pwd", "echo", "kill", "ps":
		return true
	default:
		return false
	}
}

// RunBuiltinParent исполняется в родителе для встроенных команд
func RunBuiltinParent(c CmdSpec, stdin io.Reader, stdout io.Writer) int {
	name := c.Args[0]
	args := c.Args[1:]

	switch name {
	case "cd":
		if len(args) != 1 {
			fmt.Fprintln(stdout, "usage: cd <path>")
			return 2
		}
		if err := os.Chdir(args[0]); err != nil {
			fmt.Fprintf(stdout, "cd: %v\n", err)
			return 1
		}
		return 0

	case "pwd":
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(stdout, "pwd: %v\n", err)
			return 1
		}
		fmt.Fprintln(stdout, wd)
		return 0

	case "echo":
		fmt.Fprintln(stdout, strings.Join(args, " "))
		return 0

	case "kill":
		if len(args) != 1 {
			fmt.Fprintln(stdout, "usage: kill <pid>")
			return 2
		}
		pid, err := strconv.Atoi(args[0])
		if err != nil || pid <= 0 {
			fmt.Fprintln(stdout, "kill: invalid pid")
			return 2
		}
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			fmt.Fprintf(stdout, "kill: %v\n", err)
			return 1
		}
		return 0

	case "ps":
		cmd := exec.Command("ps", "-e", "-o", "pid,comm")
		cmd.Stdin = stdin
		cmd.Stdout = stdout
		cmd.Stderr = stdout
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(stdout, "ps: %v\n", err)
			return 1
		}
		return 0

	default:
		fmt.Fprintf(stdout, "unknown builtin: %s\n", name)
		return 127
	}
}
