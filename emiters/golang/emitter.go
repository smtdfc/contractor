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

func (e *GoEmitter) EmitModelField(ir *generator.ModelField) (string, exception.IException) {
	var sb strings.Builder
	jsonTag := fmt.Sprintf(`json:"%s"`, helpers.ToCamelCase(ir.Name))
	data := map[string]string{
		"Name": helpers.ToPascalCase(ir.Name),
		"Type": "any",
		"Tag":  jsonTag,
	}

	tmpl, _ := template.New("test").Parse(ModelFieldTemplate)

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
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
