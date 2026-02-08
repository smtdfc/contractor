package generator

import (
	"fmt"
	"strings"

	"github.com/smtdfc/contractor/parser"
)

type CodeBuffer struct {
	strings.Builder
	indentSize  int
	indentLevel int
}

func NewCodeBuffer(indentSize int) *CodeBuffer {
	return &CodeBuffer{indentSize: indentSize}
}

func (cb *CodeBuffer) Indent() { cb.indentLevel++ }
func (cb *CodeBuffer) Outdent() {
	if cb.indentLevel > 0 {
		cb.indentLevel--
	}
}

func (cb *CodeBuffer) WriteLine(format string, args ...interface{}) {
	indent := strings.Repeat(" ", cb.indentLevel*cb.indentSize)
	cb.WriteString(indent)
	fmt.Fprintf(&cb.Builder, format, args...)
	cb.WriteString("\n")
}

var TypescriptPrimitiveTypes = map[string]string{
	"String": "string",
	"Number": "number",
	"Bool":   "boolean",
	"Null":   "null",
	"Any":    "any",
}

var SupportedValidators = map[string]bool{
	"IsEmail":    true,
	"Max":        true,
	"Min":        true,
	"Length":     true,
	"IsPattern":  true,
	"IsNotEmpty": true,
}

type FieldValidateRule struct {
	RuleName string
	Value    string
	Message  string
}

type FieldValidationMetadata struct {
	Name  string
	Rules []FieldValidateRule
}

type TypescriptGenerator struct{}

func NewTypescriptGenerator() *TypescriptGenerator {
	return &TypescriptGenerator{}
}

