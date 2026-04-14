package main

import (
	"fmt"
	"go_compiler/internal/ast"
	"go_compiler/internal/lexer"
	"go_compiler/internal/semantic"
)

func main() {
	src := `
BeginProject testprog;

Setup :
  define x | y : integer;
  define moyenne : float = 0.0;
  define i : integer;
  define Tab : [integer ;20] ;
  const Pi : float = 3.14159;

Run :
{
  x <- 5;
  y <- x + 10;
  Tab[0] <- y * 2;
  moyenne <- (x + y) / 2;

  if (x > y) then:
  {
    x <- x - 1;
  } else {
    y <- y + 1;
  } endIf;

  loop while (x < 10)
  {
    x <- x + 1;
  } endloop;

  for i in 1 to 10 {
    y <- y * i;
  } endfor;

  in(x);
  out("Resultat: ", moyenne);
}

EndProject ;
`

	// ── Analyse Lexicale ─────────────────────────────────────
	fmt.Println("══════════════════════════════════════════════")
	fmt.Println("          ANALYSE LEXICALE — ProLang          ")
	fmt.Println("══════════════════════════════════════════════")

	lx := lexer.New(src)
	lx.Lex()

	for _, tok := range lx.Tokens {
		fmt.Printf("Token{%-15s %-22q ligne:%-3d col:%d}\n", tok.Kind, tok.Value, tok.Line, tok.Col)
	}

	fmt.Println()
	if len(lx.Errors) > 0 {
		fmt.Println("── Erreurs Lexicales ──────────────────────────")
		for _, e := range lx.Errors {
			fmt.Println(" ✗", e)
		}
	} else {
		fmt.Println(" ✓ Aucune erreur lexicale détectée.")
	}

	// ── Analyse Sémantique ───────────────────────────────────
	fmt.Println("\n══════════════════════════════════════════════")
	fmt.Println("         ANALYSE SÉMANTIQUE — ProLang         ")
	fmt.Println("══════════════════════════════════════════════")

	table := ast.NewSymbolTable()
	analyzer := semantic.NewAnalyzer(lx.Tokens, table)
	analyzer.Analyze()

	table.Print()

	if len(analyzer.Errors) > 0 {
		fmt.Println("── Erreurs Sémantiques ────────────────────────")
		for _, e := range analyzer.Errors {
			fmt.Println(" ✗", e)
		}
	} else {
		fmt.Println(" ✓ Aucune erreur sémantique détectée.")
	}
}
