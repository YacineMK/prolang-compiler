package semantic

import (
	"fmt"
	"go_compiler/internal/ast"
	"go_compiler/internal/lexer"
	"strconv"
)

// ════════════════════════════════════════════════════════════════
//  ANALYSEUR SÉMANTIQUE (SEMANTIC ANALYZER)
// ════════════════════════════════════════════════════════════════

type Analyzer struct {
	tokens []lexer.Token
	pos    int
	Table  *ast.SymbolTable
	Errors []string
}

// NewAnalyzer crée un nouvel analyseur sémantique
func NewAnalyzer(tokens []lexer.Token, table *ast.SymbolTable) *Analyzer {
	return &Analyzer{tokens: tokens, Table: table}
}

// ─── Navigation dans les tokens ──────────────────────────────────

func (a *Analyzer) cur() lexer.Token {
	if a.pos >= len(a.tokens) {
		return lexer.Token{Kind: lexer.EOF}
	}
	return a.tokens[a.pos]
}

func (a *Analyzer) consume() lexer.Token {
	t := a.cur()
	a.pos++
	return t
}

func (a *Analyzer) expect(kind lexer.TokenKind) (lexer.Token, bool) {
	t := a.cur()
	if t.Kind != kind {
		a.errorf(t, "attendu '%s', trouvé '%s'", kind, t.Value)
		return t, false
	}
	return a.consume(), true
}

// ─── Gestion des erreurs ────────────────────────────────────────

func (a *Analyzer) errorf(t lexer.Token, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	a.Errors = append(a.Errors, fmt.Sprintf("SEMANTIC_ERROR, ligne %d, colonne %d, %s [entité: '%s']",
		t.Line, t.Col, msg, t.Value))
}

// ─── Analyse principale ──────────────────────────────────────────

// Analyze effectue l'analyse sémantique complète du programme
func (a *Analyzer) Analyze() {
	a.expect(lexer.BEGINPROJECT)
	a.expect(lexer.IDENT) // nom du programme
	a.expect(lexer.SEMI)

	a.expect(lexer.SETUP)
	a.expect(lexer.COLON)
	a.parseDeclarations()

	a.expect(lexer.RUN)
	a.expect(lexer.COLON)
	a.expect(lexer.LBRACE)
	a.parseInstructions()
	a.expect(lexer.RBRACE)

	a.expect(lexer.ENDPROJECT)
	a.expect(lexer.SEMI)
}

// ─── Déclarations ───────────────────────────────────────────────

func (a *Analyzer) parseDeclarations() {
	for {
		switch a.cur().Kind {
		case lexer.DEFINE:
			a.parseDefine()
		case lexer.CONST:
			a.parseConst()
		default:
			return
		}
	}
}

func (a *Analyzer) parseDefine() {
	a.consume() // define

	// Noms : a | b | c
	var names []lexer.Token
	if tok, ok := a.expect(lexer.IDENT); ok {
		names = append(names, tok)
	} else {
		return
	}
	for a.cur().Kind == lexer.PIPE {
		a.consume()
		if tok, ok := a.expect(lexer.IDENT); ok {
			names = append(names, tok)
		}
	}

	a.expect(lexer.COLON)

	// Tableau : [type ; taille]
	if a.cur().Kind == lexer.LBRACK {
		a.consume()
		typeTok, ok := a.parseType()
		if !ok {
			return
		}
		a.expect(lexer.SEMI)
		sizeTok, ok := a.expect(lexer.INT_CONST)
		if !ok {
			return
		}
		a.expect(lexer.RBRACK)
		a.expect(lexer.SEMI)

		size, _ := strconv.Atoi(sizeTok.Value)
		if size <= 0 {
			a.errorf(sizeTok, "la taille du tableau doit être un entier positif")
		}
		for _, n := range names {
			sym := ast.Symbol{
				Name:        n.Value,
				Kind:        ast.KIND_ARRAY,
				Type:        ast.DataType(typeTok.Value),
				ArraySize:   size,
				Line:        n.Line,
				Column:      n.Col,
				Initialized: true,
			}
			if err := a.Table.Insert(sym); err != nil {
				a.errorf(n, err.Error())
			}
		}
		return
	}

	// Variable simple : define a : integer [= val] ;
	typeTok, ok := a.parseType()
	if !ok {
		return
	}
	initVal := ""
	initialized := false
	if a.cur().Kind == lexer.EQUAL {
		a.consume()
		v, ok := a.parseLiteral(ast.DataType(typeTok.Value))
		if ok {
			initVal = v
			initialized = true
		}
	}
	a.expect(lexer.SEMI)

	for _, n := range names {
		sym := ast.Symbol{
			Name:        n.Value,
			Kind:        ast.KIND_VAR,
			Type:        ast.DataType(typeTok.Value),
			Value:       initVal,
			Initialized: initialized,
			Line:        n.Line,
			Column:      n.Col,
		}
		if err := a.Table.Insert(sym); err != nil {
			a.errorf(n, err.Error())
		}
	}
}