func (g *TypescriptGenerator) GetAnnotation(node *parser.AnnotationChainNode, name string) (*parser.AnnotationNode, parser.BaseError) {
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

func (g *TypescriptGenerator) GenerateType(node parser.Node) (string, parser.BaseError) {
	switch v := node.(type) {
	case *parser.TypeVarNode:
		return fmt.Sprintf("<%s>", v.Name), nil
	case *parser.TypeDeclarationNode:
		name := v.Name
		if tsType, ok := TypescriptPrimitiveTypes[name]; ok {
			name = tsType
		}
		if v.Generic != nil {
			gen, err := g.GenerateType(v.Generic)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%s<%s>", name, gen), nil
		}
		return name, nil
	}
	return "", nil
}

func (g *TypescriptGenerator) ExtractValidationMetadata(node *parser.ModelFieldNode) FieldValidationMetadata {
	metadata := FieldValidationMetadata{
		Name:  node.Name,
		Rules: make([]FieldValidateRule, 0),
	}
	if node.Annotations == nil {
		return metadata
	}

	for _, anno := range node.Annotations.List {
		if SupportedValidators[anno.Name] {
			rule := FieldValidateRule{RuleName: anno.Name}
			if anno.Name == "IsEmail" || anno.Name == "IsNotEmpty" {
				if len(anno.Args) > 0 {
					rule.Message = anno.Args[0].(*parser.LiteralNode).Value
					metadata.Rules = append(metadata.Rules, rule)
				}
			} else if len(anno.Args) >= 2 {
				rule.Value = anno.Args[0].(*parser.LiteralNode).Value
				rule.Message = anno.Args[1].(*parser.LiteralNode).Value
				metadata.Rules = append(metadata.Rules, rule)
			}
		}
	}
	return metadata
}

func (g *TypescriptGenerator) Generate(ast *parser.AST) (string, parser.BaseError) {
	cb := NewCodeBuffer(2)
	cb.WriteLine("import {ContractorRuntime} from 'contractor';")
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

func (g *TypescriptGenerator) GenerateModel(node *parser.ModelStatementNode) (string, parser.BaseError) {
	cb := NewCodeBuffer(2)

	genericCode := ""
	if node.TypeVar != nil {
		genericCode, _ = g.GenerateType(node.TypeVar)
	}

	cb.WriteLine("export class %s%s {", node.Name, genericCode)
	cb.Indent()

	isGlobalData, _ := g.GetAnnotation(node.Annotations, "Data")
	isCreateConstructor, _ := g.GetAnnotation(node.Annotations, "CreateConstructor")

	for _, field := range node.Fields {
		fType, _ := g.GenerateType(field.Type)
		access := "public"
		optional := ""

		if priv, _ := g.GetAnnotation(field.Annotations, "Private"); priv != nil {
			access = "private"
		}
		if opt, _ := g.GetAnnotation(field.Annotations, "Optional"); opt != nil {
			optional = "?"
		}

		cb.WriteLine("%s %s%s: %s;", access, field.Name, optional, fType)
	}
	cb.WriteString("\n")

	if isCreateConstructor != nil {
		cb.WriteLine("constructor(")
		cb.Indent()
		for i, field := range node.Fields {
			fType, _ := g.GenerateType(field.Type)
			suffix := ","
			if i == len(node.Fields)-1 {
				suffix = ""
			}
			cb.WriteLine("%s: %s%s", field.Name, fType, suffix)
		}
		cb.Outdent()
		cb.WriteLine(") {")
		cb.Indent()
		for _, field := range node.Fields {
			cb.WriteLine("this.%s = %s;", field.Name, field.Name)
		}
		cb.Outdent()
		cb.WriteLine("}")
		cb.WriteString("\n")
	}

	for _, field := range node.Fields {
		fType, _ := g.GenerateType(field.Type)
		pName := parser.AnyToPascalCase(field.Name)

		hasGetter, _ := g.GetAnnotation(field.Annotations, "Getter")
		hasSetter, _ := g.GetAnnotation(field.Annotations, "Setter")

		if isGlobalData != nil || hasGetter != nil {
			cb.WriteLine("public get%s(): %s {", pName, fType)
			cb.Indent()
			cb.WriteLine("return this.%s;", field.Name)
			cb.Outdent()
			cb.WriteLine("}")
			cb.WriteString("\n")
		}
		if isGlobalData != nil || hasSetter != nil {
			cb.WriteLine("public set%s(v: %s): void {", pName, fType)
			cb.Indent()
			cb.WriteLine("this.%s = v;", field.Name)
			cb.Outdent()
			cb.WriteLine("}")
			cb.WriteString("\n")
		}
	}

	cb.WriteString(g.GenerateStaticValidate(node.Fields))

	cb.Outdent()
	cb.WriteLine("}")
	return cb.String(), nil
}

func (g *TypescriptGenerator) GenerateStaticValidate(fields []*parser.ModelFieldNode) string {
	cb := NewCodeBuffer(2)
	cb.Indent()

	cb.WriteLine("public static validate(obj: any): ContractorRuntime.ValidationError | null {")
	cb.Indent()
	cb.WriteLine("const errors: Record<string, string> = {};")
	cb.WriteString("\n")

	for _, field := range fields {
		isOptional, _ := g.GetAnnotation(field.Annotations, "Optional")

		if isOptional == nil {
			cb.WriteLine("if (!ContractorRuntime.Validators.IsRequired(obj.%s)) {", field.Name)
			cb.Indent()
			cb.WriteLine("errors['%s'] = '%s is required';", field.Name, field.Name)
			cb.Outdent()
			cb.WriteLine("}")
		}

		metadata := g.ExtractValidationMetadata(field)
		if len(metadata.Rules) > 0 {
			cb.WriteLine("if (obj.%s !== undefined && obj.%s !== null) {", field.Name, field.Name)
			cb.Indent()
			for _, rule := range metadata.Rules {
				valPart := ""
				if rule.Value != "" {
					valPart = ", " + rule.Value
				}

				cb.WriteLine("if (!ContractorRuntime.Validators.%s(obj.%s%s)) {", rule.RuleName, field.Name, valPart)
				cb.Indent()
				cb.WriteLine("errors['%s'] = '%s';", field.Name, rule.Message)
				cb.Outdent()
				cb.WriteLine("}")
			}
			cb.Outdent()
			cb.WriteLine("}")
		}
	}

	cb.WriteString("\n")
	cb.WriteLine("if (Object.keys(errors).length > 0) {")
	cb.Indent()
	cb.WriteLine("return new ContractorRuntime.ValidationError(errors);")
	cb.Outdent()
	cb.WriteLine("}")
	cb.WriteLine("return null;")

	cb.Outdent()
	cb.WriteLine("}")

	return cb.String()
}
