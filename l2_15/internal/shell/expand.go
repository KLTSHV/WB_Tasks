package shell

import (
	"os"
	"strings"
)

// ExpandEnvWord простая подстановка $VAR внутри слова
func ExpandEnvWord(w string) string {
	if !strings.Contains(w, "$") {
		return w
	}
	var out strings.Builder
	for i := 0; i < len(w); {
		if w[i] != '$' {
			out.WriteByte(w[i])
			i++
			continue
		}
		i++ // skip $
		if i >= len(w) {
			out.WriteByte('$')
			break
		}
		j := i
		for j < len(w) {
			c := w[j]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
				j++
			} else {
				break
			}
		}
		if j == i {
			out.WriteByte('$')
			continue
		}
		name := w[i:j]
		out.WriteString(os.Getenv(name))
		i = j
	}
	return out.String()
}
