package main

import (
	"fmt"

	"github.com/smtdfc/contractor/emiters/csharp"
	"github.com/smtdfc/contractor/emiters/golang"
	"github.com/smtdfc/contractor/emiters/java"
	"github.com/smtdfc/contractor/emiters/kotlin"
	"github.com/smtdfc/contractor/emiters/typescript"
	"github.com/smtdfc/contractor/exception"
	"github.com/smtdfc/contractor/generator"
	"github.com/smtdfc/contractor/parser"
)

func main() {
	fileName := "test.contract"

	code := `
		@CreateConstructor
		model Address<T>{
			province?: String
			district?: String
			village?: T
		}

		@CreateConstructor
		model User<T,U,K> {
			@IsEmail("error")
			name: String

			@IsModel("hello")
			address: Address<T>
			data: T
		}

		rest HelloW{
			path:"/get-products"
			method:"GET"
			requestBody: Address<String>
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

	tsEmitter := typescript.NewTypescriptEmitter()
	output, err = tsEmitter.Emit(ir)
	if err != nil {
		exception.PrintException(err, code)
		return
	}

	fmt.Println(output)

	javaEmitter := java.NewJavaEmitter()
	output, err = javaEmitter.Emit(ir)
	if err != nil {
		exception.PrintException(err, code)
		return
	}

	fmt.Println(output)

	kolinEmitter := kotlin.NewKotlinEmitter()
	output, err = kolinEmitter.Emit(ir)
	if err != nil {
		exception.PrintException(err, code)
		return
	}

	fmt.Println(output)

	csharpEmitter := csharp.NewCSharpEmitter()
	output, err = csharpEmitter.Emit(ir)
	if err != nil {
		exception.PrintException(err, code)
		return
	}

	fmt.Println(output)
}
