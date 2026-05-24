package asm

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"go_compiler/internal/ast"
	"go_compiler/internal/quad"
)

type Generator struct {
	quads   []quad.Quad
	table   *ast.SymbolTable
	temps   map[string]bool
	strings map[string]string // label -> content
	strIdx  int
}

func NewGenerator(quads []quad.Quad, table *ast.SymbolTable) *Generator {
	return &Generator{
		quads:   quads,
		table:   table,
		temps:   make(map[string]bool),
		strings: make(map[string]string),
	}
}

func (g *Generator) Generate(path string) error {
	g.collectTemps()
	g.collectStrings()

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create %s: %w", path, err)
	}
	defer f.Close()

	g.writeHeader(f)
	g.writeStackSegment(f)
	g.writeDataSegment(f)
	g.writeCodeSegment(f)

	return nil
}

// collectTemps scans quads for temporary variables that need declarations.
func (g *Generator) collectTemps() {
	for _, q := range g.quads {
		for _, field := range []string{q.Arg1, q.Arg2, q.Result} {
			if isTemp(field) {
				g.temps[field] = true
			}
		}
	}
}

// collectStrings scans out_str quads and assigns labels.
func (g *Generator) collectStrings() {
	for _, q := range g.quads {
		if q.Op == "out_str" {
			content, err := strconv.Unquote(q.Arg1)
			if err != nil {
				content = strings.Trim(q.Arg1, `"`)
			}
			if _, exists := g.strings[content]; !exists {
				label := fmt.Sprintf("str%d", g.strIdx)
				g.strIdx++
				g.strings[content] = label
			}
		}
	}
}

func (g *Generator) writeHeader(f *os.File) {
	f.WriteString("; Generated_Assembly\n\n")
}

func (g *Generator) writeStackSegment(f *os.File) {
	f.WriteString("STACK SEGMENT STACK\n")
	f.WriteString("    DW 100 DUP(?)\n")
	f.WriteString("STACK ENDS\n\n")
}

func (g *Generator) writeDataSegment(f *os.File) {
	f.WriteString("DATA SEGMENT\n")

	g.writeUserSymbols(f)
	g.writeTemporaries(f)
	g.writeStrConstants(f)

	f.WriteString("DATA ENDS\n\n")
}

func (g *Generator) writeUserSymbols(f *os.File) {
	f.WriteString("    ; user variables\n")

	for _, sym := range g.table.Symbols() {
		switch sym.Kind {
		case ast.KIND_VAR:
			if sym.Initialized && sym.Value != "" {
				val := integerize(sym.Value)
				f.WriteString(fmt.Sprintf("%s DW %s\n", sym.Name, val))
			} else {
				f.WriteString(fmt.Sprintf("%s DW ?\n", sym.Name))
			}
		case ast.KIND_CONST:
			val := integerize(sym.Value)
			f.WriteString(fmt.Sprintf("%s DW %s\n", sym.Name, val))
		case ast.KIND_ARRAY:
			if sym.ArraySize > 0 {
				f.WriteString(fmt.Sprintf("%s DW %d DUP(?)\n", sym.Name, sym.ArraySize))
			}
		}
	}
}

func (g *Generator) writeTemporaries(f *os.File) {
	if len(g.temps) == 0 {
		return
	}
	f.WriteString("\n    ; compiler-generated temporaries\n")

	names := make([]string, 0, len(g.temps))
	for t := range g.temps {
		names = append(names, t)
	}
	sortTemps(names)
	for _, t := range names {
		f.WriteString(fmt.Sprintf("%s DW ?\n", t))
	}
}

func (g *Generator) writeStrConstants(f *os.File) {
	if len(g.strings) == 0 {
		return
	}
	f.WriteString("\n    ; string constants\n")

	keys := make([]string, 0, len(g.strings))
	for k := range g.strings {
		keys = append(keys, k)
	}
	sortStrings(keys, g.strings)
	for _, content := range keys {
		label := g.strings[content]
		f.WriteString(fmt.Sprintf("%s DB '%s', '$'\n", label, content))
	}
}

