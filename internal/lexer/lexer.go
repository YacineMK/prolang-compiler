package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

type Lexer struct {
	source []rune
	pos    int
	line   int
	col    int

	Tokens []Token
	Errors []string
}

func New(src string) *Lexer {
	return &Lexer{
		source: []rune(src),
		pos:    0,
		line:   1,
		col:    1,
	}
}

// ───────────────────────── core ─────────────────────────

func (l *Lexer) Lex() {
	for l.pos < len(l.source) {
		l.scan()
	}
	l.addToken(EOF, "")
}

// ───────────────────────── helpers ─────────────────────────

func (l *Lexer) peek() rune {
	if l.pos >= len(l.source) {
		return 0
	}
	return l.source[l.pos]
}

func (l *Lexer) peekAt(offset int) rune {
	i := l.pos + offset
	if i >= len(l.source) {
		return 0
	}
	return l.source[i]
}

func (l *Lexer) advance() rune {
	ch := l.source[l.pos]
	l.pos++

	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}

	return ch
}

func (l *Lexer) addToken(kind TokenKind, value string) {
	l.Tokens = append(l.Tokens, Token{
		Kind:  kind,
		Value: value,
		Line:  l.line,
		Col:   l.col,
	})
}

func (l *Lexer) addError(msg string) {
	l.Errors = append(l.Errors,
		fmt.Sprintf("LEXICAL_ERROR line %d col %d %s", l.line, l.col, msg))
}

// ───────────────────────── scanning ─────────────────────────

func (l *Lexer) scan() {
	ch := l.peek()

	// whitespace
	if unicode.IsSpace(ch) {
		l.advance()
		return
	}

	// %% comment
	if ch == '%' && l.peekAt(1) == '%' {
		l.advance()
		l.advance()
		for l.peek() != '\n' && l.peek() != 0 {
			l.advance()
		}
		return
	}

	// string
	if ch == '"' {
		l.scanString()
		return
	}

	// number
	if unicode.IsDigit(ch) {
		l.scanNumber()
		return
	}

	// identifier / keyword
	if unicode.IsLetter(ch) {
		l.scanIdent()
		return
	}

	// operators 2-char
	switch ch {
	case '<':
		l.advance()
		if l.peek() == '-' {
			l.advance()
			l.addToken(ASSIGN, "<-")
		} else if l.peek() == '=' {
			l.advance()
			l.addToken(LTE, "<=")
		} else {
			l.addToken(LT, "<")
		}
		return

	case '>':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			l.addToken(GTE, ">=")
		} else {
			l.addToken(GT, ">")
		}
		return

	case '=':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			l.addToken(EQ, "==")
		} else {
			l.addToken(EQUAL, "=")
		}
		return

	case '!':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			l.addToken(NEQ, "!=")
		} else {
			l.addError("unexpected '!'")
		}
		return
	}

	// single char tokens
	single := map[rune]TokenKind{
		'+': PLUS,
		'-': MINUS,
		'*': STAR,
		'/': SLASH,
		'|': PIPE,
		':': COLON,
		';': SEMI,
		',': COMMA,
		'(': LPAREN,
		')': RPAREN,
		'{': LBRACE,
		'}': RBRACE,
		'[': LBRACK,
		']': RBRACK,
	}

	if k, ok := single[ch]; ok {
		l.advance()
		l.addToken(k, string(ch))
		return
	}

	l.addError(fmt.Sprintf("illegal char '%c'", ch))
	l.advance()
}

// ───────────────────────── sub scanners ─────────────────────────

func (l *Lexer) scanNumber() {
	start := l.line
	var sb strings.Builder

	isFloat := false

	for unicode.IsDigit(l.peek()) {
		sb.WriteRune(l.advance())
	}

	if l.peek() == '.' {
		isFloat = true
		sb.WriteRune(l.advance())
		for unicode.IsDigit(l.peek()) {
			sb.WriteRune(l.advance())
		}
	}

	if isFloat {
		l.Tokens = append(l.Tokens, Token{FLOAT_CONST, sb.String(), start, l.col})
	} else {
		l.Tokens = append(l.Tokens, Token{INT_CONST, sb.String(), start, l.col})
	}
}

func (l *Lexer) scanString() {
	l.advance() // "

	var sb strings.Builder

	for l.peek() != '"' && l.peek() != 0 {
		if l.peek() == '\n' {
			l.addError("unterminated string")
			return
		}
		sb.WriteRune(l.advance())
	}

	if l.peek() == '"' {
		l.advance()
		l.addToken(STRING, sb.String())
	} else {
		l.addError("unterminated string")
	}
}

func (l *Lexer) scanIdent() {
	start := l.line
	startCol := l.col

	var sb strings.Builder
	for unicode.IsLetter(l.peek()) || unicode.IsDigit(l.peek()) || l.peek() == '_' {
		sb.WriteRune(l.advance())
	}

	word := sb.String()

	// context keywords
	if word == "in" {
		// Look ahead to see if next non-space char is '('
		p := l.pos
		for p < len(l.source) && unicode.IsSpace(l.source[p]) {
			p++
		}
		if p < len(l.source) && l.source[p] == '(' {
			l.addToken(IN_FUNC, word)
		} else {
			l.addToken(IN, word)
		}
		return
	}
	if word == "out" {
		l.addToken(OUT_FUNC, word)
		return
	}

	if k, ok := Keywords[word]; ok {
		l.Tokens = append(l.Tokens, Token{k, word, start, startCol})
		return
	}

	l.Tokens = append(l.Tokens, Token{IDENT, word, start, startCol})
}