func (a *Analyzer) parseConst() {
	a.consume() // const
	nameTok, ok := a.expect(lexer.IDENT)
	if !ok {
		return
	}
	a.expect(lexer.COLON)
	typeTok, ok := a.parseType()
	if !ok {
		return
	}
	a.expect(lexer.EQUAL)
	val, ok := a.parseLiteral(ast.DataType(typeTok.Value))
	if !ok {
		return
	}
	a.expect(lexer.SEMI)

	sym := ast.Symbol{
		Name:        nameTok.Value,
		Kind:        ast.KIND_CONST,
		Type:        ast.DataType(typeTok.Value),
		Value:       val,
		Initialized: true,
		Line:        nameTok.Line,
		Column:      nameTok.Col,
	}
	if err := a.Table.Insert(sym); err != nil {
		a.errorf(nameTok, err.Error())
	}
}

// ─── Instructions ────────────────────────────────────────────────

func (a *Analyzer) parseInstructions() {
	for {
		switch a.cur().Kind {
		case lexer.IDENT:
			a.parseAssignment()
		case lexer.IF:
			a.parseIf()
		case lexer.LOOP:
			a.parseLoopWhile()
		case lexer.FOR:
			a.parseFor()
		case lexer.IN_FUNC:
			a.parseInFunc()
		case lexer.OUT_FUNC:
			a.parseOutFunc()
		default:
			return
		}
	}
}

func (a *Analyzer) parseAssignment() {
	nameTok := a.consume()

	sym := a.Table.Lookup(nameTok.Value)
	if sym == nil {
		a.errorf(nameTok, "identificateur '%s' non déclaré", nameTok.Value)
	}
	if sym != nil && sym.Kind == ast.KIND_CONST {
		a.errorf(nameTok, "'%s' est une constante, affectation interdite", nameTok.Value)
	}

	// Accès tableau ?
	if a.cur().Kind == lexer.LBRACK {
		a.consume()
		a.parseExpr()
		a.expect(lexer.RBRACK)
		if sym != nil && sym.Kind != ast.KIND_ARRAY {
			a.errorf(nameTok, "'%s' n'est pas un tableau", nameTok.Value)
		}
	}

	a.expect(lexer.ASSIGN)
	exprType := a.parseExpr()
	a.expect(lexer.SEMI)

	if sym != nil && sym.Kind != ast.KIND_ARRAY {
		if exprType == ast.TYPE_FLOAT && sym.Type == ast.TYPE_INTEGER {
			a.errorf(nameTok, "incompatibilité de type : float affecté à integer '%s'", nameTok.Value)
		}
	}
	if sym != nil {
		a.Table.MarkInitialized(nameTok.Value)
	}
}

func (a *Analyzer) parseIf() {
	a.consume() // if
	a.expect(lexer.LPAREN)
	a.parseCondition()
	a.expect(lexer.RPAREN)
	a.expect(lexer.THEN)
	a.expect(lexer.COLON)
	a.expect(lexer.LBRACE)
	a.parseInstructions()
	a.expect(lexer.RBRACE)
	if a.cur().Kind == lexer.ELSE {
		a.consume()
		a.expect(lexer.LBRACE)
		a.parseInstructions()
		a.expect(lexer.RBRACE)
	}
	a.expect(lexer.ENDIF)
	a.expect(lexer.SEMI)
}

func (a *Analyzer) parseLoopWhile() {
	a.consume() // loop
	a.expect(lexer.WHILE)
	a.expect(lexer.LPAREN)
	a.parseCondition()
	a.expect(lexer.RPAREN)
	a.expect(lexer.LBRACE)
	a.parseInstructions()
	a.expect(lexer.RBRACE)
	a.expect(lexer.ENDLOOP)
	a.expect(lexer.SEMI)
}

func (a *Analyzer) parseFor() {
	a.consume() // for
	idfTok, ok := a.expect(lexer.IDENT)
	if !ok {
		return
	}
	sym := a.Table.Lookup(idfTok.Value)
	if sym == nil {
		a.errorf(idfTok, "variable de boucle '%s' non déclarée", idfTok.Value)
	} else if sym.Type != ast.TYPE_INTEGER {
		a.errorf(idfTok, "variable de boucle '%s' doit être de type integer", idfTok.Value)
	}
	a.expect(lexer.IN)
	a.parseExpr() // val_init
	a.expect(lexer.TO)
	a.parseExpr() // val_limit
	a.expect(lexer.LBRACE)
	a.parseInstructions()
	a.expect(lexer.RBRACE)
	a.expect(lexer.ENDFOR)
	a.expect(lexer.SEMI)
}

func (a *Analyzer) parseInFunc() {
	a.consume() // in
	a.expect(lexer.LPAREN)
	nameTok, ok := a.expect(lexer.IDENT)
	if ok && a.Table.Lookup(nameTok.Value) == nil {
		a.errorf(nameTok, "identificateur '%s' non déclaré", nameTok.Value)
	}
	a.expect(lexer.RPAREN)
	a.expect(lexer.SEMI)
}