func (g *Generator) writeCodeSegment(f *os.File) {
	f.WriteString("CODE SEGMENT\n")
	f.WriteString("    ASSUME CS:CODE, DS:DATA, SS:STACK\n\n")

	f.WriteString("START:\n")
	f.WriteString("    MOV AX, DATA\n")
	f.WriteString("    MOV DS, AX\n\n")

	for i, q := range g.quads {
		f.WriteString(fmt.Sprintf("L%d:\n", i))
		g.emitQuad(f, i, q)
		f.WriteString("\n")
	}

	f.WriteString("    MOV AH, 4Ch\n")
	f.WriteString("    INT 21h\n\n")

	g.writePrintRoutine(f)

	f.WriteString("CODE ENDS\n")
	f.WriteString("END START\n")
}

func (g *Generator) emitQuad(f *os.File, idx int, q quad.Quad) {
	switch q.Op {
	case ":=":
		g.emitAssign(f, q)
	case "+":
		g.emitArith(f, q, "ADD")
	case "-":
		g.emitArith(f, q, "SUB")
	case "*":
		g.emitMul(f, q)
	case "/":
		g.emitDiv(f, q)
	case "NEG":
		g.emitNeg(f, q)
	case ">", "<", ">=", "<=", "==", "!=":
		g.emitCmp(f, idx, q)
	case "AND":
		g.emitLogical(f, q, "AND")
	case "OR":
		g.emitLogical(f, q, "OR")
	case "NON":
		g.emitNot(f, q)
	case "BF":
		g.emitBF(f, q)
	case "BR":
		f.WriteString(fmt.Sprintf("    JMP L%s\n", q.Result))
	case "[]":
		g.emitArrayRead(f, q)
	case "[]=":
		g.emitArrayWrite(f, q)
	case "in":
		g.emitIn(f, q)
	case "out":
		g.emitOut(f, q)
	case "out_str":
		g.emitOutStr(f, q)
	case "LABEL":
		// label is informational only
	}
}

func (g *Generator) emitAssign(f *os.File, q quad.Quad) {
	src := asmVal(q.Arg1)
	dst := q.Result

	if isNum(q.Arg1) {
		f.WriteString(fmt.Sprintf("    MOV %s, %s\n", dst, src))
	} else {
		f.WriteString(fmt.Sprintf("    MOV AX, %s\n", src))
		f.WriteString(fmt.Sprintf("    MOV %s, AX\n", dst))
	}
}

func (g *Generator) emitArith(f *os.File, q quad.Quad, insn string) {
	a1, a2, res := asmVal(q.Arg1), asmVal(q.Arg2), q.Result

	f.WriteString(fmt.Sprintf("    MOV AX, %s\n", a1))
	if isNum(q.Arg2) {
		f.WriteString(fmt.Sprintf("    %s AX, %s\n", insn, a2))
	} else {
		f.WriteString(fmt.Sprintf("    MOV CX, %s\n", a2))
		f.WriteString(fmt.Sprintf("    %s AX, CX\n", insn))
	}
	f.WriteString(fmt.Sprintf("    MOV %s, AX\n", res))
}

func (g *Generator) emitMul(f *os.File, q quad.Quad) {
	a1, a2, res := asmVal(q.Arg1), asmVal(q.Arg2), q.Result

	f.WriteString(fmt.Sprintf("    MOV AX, %s\n", a1))
	f.WriteString(fmt.Sprintf("    MOV CX, %s\n", a2))
	f.WriteString("    IMUL CX\n")
	f.WriteString(fmt.Sprintf("    MOV %s, AX\n", res))
}

