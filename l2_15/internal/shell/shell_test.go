package shell

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// helper: быстрые Options для тестов
func testOpts(stdin string) (Options, *bytes.Buffer, *bytes.Buffer) {
	var out bytes.Buffer
	var err bytes.Buffer
	return Options{
		Prompt: false,
		Stdin:  strings.NewReader(stdin),
		Stdout: &out,
		Stderr: &err,
	}, &out, &err
}

func TestTokenizeOperators(t *testing.T) {
	toks, err := Tokenize(`ps | grep x && echo ok || echo bad < in.txt > out.txt`)
	if err != nil {
		t.Fatalf("Tokenize error: %v", err)
	}

	// Проверим, что ключевые операторы присутствуют
	has := func(k TokenKind, text string) bool {
		for _, tk := range toks {
			if tk.Kind == k && tk.Text == text {
				return true
			}
		}
		return false
	}

	if !has(TPipe, "|") {
		t.Fatalf("expected token |")
	}
	if !has(TAndAnd, "&&") {
		t.Fatalf("expected token &&")
	}
	if !has(TOrOr, "||") {
		t.Fatalf("expected token ||")
	}
	if !has(TIn, "<") {
		t.Fatalf("expected token <")
	}
	if !has(TOut, ">") {
		t.Fatalf("expected token >")
	}
}

func TestParseChainAndPipeline(t *testing.T) {
	toks, err := Tokenize(`printf "a\nb\n" | wc -l && echo ok`)
	if err != nil {
		t.Fatalf("Tokenize error: %v", err)
	}

	chain, err := ParseChain(toks)
	if err != nil {
		t.Fatalf("ParseChain error: %v", err)
	}
	if len(chain) != 2 {
		t.Fatalf("expected 2 chain items, got %d", len(chain))
	}
	if chain[0].Op != OpNone {
		t.Fatalf("expected first op OpNone")
	}
	if chain[1].Op != OpAnd {
		t.Fatalf("expected second op OpAnd")
	}
	if len(chain[0].P.Cmds) != 2 {
		t.Fatalf("expected pipeline of 2 commands")
	}
	if chain[0].P.Cmds[0].Args[0] != "printf" {
		t.Fatalf("expected first cmd printf, got %s", chain[0].P.Cmds[0].Args[0])
	}
	if chain[0].P.Cmds[1].Args[0] != "wc" {
		t.Fatalf("expected second cmd wc, got %s", chain[0].P.Cmds[1].Args[0])
	}
}

func TestExpandEnvWord(t *testing.T) {
	_ = os.Setenv("TEST_SHELL_VAR", "XYZ")
	got := ExpandEnvWord("a$TEST_SHELL_VAR:b")
	if got != "aXYZ:b" {
		t.Fatalf("expected aXYZ:b, got %q", got)
	}
}

func TestBuiltinCdPwdEcho(t *testing.T) {
	orig, _ := os.Getwd()
	defer func() { _ = os.Chdir(orig) }()

	tmp := t.TempDir()

	opt, out, errBuf := testOpts("")
	sigint := make(chan os.Signal, 1)

	// cd tmp && pwd
	toks, _ := Tokenize("cd " + tmp + " && pwd")
	chain, parseErr := ParseChain(toks)
	if parseErr != nil {
		t.Fatalf("ParseChain error: %v", parseErr)
	}

	code := RunChain(chain, sigint, opt)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, stderr=%q", code, errBuf.String())
	}

	gotPwd := strings.TrimSpace(out.String())

	tmpCanon, err := filepath.EvalSymlinks(tmp)
	if err != nil {
		t.Fatalf("EvalSymlinks(tmp): %v", err)
	}
	gotCanon, err := filepath.EvalSymlinks(gotPwd)
	if err != nil {
		t.Fatalf("EvalSymlinks(gotPwd): %v", err)
	}

	if gotCanon != tmpCanon {
		t.Fatalf("expected pwd %q (canon %q), got %q (canon %q)",
			tmp, tmpCanon, gotPwd, gotCanon)
	}

	// echo
	out.Reset()
	toks, _ = Tokenize("echo hello world")
	chain, _ = ParseChain(toks)
	code = RunChain(chain, sigint, opt)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if strings.TrimSpace(out.String()) != "hello world" {
		t.Fatalf("expected echo output, got %q", out.String())
	}
}

func TestExecExternalSimple(t *testing.T) {
	opt, out, errBuf := testOpts("")
	_ = errBuf
	// Прямой вызов RunExternal
	code := RunExternal([]string{"printf", "hello"}, opt.Stdin, opt.Stdout, opt.Stderr, 0)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if out.String() != "hello" {
		t.Fatalf("expected %q, got %q", "hello", out.String())
	}
}

func TestExecPipeline(t *testing.T) {
	opt, out, errBuf := testOpts("")
	sigint := make(chan os.Signal, 1)

	toks, err := Tokenize(`echo a a a | wc -w`)
	if err != nil {
		t.Fatalf("Tokenize error: %v", err)
	}
	chain, err := ParseChain(toks)
	if err != nil {
		t.Fatalf("ParseChain error: %v", err)
	}

	code := RunChain(chain, sigint, opt)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, stderr=%q", code, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != "3" {
		t.Fatalf("expected 3, got %q", got)
	}
}

func TestRedirectsInOut(t *testing.T) {
	tmp := t.TempDir()
	inPath := filepath.Join(tmp, "in.txt")
	outPath := filepath.Join(tmp, "out.txt")

	// Подготовим входной файл
	if err := os.WriteFile(inPath, []byte("hello\nworld\n"), 0644); err != nil {
		t.Fatalf("write in.txt: %v", err)
	}

	opt, out, errBuf := testOpts("")
	sigint := make(chan os.Signal, 1)
	_ = out // stdout тут не нужен, будет редирект

	// cat < in.txt | wc -l > out.txt
	cmdLine := "cat < " + inPath + " | wc -l > " + outPath
	toks, err := Tokenize(cmdLine)
	if err != nil {
		t.Fatalf("Tokenize error: %v", err)
	}
	chain, err := ParseChain(toks)
	if err != nil {
		t.Fatalf("ParseChain error: %v", err)
	}

	code := RunChain(chain, sigint, opt)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, stderr=%q", code, errBuf.String())
	}

	b, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read out.txt: %v", err)
	}
	got := strings.TrimSpace(string(b))
	if got != "2" {
		t.Fatalf("expected 2, got %q", got)
	}
}

func TestAndOrLogic(t *testing.T) {
	opt, out, errBuf := testOpts("")
	sigint := make(chan os.Signal, 1)

	// false && echo X || echo Y должен вывести Y
	out.Reset()
	toks, _ := Tokenize(`false && echo X || echo Y`)
	chain, err := ParseChain(toks)
	if err != nil {
		t.Fatalf("ParseChain error: %v", err)
	}
	code := RunChain(chain, sigint, opt)
	_ = code
	if strings.TrimSpace(out.String()) != "Y" {
		t.Fatalf("expected Y, got %q, stderr=%q", out.String(), errBuf.String())
	}

	// true && echo OK || echo BAD должен вывести OK
	out.Reset()
	errBuf.Reset()
	toks, _ = Tokenize(`true && echo OK || echo BAD`)
	chain, _ = ParseChain(toks)
	_ = RunChain(chain, sigint, opt)
	if strings.TrimSpace(out.String()) != "OK" {
		t.Fatalf("expected OK, got %q, stderr=%q", out.String(), errBuf.String())
	}
}
