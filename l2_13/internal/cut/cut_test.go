package cut

import (
	"flag"
	"os"
	"testing"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestParseFlags_SingleFields(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "-f", "1,3,5"}
	opt := ParseFlags()

	expected := map[int]struct{}{1: {}, 3: {}, 5: {}}
	if len(opt.Fields) != len(expected) {
		t.Fatalf("expected %d fields, got %d", len(expected), len(opt.Fields))
	}

	for k := range expected {
		if _, ok := opt.Fields[k]; !ok {
			t.Errorf("expected field %d to be present", k)
		}
	}
}

func TestParseFlags_Ranges(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "-f", "2-4"}
	opt := ParseFlags()

	for i := 2; i <= 4; i++ {
		if _, ok := opt.Fields[i]; !ok {
			t.Errorf("expected field %d to be present", i)
		}
	}
}

func TestParseFlags_Combination(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "-f", "1,3-5"}
	opt := ParseFlags()

	expected := map[int]struct{}{1: {}, 3: {}, 4: {}, 5: {}}
	if len(opt.Fields) != len(expected) {
		t.Fatalf("expected %d fields, got %d", len(expected), len(opt.Fields))
	}
}

func TestProcessLine_Basic(t *testing.T) {
	opt := Options{
		Fields:    map[int]struct{}{1: {}, 3: {}},
		Delimiter: "\t",
	}

	input := "a\tb\tc\td"
	out := ProcessLine(input, opt)
	expected := "a\tc"

	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestProcessLine_OnlySeparated(t *testing.T) {
	opt := Options{
		Fields:        map[int]struct{}{1: {}, 2: {}},
		Delimiter:     "\t",
		OnlySeparated: true,
	}

	input := "no_separator_here"
	out := ProcessLine(input, opt)

	if out != "" {
		t.Errorf("expected empty output, got %q", out)
	}
}

func TestProcessLine_FieldOutOfBounds(t *testing.T) {
	opt := Options{
		Fields:    map[int]struct{}{1: {}, 3: {}, 5: {}},
		Delimiter: ",",
	}

	input := "x,y"
	out := ProcessLine(input, opt)
	expected := "x"

	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestProcessLine_EmptyFields(t *testing.T) {
	opt := Options{
		Fields:    map[int]struct{}{2: {}, 4: {}},
		Delimiter: ",",
	}

	input := "a,,b,"
	out := ProcessLine(input, opt)
	expected := "," // 2 и 4 поле пустые

	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestProcessLine_DifferentDelimiter(t *testing.T) {
	opt := Options{
		Fields:    map[int]struct{}{2: {}, 3: {}},
		Delimiter: ";",
	}

	input := "one;two;three;four"
	out := ProcessLine(input, opt)
	expected := "two;three"

	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}