func (a *Analyzer) parseOutFunc() {
	a.consume() // out
	a.expect(lexer.LPAREN)
	for {
		t := a.cur()
		if t.Kind == lexer.STRING {
			a.consume()
		} else if t.Kind == lexer.IDENT {
			nameTok := a.consume()
			if a.Table.Lookup(nameTok.Value) == nil {
				a.errorf(nameTok, "identificateur '%s' non déclaré", nameTok.Value)
			}
		} else {
			break
		}
		if a.cur().Kind != lexer.COMMA {
			break
		}
		a.consume()
	}
	a.expect(lexer.RPAREN)
	a.expect(lexer.SEMI)
}

// ─── Expressions ─────────────────────────────────────────────────

func (a *Analyzer) parseExpr() ast.DataType {
	return a.parseAddSub()
}

func (a *Analyzer) parseAddSub() ast.DataType {
	t := a.parseMulDiv()
	for a.cur().Kind == lexer.PLUS || a.cur().Kind == lexer.MINUS {
		a.consume()
		t2 := a.parseMulDiv()
		t = ast.MergeTypes(t, t2)
	}
	return t
}

func (a *Analyzer) parseMulDiv() ast.DataType {
	t := a.parsePrimary()
	for a.cur().Kind == lexer.STAR || a.cur().Kind == lexer.SLASH {
		a.consume()
		t2 := a.parsePrimary()
		t = ast.MergeTypes(t, t2)
	}
	return t
}

func (a *Analyzer) parsePrimary() ast.DataType {
	t := a.cur()
	switch t.Kind {
	case lexer.INT_CONST:
		a.consume()
		return ast.TYPE_INTEGER
	case lexer.FLOAT_CONST:
		a.consume()
		return ast.TYPE_FLOAT
	case lexer.IDENT:
		a.consume()
		sym := a.Table.Lookup(t.Value)
		if sym == nil {
			a.errorf(t, "identificateur '%s' non déclaré", t.Value)
			return ast.TYPE_UNKNOWN
		}
		if a.cur().Kind == lexer.LBRACK {
			a.consume()
			a.parseExpr()
			a.expect(lexer.RBRACK)
		}
		return sym.Type
	case lexer.LPAREN:
		a.consume()
		// Constante signée : (+5) ou (-3.14)
		if a.cur().Kind == lexer.PLUS || a.cur().Kind == lexer.MINUS {
			a.consume()
		}
		inner := a.parseExpr()
		a.expect(lexer.RPAREN)
		return inner
	default:
		a.errorf(t, "expression attendue, trouvé '%s'", t.Value)
		a.consume()
		return ast.TYPE_UNKNOWN
	}
}

func (a *Analyzer) parseCondition() {
	if a.cur().Kind == lexer.NON {
		a.consume()
		a.expect(lexer.LPAREN)
		a.parseCondition()
		a.expect(lexer.RPAREN)
		return
	}
	a.parseExpr()
	switch a.cur().Kind {
	case lexer.LT, lexer.GT, lexer.LTE, lexer.GTE, lexer.EQ, lexer.NEQ:
		a.consume()
		a.parseExpr()
	case lexer.AND, lexer.OR:
		a.consume()
		a.parseCondition()
	}
}

// ─── Utilitaires ─────────────────────────────────────────────────

func (a *Analyzer) parseType() (lexer.Token, bool) {
	t := a.cur()
	if t.Kind == lexer.INTEGER || t.Kind == lexer.FLOAT {
		return a.consume(), true
	}
	a.errorf(t, "type attendu (integer|float), trouvé '%s'", t.Value)
	return t, false
}

func (a *Analyzer) parseLiteral(expected ast.DataType) (string, bool) {
	t := a.cur()
	// Constante signée entre parenthèses : (+5) ou (-3.14)
	if t.Kind == lexer.LPAREN {
		a.consume()
		sign := ""
		if a.cur().Kind == lexer.PLUS || a.cur().Kind == lexer.MINUS {
			sign = a.consume().Value
		}
		val := a.cur()
		if val.Kind != lexer.INT_CONST && val.Kind != lexer.FLOAT_CONST {
			a.errorf(val, "constante numérique attendue")
			return "", false
		}
		a.consume()
		a.expect(lexer.RPAREN)
		if expected == ast.TYPE_INTEGER && val.Kind == lexer.FLOAT_CONST {
			a.errorf(val, "type incompatible : float assigné à integer")
		}
		return sign + val.Value, true
	}
	if t.Kind == lexer.INT_CONST || t.Kind == lexer.FLOAT_CONST {
		a.consume()
		if expected == ast.TYPE_INTEGER && t.Kind == lexer.FLOAT_CONST {
			a.errorf(t, "type incompatible : float assigné à integer")
		}
		return t.Value, true
	}
	a.errorf(t, "littéral numérique attendu, trouvé '%s'", t.Value)
	return "", false
}
