package main

import (
	"github.com/smtdfc/contractor/exception"
	"github.com/smtdfc/contractor/parser"
)

func main() {
	code := `
	model A{}
	model A{}
	`

	lexer := parser.NewLexer("test.contract")
	tokens, err := lexer.Start(code)
	if err != nil {
		exception.PrintException(err, code)
	}

	parser.PrintTokenList(tokens)
	p := parser.NewParser("test.contract", tokens)
	ast, err := p.Parse()
	if err != nil {
		exception.PrintException(err, code)
	} else {
		parser.PrintAST(ast, 1)
	}

	checker := parser.NewTypeChecker()
	err = checker.Check(ast)
	if err != nil {
		exception.PrintException(err, code)
	}
}
