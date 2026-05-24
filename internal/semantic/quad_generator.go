package semantic

import (
	"fmt"
	"go_compiler/internal/lexer"
	"go_compiler/internal/quad"
)

type QuadGenerator struct {
	analyzer *Analyzer
	qm       *quad.Manager
	tokens   []lexer.Token
	pos      int
}

func NewQuadGenerator(analyzer *Analyzer) *QuadGenerator {
	return &QuadGenerator{
		analyzer: analyzer,
		qm:       quad.NewManager(),
		tokens:   analyzer.tokens,
	}
}

func (qg *QuadGenerator) Generate() *quad.Manager {
	qg.skipToRunBlock()
	qg.generateStatements()
	return qg.qm
}

func (qg *QuadGenerator) cur() lexer.Token {
	if qg.pos >= len(qg.tokens) {
		return lexer.Token{Kind: lexer.EOF}
	}
	return qg.tokens[qg.pos]
}

func (qg *QuadGenerator) peek(offset int) lexer.Token {
	idx := qg.pos + offset
	if idx >= len(qg.tokens) {
		return lexer.Token{Kind: lexer.EOF}
	}
	return qg.tokens[idx]
}

func (qg *QuadGenerator) consume() lexer.Token {
	t := qg.cur()
	qg.pos++
	return t
}

func (qg *QuadGenerator) expect(kind lexer.TokenKind) lexer.Token {
	if qg.cur().Kind == kind {
		return qg.consume()
	}
	return lexer.Token{}
}

func (qg *QuadGenerator) skipToRunBlock() {
	for qg.cur().Kind != lexer.EOF {
		if qg.cur().Kind == lexer.RUN {
			qg.consume()
			qg.expect(lexer.COLON)
			qg.expect(lexer.LBRACE)
			return
		}
		qg.consume()
	}
}

func (qg *QuadGenerator) generateStatements() {
	for {
		switch qg.cur().Kind {
		case lexer.RBRACE, lexer.EOF,
			lexer.ENDLOOP, lexer.ENDFOR, lexer.ENDIF, lexer.ELSE:
			return

		case lexer.IDENT:
			qg.generateAssignment()

		case lexer.IF:
			qg.generateIf()

		case lexer.LOOP:
			qg.generateLoopWhile()

		case lexer.FOR:
			qg.generateFor()

		case lexer.IN_FUNC:
			qg.generateInCall()

		case lexer.OUT_FUNC:
			qg.generateOutCall()

		default:
			qg.consume()
		}
	}
}

func (qg *QuadGenerator) generateAssignment() {
	nameTok := qg.consume()
	varName := nameTok.Value

	if qg.cur().Kind == lexer.LBRACK {
		qg.consume()
		idxPlace := qg.generateExpr()
		qg.expect(lexer.RBRACK)
		qg.expect(lexer.ASSIGN)
		valPlace := qg.generateExpr()
		qg.expect(lexer.SEMI)
		qg.qm.Emit("[]=", valPlace, idxPlace, varName)
		return
	}

	qg.expect(lexer.ASSIGN)
	result := qg.generateExpr()
	qg.expect(lexer.SEMI)
	qg.qm.Emit(":=", result, "", varName)
}

func (qg *QuadGenerator) generateIf() {
	qg.consume()
	qg.expect(lexer.LPAREN)
	cond := qg.generateCondition()
	qg.expect(lexer.RPAREN)
	qg.expect(lexer.THEN)
	qg.expect(lexer.COLON)
	qg.expect(lexer.LBRACE)

	bfIdx := qg.qm.NextIndex()
	qg.qm.Emit("BF", cond, "", "")

	qg.generateStatements()
	qg.expect(lexer.RBRACE)

	hasElse := qg.cur().Kind == lexer.ELSE

	if hasElse {
		brIdx := qg.qm.NextIndex()
		qg.qm.Emit("BR", "", "", "")

		qg.qm.Backpatch(bfIdx, fmt.Sprintf("%d", qg.qm.NextIndex()))

		qg.consume()
		qg.expect(lexer.LBRACE)
		qg.generateStatements()
		qg.expect(lexer.RBRACE)

		qg.qm.Backpatch(brIdx, fmt.Sprintf("%d", qg.qm.NextIndex()))
	} else {
		qg.qm.Backpatch(bfIdx, fmt.Sprintf("%d", qg.qm.NextIndex()))
	}

	qg.expect(lexer.ENDIF)
	qg.expect(lexer.SEMI)
}

