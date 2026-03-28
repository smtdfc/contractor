package golang

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/smtdfc/contractor/exception"
	"github.com/smtdfc/contractor/generator"
	"github.com/smtdfc/contractor/internal/helpers"
)

type GoEmitter struct{}

func (e *GoEmitter) EmitBuildInType(ir *generator.TypeIR) (string, exception.IException) {

	switch ir.Name {
	case "Int":
		return "int", nil
	case "Float":
		return "float", nil
	case "Bool":
		return "bool", nil
	case "Any":
		return "any", nil
	case "String":
		return "string", nil
	case "Array":
		genericType, err := e.EmitType(ir.Generics[0])
		if err != nil {
			return "", err
		}

		return fmt.Sprintf(`[]%s`, genericType), nil
	}
	return "any", nil
}

func (e *GoEmitter) EmitType(ir *generator.TypeIR) (string, exception.IException) {

	if ir.Kind == generator.TypeKindBuiltin {
		return e.EmitBuildInType(ir)
	}

	return "any", nil
}

func (e *GoEmitter) EmitModelField(ir *generator.ModelField) (string, exception.IException) {
	var sb strings.Builder
	jsonTag := fmt.Sprintf(`json:"%s"`, helpers.ToCamelCase(ir.Name))

	typeStr, err := e.EmitType(ir.Type)
	if err != nil {
		return "", err
	}

	if ir.IsOptional {
		typeStr = fmt.Sprintf("*%s", typeStr)
	}

	data := map[string]string{
		"Name": helpers.ToPascalCase(ir.Name),
		"Type": typeStr,
		"Tag":  jsonTag,
	}

	tmpl, _ := template.New("test").Parse(ModelFieldTemplate)

	var tpl bytes.Buffer
	err_ := tmpl.Execute(&tpl, data)
	if err_ != nil {
		return "", exception.NewEmitException("Error when emit go code", ir.Span.ToLocation())
	}

	sb.WriteString(tpl.String())
	return sb.String(), nil
}

func (e *GoEmitter) EmitModel(ir *generator.ModelIR) (string, exception.IException) {
	var sb strings.Builder
	var fields = []string{}

	for _, field := range ir.Fields {
		c, err := e.EmitModelField(field)
		if err != nil {
			return "", err
		}

		fields = append(fields, c)
	}

	data := map[string]any{
		"Name":   ir.Name,
		"Fields": fields,
	}

	tmpl, _ := template.New("test").Parse(ModelTemplate)

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit go code", ir.Span.ToLocation())
	}

	sb.WriteString(tpl.String())
	return sb.String(), nil
}

func (e *GoEmitter) Emit(ir *generator.ProgramIR) (string, exception.IException) {
	var sb strings.Builder

	for _, model := range ir.Models {
		code, err := e.EmitModel(model)
		if err != nil {
			return "", err
		}
		sb.WriteString(code)
	}

	return sb.String(), nil
}

func NewGoEmitter() *GoEmitter {
	return &GoEmitter{}
}
