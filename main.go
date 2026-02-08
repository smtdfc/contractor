package main

import (
	"fmt"
	"strings"

	g "github.com/smtdfc/contractor/generator"
	p "github.com/smtdfc/contractor/parser"
)

func printError(err p.BaseError, code string) {
	if err == nil {
		return
	}

	loc := err.GetLocation()
	fmt.Printf("\033[1;31m[%s]\033[0m: %s\n", err.GetName(), err.GetMessage())

	if loc != nil {
		fmt.Printf("  \033[1;33m->\033[0m %s:%d:%d\n", loc.File, loc.Start.Line, loc.Start.Column)
		fmt.Println("   |")

		lines := strings.Split(code, "\n")
		lineIdx := loc.Start.Line - 1

		if lineIdx >= 0 && lineIdx < len(lines) {
			rawLine := lines[lineIdx]
			displayLine := strings.ReplaceAll(rawLine, "\t", "    ")

			fmt.Printf("%2d |  %s\n", loc.Start.Line, displayLine)
			padding := ""
			//tabCount := 0
			for i := 0; i < loc.Start.Column-1 && i < len(rawLine); i++ {
				if rawLine[i] == '\t' {
					padding += "    "
				} else {
					padding += " "
				}
			}

			length := 1
			if loc.End.Line == loc.Start.Line {
				length = loc.End.Column - loc.Start.Column
			}
			if length <= 0 {
				length = 1
			}

			underline := strings.Repeat("^", length)
			fmt.Printf("   |  %s\033[1;31m%s\033[0m\n", padding, underline)
		}
		fmt.Println("   |")
	}
	fmt.Println()
}

func main() {
	code := `
	@CreateConstructor
	@Data
	model LoginDTO{
		@Private
		@Optional
		@IsEmail("shssjs")
		String email
		
		@Required
		String password
	}
	
	@Data
	model Response<T>{
		Array<T> Data
		T Hello
	}
	`
	lexer := p.NewLexer()
	tokens, err := lexer.Parse(code, "<test>")

	if err != nil {
		printError(err, code)
	}

	for _, tok := range tokens {
		fmt.Println(tok.Loc)
		fmt.Printf("Token: %s |Loc: %s\n", tok.ToString(), tok.Loc.ToString())
	}

	parser := p.NewParser()
	ast, err := parser.Parse(tokens, "h")
	if err != nil {
		printError(err, code)
	}
	p.PrintAST(ast, 1)

	typeChecker := p.NewTypeChecker()
	err = typeChecker.Check(ast)
	if err != nil {
		printError(err, code)
	}

	tsGenerator := g.NewTypescriptGenerator()
	code, err = tsGenerator.Generate(ast)
	if err != nil {
		printError(err, code)
	}

	fmt.Println(code)
}
