package generator

import (
	"github.com/smtdfc/contractor/generator/helpers"
	"github.com/smtdfc/contractor/parser"
)

type GoGenerator struct{}

func NewGoGenerator() *GoGenerator {
	return &GoGenerator{}
}

func (g *GoGenerator) Generate(ast *parser.AST, packageName string, goModulePath string) (string, parser.BaseError) {
	cb := helpers.NewCodeBuffer(2)
	cb.WriteLine("package %s", packageName)
	cb.WriteString("\n")

	return cb.String(), nil
}
