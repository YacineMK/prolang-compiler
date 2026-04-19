package semantic

import (
	"go_compiler/internal/lexer"
	"go_compiler/internal/quad"
)

// ════════════════════════════════════════════════════════════════
//  QUAD GENERATOR (IR GENERATION)
// ════════════════════════════════════════════════════════════════

type QuadGenerator struct {
	analyzer *Analyzer
	qm       *quad.Manager
	tempVars map[string]string // Map of parsed variables to track
}

// NewQuadGenerator creates a new quad generator
func NewQuadGenerator(analyzer *Analyzer) *QuadGenerator {
	return &QuadGenerator{
		analyzer: analyzer,
		qm:       quad.NewManager(),
		tempVars: make(map[string]string),
	}
}

// Generate generates quads from the analyzed program
func (qg *QuadGenerator) Generate() *quad.Manager {
	// We need to re-parse and track instructions
	// Since semantic analyzer only validates, we need to walk through tokens again
	qg.generateFromTokens()
	return qg.qm
}

// generateFromTokens walks through tokens and generates quads
func (qg *QuadGenerator) generateFromTokens() {
	tokens := qg.analyzer.tokens
	pos := 0

	// Find RUN block
	for pos < len(tokens) {
		if tokens[pos].Kind == lexer.RUN {
			pos++ // skip RUN
			if pos < len(tokens) && tokens[pos].Kind == lexer.COLON {
				pos++ // skip :
			}
			if pos < len(tokens) && tokens[pos].Kind == lexer.LBRACE {
				pos++ // skip {
				// Generate quads for statements in Run block
				pos = qg.generateStatements(tokens, pos)
			}
			break
		}
		pos++
	}
}

// generateStatements processes statements and generates quads
func (qg *QuadGenerator) generateStatements(tokens []lexer.Token, pos int) int {
	// Process statements until we hit }
	for pos < len(tokens) {
		tok := tokens[pos]

		switch tok.Kind {
		case lexer.RBRACE:
			return pos + 1

		case lexer.IDENT:
			// Check if assignment
			if pos+1 < len(tokens) {
				nextTok := tokens[pos+1]
				if nextTok.Kind == lexer.ASSIGN {
					varName := tok.Value
					pos += 2 // skip var and <-

					// Generate quads for the right-hand side expression
					result := qg.generateExpression(tokens, &pos)

					// Emit assignment
					qg.qm.Emit(":=", result, "", varName)

					// Skip semicolon
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

// generateExpression handles expressions and returns the result variable
func (qg *QuadGenerator) generateExpression(tokens []lexer.Token, pos *int) string {
	if *pos >= len(tokens) {
		return ""
	}

	// Parse left operand
	left := qg.getPrimaryValue(tokens, pos)

	// Check for binary operations
	for *pos < len(tokens) {
		if *pos >= len(tokens) {
			break
		}

		tok := tokens[*pos]

		switch tok.Kind {
		case lexer.PLUS, lexer.MINUS, lexer.STAR, lexer.SLASH,
			lexer.LT, lexer.GT, lexer.LTE, lexer.GTE, lexer.EQ, lexer.NEQ:

			op := tok.Value
			*pos++ // skip operator

			right := qg.getPrimaryValue(tokens, pos)

			// Emit binary operation
			result := qg.qm.NewTemp()
			qg.qm.Emit(op, left, right, result)
			left = result

		default:
			return left
		}
	}

	return left
}

// getPrimaryValue gets a primary value (constant, variable, or parenthesized expression)
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
		// Handle parenthesized expression
		result := qg.generateExpression(tokens, pos)
		// Skip closing )
		if *pos < len(tokens) && tokens[*pos].Kind == lexer.RPAREN {
			(*pos)++
		}
		return result
	}

	return ""
}

// generateOutCall handles output function calls
func (qg *QuadGenerator) generateOutCall(tokens []lexer.Token, pos int) int {
	pos++ // skip 'out'

	if pos < len(tokens) && tokens[pos].Kind == lexer.LPAREN {
		pos++ // skip (

		// Collect arguments
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

	// Skip semicolon
	if pos < len(tokens) && tokens[pos].Kind == lexer.SEMI {
		pos++
	}

	return pos
}

// generateInCall handles input function calls
func (qg *QuadGenerator) generateInCall(tokens []lexer.Token, pos int) int {
	pos++ // skip 'in'

	if pos < len(tokens) && tokens[pos].Kind == lexer.LPAREN {
		pos++ // skip (

		// Collect arguments
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

	// Skip semicolon
	if pos < len(tokens) && tokens[pos].Kind == lexer.SEMI {
		pos++
	}

	return pos
}
