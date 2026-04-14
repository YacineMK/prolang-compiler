package main

import (
	"fmt"
	"go_compiler/internal/lexer"
)

func main() {
	src := `
BeginProject testprog;

Setup :
  define x | y : integer;
  define moyenne : float = 0.0;
  const Pi : float = 3.14159;

Run :
{
  x <- 5;
  y <- x + 10;
  out("Resultat:", moyenne);
}

EndProject;
`

	lx := lexer.New(src)
	lx.Lex()

	fmt.Println("════ TOKENS ════")
	for _, t := range lx.Tokens {
		fmt.Printf("%+v\n", t)
	}

	fmt.Println("\n════ ERRORS ════")
	if len(lx.Errors) == 0 {
		fmt.Println("No lexical errors")
	} else {
		for _, e := range lx.Errors {
			fmt.Println(e)
		}
	}

}