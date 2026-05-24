package lexer

type TokenKind string

const (
	INT_CONST   TokenKind = "INT_CONST"
	FLOAT_CONST TokenKind = "FLOAT_CONST"
	STRING      TokenKind = "STRING"
	IDENT       TokenKind = "IDENT"

	BEGINPROJECT TokenKind = "BEGINPROJECT"
	ENDPROJECT   TokenKind = "ENDPROJECT"
	SETUP        TokenKind = "SETUP"
	RUN          TokenKind = "RUN"
	DEFINE       TokenKind = "DEFINE"
	CONST        TokenKind = "CONST"
	INTEGER      TokenKind = "INTEGER"
	FLOAT        TokenKind = "FLOAT"
	IF           TokenKind = "IF"
	THEN         TokenKind = "THEN"
	ELSE         TokenKind = "ELSE"
	ENDIF        TokenKind = "ENDIF"
	LOOP         TokenKind = "LOOP"
	WHILE        TokenKind = "WHILE"
	ENDLOOP      TokenKind = "ENDLOOP"
	FOR          TokenKind = "FOR"
	IN           TokenKind = "IN"
	TO           TokenKind = "TO"
	ENDFOR       TokenKind = "ENDFOR"
	AND          TokenKind = "AND"
	OR           TokenKind = "OR"
	NON          TokenKind = "NON"

	ASSIGN TokenKind = "ASSIGN"
	PLUS   TokenKind = "PLUS"
	MINUS  TokenKind = "MINUS"
	STAR   TokenKind = "STAR"
	SLASH  TokenKind = "SLASH"

	LT  TokenKind = "LT"
	GT  TokenKind = "GT"
	LTE TokenKind = "LTE"
	GTE TokenKind = "GTE"
	EQ  TokenKind = "EQ"
	NEQ TokenKind = "NEQ"

	PIPE  TokenKind = "PIPE"
	COLON TokenKind = "COLON"
	SEMI  TokenKind = "SEMI"
	COMMA TokenKind = "COMMA"

	LPAREN TokenKind = "LPAREN"
	RPAREN TokenKind = "RPAREN"
	LBRACE TokenKind = "LBRACE"
	RBRACE TokenKind = "RBRACE"
	LBRACK TokenKind = "LBRACK"
	RBRACK TokenKind = "RBRACK"

	EQUAL TokenKind = "EQUAL"

	IN_FUNC  TokenKind = "IN_FUNC"
	OUT_FUNC TokenKind = "OUT_FUNC"

	EOF     TokenKind = "EOF"
	ILLEGAL TokenKind = "ILLEGAL"
)

type Token struct {
	Kind  TokenKind
	Value string
	Line  int
	Col   int
}

var Keywords = map[string]TokenKind{
	"BeginProject": BEGINPROJECT,
	"EndProject":   ENDPROJECT,
	"Setup":        SETUP,
	"Run":          RUN,
	"define":       DEFINE,
	"const":        CONST,
	"integer":      INTEGER,
	"float":        FLOAT,
	"if":           IF,
	"then":         THEN,
	"else":         ELSE,
	"endIf":        ENDIF,
	"loop":         LOOP,
	"while":        WHILE,
	"endloop":      ENDLOOP,
	"for":          FOR,
	"to":           TO,
	"endfor":       ENDFOR,
	"AND":          AND,
	"OR":           OR,
	"NON":          NON,
	"input":        IN_FUNC,
}
