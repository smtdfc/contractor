package golang

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/smtdfc/contractor/exception"
	"github.com/smtdfc/contractor/generator"
)

type GoEmitter struct{}

func (e *GoEmitter) EmitModelField(ir *generator.ModelField) (string, exception.IException) {

}

func (e *GoEmitter) EmitModel(ir *generator.ModelIR) (string, exception.IException) {
	var sb strings.Builder

	data := map[string]string{
		"Name": ir.Name,
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
