package semantic

import (
	"go_compiler/internal/lexer"
	"go_compiler/internal/quad"
)

type QuadGenerator struct {
	analyzer *Analyzer
	qm       *quad.Manager
	tempVars map[string]string
}

func NewQuadGenerator(analyzer *Analyzer) *QuadGenerator {
	return &QuadGenerator{
		analyzer: analyzer,
		qm:       quad.NewManager(),
		tempVars: make(map[string]string),
	}
}

func (qg *QuadGenerator) Generate() *quad.Manager {
	qg.generateFromTokens()
	return qg.qm
}

func (qg *QuadGenerator) generateFromTokens() {
	tokens := qg.analyzer.tokens
	pos := 0

	for pos < len(tokens) {
		if tokens[pos].Kind == lexer.RUN {
			pos++
			if pos < len(tokens) && tokens[pos].Kind == lexer.COLON {
				pos++
			}
			if pos < len(tokens) && tokens[pos].Kind == lexer.LBRACE {
				pos++
				pos = qg.generateStatements(tokens, pos)
			}
			break
		}
		pos++
	}
}

func (qg *QuadGenerator) generateStatements(tokens []lexer.Token, pos int) int {
	for pos < len(tokens) {
		tok := tokens[pos]

		switch tok.Kind {
		case lexer.RBRACE:
			return pos + 1

		case lexer.IDENT:
			if pos+1 < len(tokens) {
				nextTok := tokens[pos+1]
				if nextTok.Kind == lexer.ASSIGN {
					varName := tok.Value
					pos += 2

					result := qg.generateExpression(tokens, &pos)

					qg.qm.Emit(":=", result, "", varName)

					if pos < len(tokens) && tokens[pos].Kind == lexer.SEMI {
						pos++
					}
					continue
				}
			}

		case lexer.OUT_FUNC:
			pos = qg.generateOutCall(tokens, pos)
			continue

		case lexer.IN_FUNC:
			pos = qg.generateInCall(tokens, pos)
			continue
		}
		pos++
	}
	return pos
}

func (qg *QuadGenerator) generateExpression(tokens []lexer.Token, pos *int) string {
	if *pos >= len(tokens) {
		return ""
	}

	left := qg.getPrimaryValue(tokens, pos)

	for *pos < len(tokens) {
		if *pos >= len(tokens) {
			break
		}

		tok := tokens[*pos]

		switch tok.Kind {
		case lexer.PLUS, lexer.MINUS, lexer.STAR, lexer.SLASH,
			lexer.LT, lexer.GT, lexer.LTE, lexer.GTE, lexer.EQ, lexer.NEQ:

			op := tok.Value
			*pos++

			right := qg.getPrimaryValue(tokens, pos)

			result := qg.qm.NewTemp()
			qg.qm.Emit(op, left, right, result)
			left = result

		default:
			return left
		}
	}

	return left
}

func (qg *QuadGenerator) getPrimaryValue(tokens []lexer.Token, pos *int) string {
	if *pos >= len(tokens) {
		return ""
	}

	tok := tokens[*pos]
	*pos++

	switch tok.Kind {
	case lexer.INT_CONST, lexer.FLOAT_CONST:
		return tok.Value
	case lexer.IDENT:
		return tok.Value
	case lexer.STRING:
		return tok.Value
	case lexer.LPAREN:
		result := qg.generateExpression(tokens, pos)
		if *pos < len(tokens) && tokens[*pos].Kind == lexer.RPAREN {
			(*pos)++
		}
		return result
	}

	return ""
}

func (qg *QuadGenerator) generateOutCall(tokens []lexer.Token, pos int) int {
	pos++

	if pos < len(tokens) && tokens[pos].Kind == lexer.LPAREN {
		pos++

		for pos < len(tokens) {
			if tokens[pos].Kind == lexer.RPAREN {
				pos++
				break
			}

			if tokens[pos].Kind != lexer.COMMA {
				arg := tokens[pos].Value
				qg.qm.Emit("out", arg, "", "")
			}

			pos++
		}
	}

	if pos < len(tokens) && tokens[pos].Kind == lexer.SEMI {
		pos++
	}

	return pos
}

func (qg *QuadGenerator) generateInCall(tokens []lexer.Token, pos int) int {
	pos++

	if pos < len(tokens) && tokens[pos].Kind == lexer.LPAREN {
		pos++

		for pos < len(tokens) {
			if tokens[pos].Kind == lexer.RPAREN {
				pos++
				break
			}

			if tokens[pos].Kind != lexer.COMMA {
				arg := tokens[pos].Value
				qg.qm.Emit("in", arg, "", "")
			}

			pos++
		}
	}

	if pos < len(tokens) && tokens[pos].Kind == lexer.SEMI {
		pos++
	}

	return pos
}
