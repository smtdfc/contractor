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
	"Array":  "ContractorRuntime.Array",
	"Map":    "map[string]interface{}",
}

type GoGenerator struct{}

func NewGoGenerator() *GoGenerator {
	return &GoGenerator{}
}

func (g *GoGenerator) ExtractValidationMetadata(node *parser.ModelFieldNode) helpers.FieldValidationMetadata {
	metadata := helpers.FieldValidationMetadata{
		Name:  node.Name,
		Rules: make([]helpers.FieldValidateRule, 0),
	}
	if node.Annotations == nil {
		return metadata
	}

	for _, anno := range node.Annotations.List {
		if validator, ok := parser.SupportedFieldValidatorAnnotations[anno.Name]; ok {
			rule := helpers.FieldValidateRule{RuleName: anno.Name}
			if len(anno.Args) > validator.Args-1 {
				rule.Message = anno.Args[validator.Args-1].(*parser.LiteralNode).Value
			}
			if validator.Args >= 2 && len(anno.Args) >= 2 {
				rule.Value = anno.Args[0].(*parser.LiteralNode).Value
			}
			metadata.Rules = append(metadata.Rules, rule)
		}
	}
	return metadata
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
	cb.WriteString(g.GenerateStaticValidate(node.Fields, modelName, modelType))
	return cb.String(), nil
}

func (g *GoGenerator) GenerateStaticValidate(fields []*parser.ModelFieldNode, modelName string, modelType string) string {
	cb := helpers.NewCodeBuffer(2)
	cb.Indent()

	cb.WriteLine("func Validate%s(obj *%s) error {", modelName, modelType)
	cb.Indent()
	cb.WriteLine("var errors map[string]string = make(map[string]string)")
	cb.WriteString("\n")
	for _, field := range fields {
		isOptional, _ := g.GetAnnotation(field.Annotations, "Optional")

		pName := parser.AnyToPascalCase(field.Name)

		if priv, _ := g.GetAnnotation(field.Annotations, "Private"); priv != nil {
			pName = parser.AnyToCamelCase(field.Name)
		}

		if isOptional == nil {
			cb.WriteLine("if (!ContractorRuntime.IsRequired(obj.%s)) {", pName)
			cb.Indent()
			cb.WriteLine("errors[\"%s\"] = \"%s is required\";", field.Name, field.Name)
			cb.Outdent()
			cb.WriteLine("}")
		}

		metadata := g.ExtractValidationMetadata(field)
		if len(metadata.Rules) > 0 {
			for _, rule := range metadata.Rules {
				valPart := ""
				if rule.Value != "" {
					valPart = ", " + rule.Value
				}

				cb.WriteLine("if (!ContractorRuntime.%s(obj.%s%s)) {", rule.RuleName, pName, valPart)
				cb.Indent()
				cb.WriteLine("errors[\"%s\"] = \"%s\";", field.Name, rule.Message)
				cb.Outdent()
				cb.WriteLine("}")
			}
		}
	}

	cb.WriteString("\n")
	cb.WriteLine("if (len(errors) > 0) {")
	cb.Indent()
	cb.WriteLine("return ContractorRuntime.NewValidationError(errors)")
	cb.Outdent()
	cb.WriteLine("}")
	cb.WriteLine("return nil")
	cb.Outdent()
	cb.WriteLine("}")
	return cb.String()
}

func (g *GoGenerator) Generate(ast *parser.AST, packageName string, goModulePath string) (string, parser.BaseError) {
	cb := helpers.NewCodeBuffer(2)
	cb.WriteLine("package %s", packageName)
	cb.WriteString("\n")
	cb.WriteLine("import ContractorRuntime \"%s\"", "github.com/smtdfc/contractor/runtime/go")
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