func (qg *QuadGenerator) generateLoopWhile() {
	qg.consume()
	qg.expect(lexer.WHILE)
	qg.expect(lexer.LPAREN)

	startIndex := qg.qm.NextIndex()

	cond := qg.generateCondition()
	qg.expect(lexer.RPAREN)
	qg.expect(lexer.LBRACE)

	bfIdx := qg.qm.NextIndex()
	qg.qm.Emit("BF", cond, "", "")

	qg.generateStatements()
	qg.expect(lexer.RBRACE)
	qg.expect(lexer.ENDLOOP)
	qg.expect(lexer.SEMI)

	qg.qm.Emit("BR", "", "", fmt.Sprintf("%d", startIndex))

	qg.qm.Backpatch(bfIdx, fmt.Sprintf("%d", qg.qm.NextIndex()))
}

func (qg *QuadGenerator) generateFor() {
	qg.consume()

	idfTok := qg.consume()
	idf := idfTok.Value

	qg.expect(lexer.IN)
	initVal := qg.generateExpr()
	qg.expect(lexer.TO)
	limitVal := qg.generateExpr()
	qg.expect(lexer.LBRACE)

	qg.qm.Emit(":=", initVal, "", idf)

	startIndex := qg.qm.NextIndex()

	tCmp := qg.qm.NewTemp()
	qg.qm.Emit("<=", idf, limitVal, tCmp)

	bfIdx := qg.qm.NextIndex()
	qg.qm.Emit("BF", tCmp, "", "")

	qg.generateStatements()
	qg.expect(lexer.RBRACE)
	qg.expect(lexer.ENDFOR)
	qg.expect(lexer.SEMI)

	tInc := qg.qm.NewTemp()
	qg.qm.Emit("+", idf, "1", tInc)
	qg.qm.Emit(":=", tInc, "", idf)

	qg.qm.Emit("BR", "", "", fmt.Sprintf("%d", startIndex))

	qg.qm.Backpatch(bfIdx, fmt.Sprintf("%d", qg.qm.NextIndex()))
}

func (qg *QuadGenerator) generateInCall() {
	qg.consume()
	qg.expect(lexer.LPAREN)

	if qg.cur().Kind != lexer.RPAREN && qg.cur().Kind != lexer.EOF {
		varTok := qg.consume()
		varName := varTok.Value

		if qg.cur().Kind == lexer.LBRACK {
			qg.consume()
			idxPlace := qg.generateExpr()
			qg.expect(lexer.RBRACK)
			tRead := qg.qm.NewTemp()
			qg.qm.Emit("in", tRead, "", "")
			qg.qm.Emit("[]=", tRead, idxPlace, varName)
		} else {
			qg.qm.Emit("in", "", "", varName)
		}
	}

	qg.expect(lexer.RPAREN)
	qg.expect(lexer.SEMI)
}

func (qg *QuadGenerator) generateOutCall() {
	qg.consume()
	qg.expect(lexer.LPAREN)

	for qg.cur().Kind != lexer.RPAREN && qg.cur().Kind != lexer.EOF {
		if qg.cur().Kind == lexer.COMMA {
			qg.consume()
			continue
		}

		if qg.cur().Kind == lexer.STRING {
			strVal := qg.consume().Value
			qg.qm.Emit("out_str", fmt.Sprintf("%q", strVal), "", "")
		} else {
			place := qg.generateExpr()
			qg.qm.Emit("out", place, "", "")
		}
	}

	qg.expect(lexer.RPAREN)
	qg.expect(lexer.SEMI)
}

func (qg *QuadGenerator) generateCondition() string {
	return qg.generateOrCond()
}

func (qg *QuadGenerator) generateOrCond() string {
	left := qg.generateAndCond()
	for qg.cur().Kind == lexer.OR {
		qg.consume()
		right := qg.generateAndCond()
		t := qg.qm.NewTemp()
		qg.qm.Emit("OR", left, right, t)
		left = t
	}
	return left
}

func (qg *QuadGenerator) generateAndCond() string {
	left := qg.generateNotCond()
	for qg.cur().Kind == lexer.AND {
		qg.consume()
		right := qg.generateNotCond()
		t := qg.qm.NewTemp()
		qg.qm.Emit("AND", left, right, t)
		left = t
	}
	return left
}

func (qg *QuadGenerator) generateNotCond() string {
	if qg.cur().Kind == lexer.NON {
		qg.consume()
		qg.expect(lexer.LPAREN)
		inner := qg.generateCondition()
		qg.expect(lexer.RPAREN)
		t := qg.qm.NewTemp()
		qg.qm.Emit("NON", inner, "", t)
		return t
	}
	return qg.generatePrimaryCond()
}

