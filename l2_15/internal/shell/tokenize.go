package shell

import (
	"bytes"
	"fmt"
)

type TokenKind int

const (
	TWord TokenKind = iota
	TPipe
	TAndAnd
	TOrOr
	TIn
	TOut
)

type Token struct {
	Kind TokenKind
	Text string
}

func Tokenize(line string) ([]Token, error) {
	var toks []Token
	i := 0

	flushWord := func(buf *bytes.Buffer) {
		if buf.Len() > 0 {
			toks = append(toks, Token{Kind: TWord, Text: buf.String()})
			buf.Reset()
		}
	}

	var buf bytes.Buffer
	for i < len(line) {
		ch := line[i]
		switch ch {
		case ' ', '\t', '\n', '\r':
			flushWord(&buf)
			i++
		case '|':
			flushWord(&buf)
			if i+1 < len(line) && line[i+1] == '|' {
				toks = append(toks, Token{Kind: TOrOr, Text: "||"})
				i += 2
			} else {
				toks = append(toks, Token{Kind: TPipe, Text: "|"})
				i++
			}
		case '&':
			flushWord(&buf)
			if i+1 < len(line) && line[i+1] == '&' {
				toks = append(toks, Token{Kind: TAndAnd, Text: "&&"})
				i += 2
			} else {
				return nil, fmt.Errorf("unexpected '&' (only && supported)")
			}
		case '<':
			flushWord(&buf)
			toks = append(toks, Token{Kind: TIn, Text: "<"})
			i++
		case '>':
			flushWord(&buf)
			toks = append(toks, Token{Kind: TOut, Text: ">"})
			i++
		default:
			buf.WriteByte(ch)
			i++
		}
	}

	flushWord(&buf)
	return toks, nil
}
