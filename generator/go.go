package generator

import (
	"fmt"

	"github.com/smtdfc/contractor/generator/helpers"
	"github.com/smtdfc/contractor/parser"
)

var GoPrimitiveTypes = map[string]string{
	"String": "string",
	"Number": "float64",
	"Bool":   "bool",
	"Null":   "nil",
	"Any":    "interface{}",
	"Array":  "Contractor.Array",
	"Map":    "map[string]interface{}",
}

type GoGenerator struct{}

func NewGoGenerator() *GoGenerator {
	return &GoGenerator{}
}
func (g *GoGenerator) GenerateType(node parser.Node, constraint bool) (string, parser.BaseError) {
	switch v := node.(type) {
	case *parser.TypeVarNode:
		if constraint {
			return fmt.Sprintf("[%s any]", v.Name), nil
		}

		return fmt.Sprintf("[%s]", v.Name), nil
	case *parser.TypeDeclarationNode:
		name := v.Name
		if goType, ok := GoPrimitiveTypes[name]; ok {
			name = goType
		}
		if v.Generic != nil {
			gen, err := g.GenerateType(v.Generic, false)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%s[%s]", name, gen), nil
		}
		return name, nil
	}
	return "", nil
}

func (g *GoGenerator) GetAnnotation(node *parser.AnnotationChainNode, name string) (*parser.AnnotationNode, parser.BaseError) {
	if node == nil {
		return nil, nil
	}
	for _, anno := range node.List {
		if anno.Name == name {
			return anno, nil
		}
	}
	return nil, nil
}

func (g *GoGenerator) GenerateModel(node *parser.ModelStatementNode) (string, parser.BaseError) {
	cb := helpers.NewCodeBuffer(2)
	// for _, field := range node.Fields {
	// 	fieldName := parser.AnyToPascalCase(field.Name)
	// 	fieldType := helpers.MapTypeToGo(field.Type)
	// 	cb.WriteLine("%s %s `json:\"%s\"`", fieldName, fieldType, field.Name)
	// }
	// cb.WriteLine("}")
	// cb.WriteString("\n")
	modelName := parser.AnyToPascalCase(node.Name)
	modelType := node.Name
	genericCode := ""
	if node.TypeVar != nil {
		genericCode, _ = g.GenerateType(node.TypeVar, false)
		modelType += genericCode

		genericCode, _ = g.GenerateType(node.TypeVar, true)
		modelName += genericCode
	}

	cb.WriteLine("type %s struct {", modelName)
	cb.Indent()
	isGlobalData, _ := g.GetAnnotation(node.Annotations, "Data")
	isCreateConstructor, _ := g.GetAnnotation(node.Annotations, "CreateConstructor")
	// isCreateMapper, _ := g.GetAnnotation(node.Annotations, "CreateMapper")

	for _, field := range node.Fields {
		fType, _ := g.GenerateType(field.Type, false)
		structFieldName := parser.AnyToPascalCase(field.Name)
		optional := ""

		if priv, _ := g.GetAnnotation(field.Annotations, "Private"); priv != nil {
			structFieldName = parser.AnyToCamelCase(field.Name)
		}

		if opt, _ := g.GetAnnotation(field.Annotations, "Optional"); opt != nil {
			optional = "*"
		}

		cb.WriteLine("%s %s%s `json:\"%s,omitempty\"`", structFieldName, optional, fType, field.Name)
	}

	cb.WriteString("\n")
	cb.Outdent()
	cb.WriteLine("}")

	if isCreateConstructor != nil {
		cb.WriteLine("func New%s(", modelName)
		cb.Indent()
		for _, field := range node.Fields {
			fType, _ := g.GenerateType(field.Type, false)

			optionalPart := ""
			if opt, _ := g.GetAnnotation(field.Annotations, "Optional"); opt != nil {
				optionalPart = "*"
			}

			cb.WriteLine("%s %s%s,", field.Name, optionalPart, fType)
		}
		cb.Outdent()
		cb.WriteLine(")*%s {", modelType)
		cb.Indent()
		cb.WriteLine("return &%s{", modelType)
		for _, field := range node.Fields {
			cb.WriteLine("%s: %s,", parser.AnyToPascalCase(field.Name), field.Name)
		}
		cb.WriteLine("}")
		cb.Outdent()
		cb.WriteLine("}")
		cb.WriteString("\n")
	}

	for _, field := range node.Fields {
		fType, _ := g.GenerateType(field.Type, false)
		pName := parser.AnyToPascalCase(field.Name)

		if priv, _ := g.GetAnnotation(field.Annotations, "Private"); priv != nil {
			pName = parser.AnyToCamelCase(field.Name)
		}

		hasGetter, _ := g.GetAnnotation(field.Annotations, "Getter")
		hasSetter, _ := g.GetAnnotation(field.Annotations, "Setter")
		optionalPart := ""
		if opt, _ := g.GetAnnotation(field.Annotations, "Optional"); opt != nil {
			optionalPart = "*"
		}

		if isGlobalData != nil || hasGetter != nil {

			cb.WriteLine("func (i *%s) Get%s() %s%s{", modelType, pName, optionalPart, fType)
			cb.Indent()
			cb.WriteLine("return i.%s", pName)
			cb.Outdent()
			cb.WriteLine("}")
			cb.WriteString("\n")
		}

		if isGlobalData != nil || hasSetter != nil {
			cb.WriteLine("func (i *%s) Set%s(v %s%s) {", modelType, pName, optionalPart, fType)
			cb.Indent()
			cb.WriteLine("i.%s = v", pName)
			cb.Outdent()
			cb.WriteLine("}")
			cb.WriteString("\n")
		}

	}

	cb.WriteString("\n")
	return cb.String(), nil
}

func (g *GoGenerator) Generate(ast *parser.AST, packageName string, goModulePath string) (string, parser.BaseError) {
	cb := helpers.NewCodeBuffer(2)
	cb.WriteLine("package %s", packageName)
	cb.WriteString("\n")

	for _, stat := range ast.Statements {
		if v, ok := stat.(*parser.ModelStatementNode); ok {
			res, err := g.GenerateModel(v)
			if err != nil {
				return "", err
			}
			cb.WriteString(res)
		}
	}

	return cb.String(), nil
}
