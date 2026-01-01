package shell

import "fmt"

type CmdSpec struct {
	Args    []string
	InFile  string
	OutFile string
}

type Pipeline struct {
	Cmds []CmdSpec
}

type ChainOp int

const (
	OpNone ChainOp = iota
	OpAnd
	OpOr
)

type ChainItem struct {
	P  Pipeline
	Op ChainOp // оператор связывает предыдущий результат с текущим
}

func ParseChain(tokens []Token) ([]ChainItem, error) {
	// chain := pipeline ( (&& || ||) pipeline )*
	var items []ChainItem
	pos := 0

	parsePipeline := func() (Pipeline, error) {
		var p Pipeline

		parseCmd := func() (CmdSpec, error) {
			var c CmdSpec
			for pos < len(tokens) {
				t := tokens[pos]
				switch t.Kind {
				case TWord:
					c.Args = append(c.Args, ExpandEnvWord(t.Text))
					pos++
				case TIn:
					pos++
					if pos >= len(tokens) || tokens[pos].Kind != TWord {
						return CmdSpec{}, fmt.Errorf("expected file after '<'")
					}
					c.InFile = ExpandEnvWord(tokens[pos].Text)
					pos++
				case TOut:
					pos++
					if pos >= len(tokens) || tokens[pos].Kind != TWord {
						return CmdSpec{}, fmt.Errorf("expected file after '>'")
					}
					c.OutFile = ExpandEnvWord(tokens[pos].Text)
					pos++
				case TPipe, TAndAnd, TOrOr:
					return c, nil
				default:
					return CmdSpec{}, fmt.Errorf("unexpected token: %s", t.Text)
				}
			}
			return c, nil
		}

		for {
			cmd, err := parseCmd()
			if err != nil {
				return Pipeline{}, err
			}
			if len(cmd.Args) == 0 {
				return Pipeline{}, fmt.Errorf("empty command in pipeline")
			}
			p.Cmds = append(p.Cmds, cmd)

			if pos >= len(tokens) || tokens[pos].Kind != TPipe {
				break
			}
			pos++ // skip |
			if pos >= len(tokens) {
				return Pipeline{}, fmt.Errorf("expected command after '|'")
			}
		}

		return p, nil
	}

	p0, err := parsePipeline()
	if err != nil {
		return nil, err
	}
	items = append(items, ChainItem{P: p0, Op: OpNone})

	for pos < len(tokens) {
		t := tokens[pos]
		var op ChainOp
		if t.Kind == TAndAnd {
			op = OpAnd
		} else if t.Kind == TOrOr {
			op = OpOr
		} else {
			return nil, fmt.Errorf("expected && or ||, got: %s", t.Text)
		}
		pos++
		if pos >= len(tokens) {
			return nil, fmt.Errorf("expected command after %s", t.Text)
		}
		pn, err := parsePipeline()
		if err != nil {
			return nil, err
		}
		items = append(items, ChainItem{P: pn, Op: op})
	}

	return items, nil
}
