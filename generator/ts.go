package generator

import (
	"fmt"

	"github.com/smtdfc/contractor/generator/helpers"
	"github.com/smtdfc/contractor/parser"
)

var TypescriptPrimitiveTypes = map[string]string{
	"String": "string",
	"Number": "number",
	"Bool":   "boolean",
	"Null":   "null",
	"Any":    "any",
}

var SupportedValidators = map[string]bool{
	"IsEmail":       true,
	"Max":           true,
	"Min":           true,
	"IsInt":         true,
	"IsFloat":       true,
	"IsBoolean":     true,
	"IsString":      true,
	"IsDateString":  true,
	"IsUUID":        true,
	"IsUrl":         true,
	"IsArray":       true,
	"MinLength":     true,
	"MaxLength":     true,
	"Length":        true,
	"IsNotEmpty":    true,
	"ArrayMinSize":  true,
	"ArrayMaxSize":  true,
	"ArrayLength":   true,
	"IsPhoneNumber": true,
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
	cb := helpers.NewCodeBuffer(2)
	cb.WriteLine("import {ContractorRuntime} from '@smtdfc/contractor';")
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
	cb := helpers.NewCodeBuffer(2)

	modelType := node.Name
	genericCode := ""
	if node.TypeVar != nil {
		genericCode, _ = g.GenerateType(node.TypeVar)
		modelType += genericCode
	}

	cb.WriteLine("export class %s%s {", node.Name, genericCode)
	cb.Indent()

	isGlobalData, _ := g.GetAnnotation(node.Annotations, "Data")
	isCreateConstructor, _ := g.GetAnnotation(node.Annotations, "CreateConstructor")
	isCreateMapper, _ := g.GetAnnotation(node.Annotations, "CreateMapper")
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
			optionalPart := ""
			if opt, _ := g.GetAnnotation(field.Annotations, "Optional"); opt != nil {
				optionalPart = "|undefined"
			}

			cb.WriteLine("public get%s(): %s%s {", pName, fType, optionalPart)
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

	if isCreateMapper != nil {
		cb.WriteString(g.GenerateMapper(node, modelType))
	}

	cb.WriteString(g.GenerateStaticValidate(node.Fields))

	cb.Outdent()
	cb.WriteLine("}")
	return cb.String(), nil
}

func (g *TypescriptGenerator) GenerateMapper(node *parser.ModelStatementNode, modelType string) string {
	cb := helpers.NewCodeBuffer(2)
	cb.Indent()

	typeVarName := "T"
	if node.TypeVar != nil {
		typeVarName = node.TypeVar.Name // Assuming node.TypeVar has a Name field
	}

	isCreateConstructor, _ := g.GetAnnotation(node.Annotations, "CreateConstructor")

	cb.WriteLine("public static fromObject<%s>(obj: any): %s<%s> {", typeVarName, node.Name, typeVarName)
	cb.Indent()

	if isCreateConstructor != nil {
		cb.WriteLine("const result = new (this as new (...args: any[]) => %s<%s>)(...([] as any));", node.Name, typeVarName)
	} else {
		cb.WriteLine("const result = new (this as new () => %s<%s>)();", node.Name, typeVarName)
	}
	cb.WriteString("\n")

	for _, field := range node.Fields {
		sourceKey := field.Name
		if mappingAnno, _ := g.GetAnnotation(field.Annotations, "Mapping"); mappingAnno != nil {
			if len(mappingAnno.Args) > 0 {
				if literal, ok := mappingAnno.Args[0].(*parser.LiteralNode); ok {
					sourceKey = literal.Value
				}
			}
		}

		cb.WriteLine("if (obj['%s'] != null) {", sourceKey)
		cb.Indent()
		cb.WriteLine("result.%s = obj['%s'];", field.Name, sourceKey)
		cb.Outdent()

		// if defAnno, _ := g.GetAnnotation(field.Annotations, "Default"); defAnno != nil {
		// 	if len(defAnno.Args) > 0 {
		// 		if literal, ok := defAnno.Args[0].(*parser.LiteralNode); ok {
		// 			cb.WriteLine("} else {")
		// 			cb.Indent()
		// 			val := literal.Value
		// 			if _, isPrimitive := TypescriptPrimitiveTypes[literal.Value]; !isPrimitive && literal.Type == parser.TOKEN_STRING {
		// 				val = fmt.Sprintf("'%s'", val)
		// 			}
		// 			cb.WriteLine("result.%s = %s;", field.Name, val)
		// 			cb.Outdent()
		// 		}
		// 	}
		// }

		cb.WriteLine("}")
	}

	cb.WriteLine("return result;")
	cb.Outdent()
	cb.WriteLine("}")
	cb.WriteString("\n")

	cb.WriteLine("public toObject(): Record<string, any> {")
	cb.Indent()
	cb.WriteLine("const result: Record<string, any> = {};")
	cb.WriteString("\n")

	for _, field := range node.Fields {
		targetKey := field.Name
		if mappingAnno, _ := g.GetAnnotation(field.Annotations, "Mapping"); mappingAnno != nil {
			if len(mappingAnno.Args) > 0 {
				if literal, ok := mappingAnno.Args[0].(*parser.LiteralNode); ok {
					targetKey = literal.Value
				}
			}
		}

		cb.WriteLine("if (this.%s !== undefined) {", field.Name)
		cb.Indent()

		cb.WriteLine("const val = this.%s as any;", field.Name)
		cb.WriteLine("if (val?.toObject && typeof val.toObject === 'function') {")
		cb.Indent()
		cb.WriteLine("result['%s'] = val.toObject();", targetKey)
		cb.Outdent()
		cb.WriteLine("} else if (Array.isArray(val)) {")
		cb.Indent()
		cb.WriteLine("result['%s'] = val.map((item: any) => item?.toObject?.() ?? item);", targetKey)
		cb.Outdent()
		cb.WriteLine("} else {")
		cb.Indent()
		cb.WriteLine("result['%s'] = val;", targetKey)
		cb.Outdent()
		cb.WriteLine("}")

		cb.Outdent()
		cb.WriteLine("}")
	}
	cb.WriteString("\n")
	cb.WriteLine("return result;")
	cb.Outdent()
	cb.WriteLine("}")

	return cb.String()
}
func (g *TypescriptGenerator) GenerateStaticValidate(fields []*parser.ModelFieldNode) string {
	cb := helpers.NewCodeBuffer(2)
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
