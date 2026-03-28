package main

import (
	"fmt"

	"github.com/smtdfc/contractor/emiters/golang"
	"github.com/smtdfc/contractor/exception"
	"github.com/smtdfc/contractor/generator"
	"github.com/smtdfc/contractor/parser"
)

func main() {
	fileName := "test.contract"

	code := `
		model Address{
			province?: String
			district?: String
			village?: String
		}


		model User {
			name: String
			address: Address
		}
	`

	lexer := parser.NewLexer(fileName)

	tokens, err := lexer.Start(code)
	if err != nil {
		exception.PrintException(err, code)
		return
	}

	p := parser.NewParser(fileName, tokens)
	ast, err := p.Parse()
	if err != nil {
		exception.PrintException(err, code)
		return
	}

	typeChecker := parser.NewTypeChecker()
	err = typeChecker.Check(ast)
	if err != nil {
		exception.PrintException(err, code)
		return
	}

	irGenerator := generator.NewIRGenerator()
	ir, err := irGenerator.GenerateProgram(ast)
	if err != nil {
		exception.PrintException(err, code)
		return
	}

	goEmitter := golang.NewGoEmitter()
	output, err := goEmitter.Emit(ir)
	if err != nil {
		exception.PrintException(err, code)
		return
	}

	fmt.Println(output)
}