func (g *Generator) emitDiv(f *os.File, q quad.Quad) {
	a1, a2, res := asmVal(q.Arg1), asmVal(q.Arg2), q.Result

	f.WriteString(fmt.Sprintf("    MOV AX, %s\n", a1))
	f.WriteString(fmt.Sprintf("    MOV CX, %s\n", a2))
	f.WriteString("    CWD\n")
	f.WriteString("    IDIV CX\n")
	f.WriteString(fmt.Sprintf("    MOV %s, AX\n", res))
}

func (g *Generator) emitNeg(f *os.File, q quad.Quad) {
	f.WriteString(fmt.Sprintf("    MOV AX, %s\n", asmVal(q.Arg1)))
	f.WriteString("    NEG AX\n")
	f.WriteString(fmt.Sprintf("    MOV %s, AX\n", q.Result))
}

func (g *Generator) emitCmp(f *os.File, idx int, q quad.Quad) {
	op := q.Op
	a1, a2, res := asmVal(q.Arg1), asmVal(q.Arg2), q.Result
	skipLabel := fmt.Sprintf("S%d_end", idx)

	f.WriteString(fmt.Sprintf("    MOV AX, %s\n", a1))
	f.WriteString(fmt.Sprintf("    CMP AX, %s\n", a2))
	f.WriteString("    MOV AX, 0\n")

	jmpInsn := cmpJump(op)
	f.WriteString(fmt.Sprintf("    %s %s\n", jmpInsn, skipLabel))
	f.WriteString("    MOV AX, 1\n")
	f.WriteString(fmt.Sprintf("%s:\n", skipLabel))
	f.WriteString(fmt.Sprintf("    MOV %s, AX\n", res))
}

func (g *Generator) emitLogical(f *os.File, q quad.Quad, insn string) {
	a1, a2, res := asmVal(q.Arg1), asmVal(q.Arg2), q.Result

	f.WriteString(fmt.Sprintf("    MOV AX, %s\n", a1))
	f.WriteString(fmt.Sprintf("    %s AX, %s\n", insn, a2))
	f.WriteString(fmt.Sprintf("    MOV %s, AX\n", res))
}

func (g *Generator) emitNot(f *os.File, q quad.Quad) {
	f.WriteString(fmt.Sprintf("    MOV AX, %s\n", asmVal(q.Arg1)))
	f.WriteString("    XOR AX, 1\n")
	f.WriteString(fmt.Sprintf("    MOV %s, AX\n", q.Result))
}

func (g *Generator) emitBF(f *os.File, q quad.Quad) {
	f.WriteString(fmt.Sprintf("    MOV AX, %s\n", asmVal(q.Arg1)))
	f.WriteString("    CMP AX, 0\n")
	f.WriteString(fmt.Sprintf("    JZ L%s\n", q.Result))
}

func (g *Generator) emitArrayRead(f *os.File, q quad.Quad) {
	arr, idx, res := q.Arg1, asmVal(q.Arg2), q.Result

	f.WriteString(fmt.Sprintf("    MOV SI, %s\n", idx))
	f.WriteString("    ADD SI, SI\n")
	f.WriteString(fmt.Sprintf("    MOV AX, %s[SI]\n", arr))
	f.WriteString(fmt.Sprintf("    MOV %s, AX\n", res))
}

func (g *Generator) emitArrayWrite(f *os.File, q quad.Quad) {
	val, idx, arr := asmVal(q.Arg1), asmVal(q.Arg2), q.Result

	f.WriteString(fmt.Sprintf("    MOV SI, %s\n", idx))
	f.WriteString("    ADD SI, SI\n")
	f.WriteString(fmt.Sprintf("    MOV AX, %s\n", val))
	f.WriteString(fmt.Sprintf("    MOV %s[SI], AX\n", arr))
}

func (g *Generator) emitIn(f *os.File, q quad.Quad) {
	f.WriteString("    MOV AH, 01h\n")
	f.WriteString("    INT 21h\n")
	f.WriteString("    SUB AL, '0'\n")
	f.WriteString("    MOV AH, 0\n")
	if q.Result != "" {
		f.WriteString(fmt.Sprintf("    MOV %s, AX\n", q.Result))
	}
}