func (qg *QuadGenerator) generatePrimaryCond() string {
	if qg.cur().Kind == lexer.LPAREN && qg.isLogicalGroup() {
		qg.consume()
		inner := qg.generateCondition()
		qg.expect(lexer.RPAREN)
		return inner
	}

	left := qg.generateExpr()

	switch qg.cur().Kind {
	case lexer.LT, lexer.GT, lexer.LTE, lexer.GTE, lexer.EQ, lexer.NEQ:
		op := qg.consume().Value
		right := qg.generateExpr()
		t := qg.qm.NewTemp()
		qg.qm.Emit(op, left, right, t)
		return t
	}

	return left
}

func (qg *QuadGenerator) isLogicalGroup() bool {
	depth := 0
	for i := qg.pos; i < len(qg.tokens); i++ {
		k := qg.tokens[i].Kind
		switch k {
		case lexer.LPAREN:
			depth++
		case lexer.RPAREN:
			depth--
			if depth == 0 {
				return false
			}
		case lexer.AND, lexer.OR, lexer.NON:
			if depth == 1 {
				return true
			}
		case lexer.LT, lexer.GT, lexer.LTE, lexer.GTE, lexer.EQ, lexer.NEQ:
			if depth == 1 {
				return true
			}
		case lexer.EOF:
			return false
		}
	}
	return false
}

func (qg *QuadGenerator) generateExpr() string {
	return qg.generateAddSub()
}

func (qg *QuadGenerator) generateAddSub() string {
	left := qg.generateMulDiv()
	for qg.cur().Kind == lexer.PLUS || qg.cur().Kind == lexer.MINUS {
		op := qg.consume().Value
		right := qg.generateMulDiv()
		t := qg.qm.NewTemp()
		qg.qm.Emit(op, left, right, t)
		left = t
	}
	return left
}

func (qg *QuadGenerator) generateMulDiv() string {
	left := qg.generateUnary()
	for qg.cur().Kind == lexer.STAR || qg.cur().Kind == lexer.SLASH {
		op := qg.consume().Value
		right := qg.generateUnary()
		t := qg.qm.NewTemp()
		qg.qm.Emit(op, left, right, t)
		left = t
	}
	return left
}

func (qg *QuadGenerator) generateUnary() string {
	if qg.cur().Kind == lexer.MINUS {
		qg.consume()
		operand := qg.generatePrimary()
		t := qg.qm.NewTemp()
		qg.qm.Emit("NEG", operand, "", t)
		return t
	}
	if qg.cur().Kind == lexer.PLUS {
		qg.consume()
	}
	return qg.generatePrimary()
}

func (qg *QuadGenerator) generatePrimary() string {
	tok := qg.cur()

	switch tok.Kind {
	case lexer.INT_CONST, lexer.FLOAT_CONST:
		qg.consume()
		return tok.Value

	case lexer.STRING:
		qg.consume()
		return fmt.Sprintf("%q", tok.Value)

	case lexer.IDENT:
		qg.consume()
		if qg.cur().Kind == lexer.LBRACK {
			qg.consume()
			idxPlace := qg.generateExpr()
			qg.expect(lexer.RBRACK)
			t := qg.qm.NewTemp()
			qg.qm.Emit("[]", tok.Value, idxPlace, t)
			return t
		}
		return tok.Value

	case lexer.LPAREN:
		qg.consume()

		isSignedLiteral := (qg.cur().Kind == lexer.PLUS || qg.cur().Kind == lexer.MINUS) &&
			qg.pos+1 < len(qg.tokens) &&
			(qg.tokens[qg.pos+1].Kind == lexer.INT_CONST ||
				qg.tokens[qg.pos+1].Kind == lexer.FLOAT_CONST) &&
			qg.pos+2 < len(qg.tokens) &&
			qg.tokens[qg.pos+2].Kind == lexer.RPAREN

		if isSignedLiteral {
			sign := qg.consume().Value
			val := qg.consume()
			qg.expect(lexer.RPAREN)
			if sign == "-" {
				t := qg.qm.NewTemp()
				qg.qm.Emit("NEG", val.Value, "", t)
				return t
			}
			return val.Value
		}

		inner := qg.generateExpr()
		qg.expect(lexer.RPAREN)
		return inner

	default:
		qg.consume()
		return ""
	}
}
