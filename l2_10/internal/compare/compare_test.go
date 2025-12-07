package compare

import "testing"

func TestDate(t *testing.T) {
	lines := []string{
		"Jan 12",
		"Feb 05",
		"Mar 23",
		"Apr 12",
		"May 15",
		"Mar 31",
		"Oct 25",
		"Nov 11",
		"Dec 31",
		"Jan 01",
	}
	sorted := SortLines(lines, 1, false, false, true, false, false, false, false)
	expected := []string{
		"Jan 01",
		"Jan 12",
		"Feb 05",
		"Mar 23",
		"Mar 31",
		"Apr 12",
		"May 15",
		"Oct 25",
		"Nov 11",
		"Dec 31",
	}
	for i, line := range sorted {
		if line != expected[i] {
			t.Errorf("expected %q, got %q", expected[i], line)
		}
	}

}

func TestCompare(t *testing.T) {
	lines := []string{
		"apple",
		"Banana",
		"cherry",
		"date",
	}
	sorted := SortLines(lines, 1, false, false, false, false, false, false, false)
	expected := []string{
		"Banana",
		"apple",
		"cherry",
		"date",
	}
	for i, line := range sorted {
		if line != expected[i] {
			t.Errorf("expected %q, got %q", expected[i], line)
		}
	}
}
func TestCompareReverse(t *testing.T) {
	lines := []string{
		"apple",
		"Banana",
		"cherry",
		"date",
	}
	sorted := SortLines(lines, 1, false, false, false, false, true, false, false)
	expected := []string{
		"date",
		"cherry",
		"apple",
		"Banana",
	}
	for i, line := range sorted {
		if line != expected[i] {
			t.Errorf("expected %q, got %q", expected[i], line)
		}
	}
}
func TestCompareNumeric(t *testing.T) {
	lines := []string{
		"10 apples",
		"2 bananas",
		"30 cherries",
		"4 dates",
	}
	sorted := SortLines(lines, 1, true, false, false, false, false, false, false)
	expected := []string{
		"2 bananas",
		"4 dates",
		"10 apples",
		"30 cherries",
	}
	for i, line := range sorted {
		if line != expected[i] {
			t.Errorf("expected %q, got %q", expected[i], line)
		}
	}
}
func TestCompareNumericReverse(t *testing.T) {
	lines := []string{
		"10 apples",
		"2 bananas",
		"30 cherries",
		"4 dates",
	}
	sorted := SortLines(lines, 1, true, false, false, false, true, false, false)
	expected := []string{
		"30 cherries",
		"10 apples",
		"4 dates",
		"2 bananas",
	}
	for i, line := range sorted {
		if line != expected[i] {
			t.Errorf("expected %q, got %q", expected[i], line)
		}
	}
}

func TestCompareBlanks(t *testing.T) {
	lines := []string{
		"   apple",
		"wb   ",
		"  wb  ",
		" hello ",
	}
	sorted := SortLines(lines, 1, false, false, false, true, false, false, false)
	expected := []string{
		"   apple",
		" hello ",
		"wb   ",
		"  wb  ",
	}
	for i, line := range sorted {
		if line != expected[i] {
			t.Errorf("expected %q, got %q", expected[i], line)
		}
	}
}
func TestCompareBlanksReverse(t *testing.T) {
	lines := []string{
		"   1apple",
		"wb   ",
		"  wb  ",
		" hello ",
		"	",
	}
	sorted := SortLines(lines, 1, false, false, false, true, true, false, false)
	expected := []string{
		"	",
		"wb   ",
		"  wb  ",
		" hello ",
		"   1apple",
	}
	for i, line := range sorted {
		if line != expected[i] {
			t.Errorf("expected %q, got %q", expected[i], line)
		}
	}
}

func TestCompareHumannumeric(t *testing.T) {
	lines := []string{
		"1K files",
		"500 files",
		"2M files",
		"1G files",
		"750K files",
	}
	sorted := SortLines(lines, 1, false, true, false, false, false, false, false)
	expected := []string{
		"500 files",
		"1K files",
		"750K files",
		"2M files",
		"1G files",
	}
	for i, line := range sorted {
		if line != expected[i] {
			t.Errorf("expected %q, got %q", expected[i], line)
		}
	}
}
func TestCompareHumannumericReverse(t *testing.T) {
	lines := []string{
		"1K files",
		"500 files",
		"2M files",
		"1G files",
		"750K files",
	}
	sorted := SortLines(lines, 1, false, true, false, false, true, false, false)
	expected := []string{
		"1G files",
		"2M files",
		"750K files",
		"1K files",
		"500 files",
	}
	for i, line := range sorted {
		if line != expected[i] {
			t.Errorf("expected %q, got %q", expected[i], line)
		}
	}
}
func TestExpandShortFlags(t *testing.T) {
	input := []string{"cmd", "-nrMu", "file.txt"}
	expected := []string{"cmd", "-n", "-r", "-M", "-u", "file.txt"}
	result := ExpandShortFlags(input)
	if len(result) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(result))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %q, got %q", i, expected[i], v)
		}
	}
}
func TestUnique(t *testing.T) {
	lines := []string{
		"apple",
		"lamp",
		"grass",
		"golang",
		"lamp",
		"date",
	}
	sorted := SortLines(lines, 1, false, false, false, false, false, true, false)
	expected := []string{
		"apple",
		"date",
		"golang",
		"grass",
		"lamp",
	}
	for i, line := range sorted {
		if line != expected[i] {
			t.Errorf("expected %q, got %q", expected[i], line)
		}
	}
}

func TestChecksort(t *testing.T) {
	lines := []string{
		"apple",
		"banana",
		"cherry",
		"date",
	}
	sorted := SortLines(lines, 1, false, false, false, false, false, false, true)
	expected := []string{
		"apple",
		"banana",
		"cherry",
		"date",
	}
	for i, line := range sorted {
		if line != expected[i] {
			t.Errorf("expected %q, got %q", expected[i], line)
		}
	}
}
func TestChecksortFail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got none")
		}
	}()

	lines := []string{
		"banana",
		"apple",
		"cherry",
		"date",
	}
	SortLines(lines, 1, false, false, false, false, false, false, true)
}
