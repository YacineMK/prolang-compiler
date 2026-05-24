package main

import (
	"go_compiler/internal/asm"
	"go_compiler/internal/ast"
	"go_compiler/internal/lexer"
	"go_compiler/internal/optimizer"
	"go_compiler/internal/quad"
	"go_compiler/internal/semantic"
	"go_compiler/internal/utils"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		utils.PrintError("Usage: go_compiler <file_path>")
		utils.PrintError("Example: go_compiler example/test.pl")
		os.Exit(1)
	}

	filePath := os.Args[1]
	utils.PrintInfo("Input file: %s", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		utils.PrintError("Error reading file %s: %v", filePath, err)
		os.Exit(1)
	}

	src := string(data)

	utils.PrintSection("ANALYSE LEXICALE — ProLang")
	utils.PrintPhase("Lexical Analysis")

	lx := lexer.New(src)
	lx.Lex()

	utils.PrintInfo("Total tokens found: %d", len(lx.Tokens))
	utils.PrintDebug("Tokens details:")
	for _, tok := range lx.Tokens {
		utils.PrintDebug("  Token{%-15s %-22q ligne:%-3d col:%d}", tok.Kind, tok.Value, tok.Line, tok.Col)
	}

	if len(lx.Errors) > 0 {
		utils.PrintWarning("Found %d lexical errors", len(lx.Errors))
		for _, e := range lx.Errors {
			utils.PrintError("%s", e)
		}
	} else {
		utils.PrintSuccess("Lexical Analysis successful")
	}

	utils.PrintSection("ANALYSE SÉMANTIQUE — ProLang")
	utils.PrintPhase("Semantic Analysis")

	table := ast.NewSymbolTable()
	analyzer := semantic.NewAnalyzer(lx.Tokens, table)
	analyzer.Analyze()

	table.Print()

	if len(analyzer.Errors) > 0 {
		utils.PrintWarning("Found %d semantic errors", len(analyzer.Errors))
		for _, e := range analyzer.Errors {
			utils.PrintError("%s", e)
		}
		utils.PrintError("Compilation failed due to semantic errors")
		os.Exit(1)
	} else {
		utils.PrintSuccess("Semantic Analysis successful")
	}

	utils.PrintSection("GÉNÉRATION DES QUADS — ProLang")
	utils.PrintPhase("Quad Generation")

	qg := semantic.NewQuadGenerator(analyzer)
	qm := qg.Generate()

	utils.PrintInfo("Quad generation from Run block statements")

	quad.Print("── QUADS GÉNÉRÉS ──", qm.Quads)
	utils.PrintStats("Total quads generated: %d", len(qm.Quads))

	utils.PrintSection("OPTIMISATION — ProLang")
	utils.PrintPhase("Code Optimization")

	optimizedQuads := optimizer.Optimize(qm.Quads)

	quad.Print("── QUADS OPTIMISÉS ──", optimizedQuads)
	utils.PrintStats("Optimization: %d before → %d after", len(qm.Quads), len(optimizedQuads))

	utils.PrintSection("GÉNÉRATION ASSEMBLY — ProLang")
	utils.PrintPhase("8086 Assembly Generation")

	os.MkdirAll("out", 0755)
	base := filepath.Base(filePath)
	asmPath := filepath.Join("out", base+".asm")
	gen := asm.NewGenerator(optimizedQuads, table)
	if err := gen.Generate(asmPath); err != nil {
		utils.PrintError("Assembly generation failed: %v", err)
		os.Exit(1)
	}
	utils.PrintSuccess("Assembly generated: %s", asmPath)

	utils.PrintSection("COMPILATION RÉUSSIE — ProLang")
	utils.PrintSuccess("Compilation successful for: %s", filePath)
}