func (g *Generator) emitOut(f *os.File, q quad.Quad) {
	f.WriteString(fmt.Sprintf("    MOV AX, %s\n", asmVal(q.Arg1)))
	f.WriteString("    CALL PRINT_INT\n")
}

func (g *Generator) emitOutStr(f *os.File, q quad.Quad) {
	content, err := strconv.Unquote(q.Arg1)
	if err != nil {
		content = strings.Trim(q.Arg1, `"`)
	}
	label := g.strings[content]
	f.WriteString(fmt.Sprintf("    LEA DX, %s\n", label))
	f.WriteString("    MOV AH, 09h\n")
	f.WriteString("    INT 21h\n")
}

func (g *Generator) writePrintRoutine(f *os.File) {
	f.WriteString("; Print integer in AX\n")
	f.WriteString("PRINT_INT PROC\n")
	f.WriteString("    PUSH AX\n")
	f.WriteString("    PUSH BX\n")
	f.WriteString("    PUSH CX\n")
	f.WriteString("    PUSH DX\n")
	f.WriteString("\n")
	f.WriteString("    CMP AX, 0\n")
	f.WriteString("    JGE print_pos\n")
	f.WriteString("    PUSH AX\n")
	f.WriteString("    MOV DL, '-'\n")
	f.WriteString("    MOV AH, 02h\n")
	f.WriteString("    INT 21h\n")
	f.WriteString("    POP AX\n")
	f.WriteString("    NEG AX\n")
	f.WriteString("print_pos:\n")
	f.WriteString("    MOV CX, 0\n")
	f.WriteString("    MOV BX, 10\n")
	f.WriteString("print_div:\n")
	f.WriteString("    MOV DX, 0\n")
	f.WriteString("    DIV BX\n")
	f.WriteString("    PUSH DX\n")
	f.WriteString("    INC CX\n")
	f.WriteString("    CMP AX, 0\n")
	f.WriteString("    JNE print_div\n")
	f.WriteString("print_digits:\n")
	f.WriteString("    POP DX\n")
	f.WriteString("    ADD DL, '0'\n")
	f.WriteString("    MOV AH, 02h\n")
	f.WriteString("    INT 21h\n")
	f.WriteString("    LOOP print_digits\n")
	f.WriteString("\n")
	f.WriteString("    POP DX\n")
	f.WriteString("    POP CX\n")
	f.WriteString("    POP BX\n")
	f.WriteString("    POP AX\n")
	f.WriteString("    RET\n")
	f.WriteString("PRINT_INT ENDP\n")
}

// ---- helpers ----

func isTemp(s string) bool {
	if len(s) < 2 || s[0] != 'T' {
		return false
	}
	for _, c := range s[1:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func isNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func asmVal(s string) string {
	if s == "" {
		return ""
	}
	if isNum(s) {
		return integerize(s)
	}
	return s
}

func integerize(s string) string {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return "0"
	}
	return fmt.Sprintf("%d", int(f))
}

func cmpJump(op string) string {
	switch op {
	case ">":
		return "JLE"
	case "<":
		return "JGE"
	case ">=":
		return "JL"
	case "<=":
		return "JG"
	case "==":
		return "JNE"
	case "!=":
		return "JE"
	}
	return "JMP"
}

func sortTemps(t []string) {
	for i := 0; i < len(t); i++ {
		for j := i + 1; j < len(t); j++ {
			ni, _ := strconv.Atoi(t[i][1:])
			nj, _ := strconv.Atoi(t[j][1:])
			if nj < ni {
				t[i], t[j] = t[j], t[i]
			}
		}
	}
}

func sortStrings(keys []string, m map[string]string) {
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if m[keys[i]] > m[keys[j]] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
}
