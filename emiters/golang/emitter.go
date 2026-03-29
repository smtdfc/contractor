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
		genericType, err := e.EmitType(ir.Generics[0], false)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf(`[]%s`, genericType), nil
	}
	return "any", nil
}

func (e *GoEmitter) EmitGenericType(ir *generator.TypeIR) (string, exception.IException) {
	return ir.Name, nil
}

func (e *GoEmitter) EmitModelType(ir *generator.TypeIR) (string, exception.IException) {
	genericTypes := []string{}

	for _, generic := range ir.Generics {
		genericType, err := e.EmitType(generic, false)
		if err != nil {
			return "", err
		}

		genericTypes = append(genericTypes, genericType)
	}

	if len(genericTypes) > 0 {
		return fmt.Sprintf(`*%s[%s]`, ir.Name, strings.Join(genericTypes, ",")), nil
	}

	return fmt.Sprintf(`*%s`, ir.Name), nil
}

func (e *GoEmitter) EmitType(ir *generator.TypeIR, isOptional bool) (string, exception.IException) {

	if ir.Kind == generator.TypeKindBuiltin {
		t, err := e.EmitBuildInType(ir)
		if isOptional {
			t = fmt.Sprintf("*%s", t)
		}

		return t, err
	}

	if ir.Kind == generator.TypeKindModel {
		return e.EmitModelType(ir)
	}

	if ir.Kind == generator.TypeKindGeneric {
		return e.EmitGenericType(ir)
	}

	return "any", nil
}

func (e *GoEmitter) EmitModelField(ir *generator.ModelField) (string, exception.IException) {
	var sb strings.Builder
	jsonTag := fmt.Sprintf(`json:"%s"`, helpers.ToCamelCase(ir.Name))

	typeStr, err := e.EmitType(ir.Type, ir.IsOptional)
	if err != nil {
		return "", err
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

func (e *GoEmitter) EmitTypeParams(params []string, constraint bool) string {
	if len(params) == 0 {
		return ""
	}

	if constraint {
		return fmt.Sprintf(`[%s any]`, strings.Join(params, ","))
	}

	return fmt.Sprintf(`[%s]`, strings.Join(params, ","))
}

func (e *GoEmitter) EmitValue(ir *generator.ValueIR) (string, exception.IException) {

	if ir.Kind == "Number" {
		return ir.Value.(string), nil
	}

	if ir.Kind == "String" {
		return fmt.Sprintf(`"%s"`, ir.Value.(string)), nil
	}

	if ir.Kind == "Bool" {
		return ir.Value.(string), nil
	}

	if ir.Kind == "Array" {
		elements := []string{}
		for _, ele := range ir.Value.([]*generator.ValueIR) {
			value, err := e.EmitValue(ele)
			if err != nil {
				return "", err
			}

			elements = append(elements, value)
		}

		return fmt.Sprintf(`[%s]`, strings.Join(elements, ",")), nil
	}

	return "", nil
}

func (e *GoEmitter) EmitCreateConstructor(ir *generator.ModelIR) (string, exception.IException) {
	args := make([]string, 0, len(ir.Fields))
	assignments := make([]string, 0, len(ir.Fields))

	for _, field := range ir.Fields {
		typeStr, err := e.EmitType(field.Type, field.IsOptional)
		if err != nil {
			return "", err
		}

		argName := helpers.ToCamelCase(field.Name)
		fieldName := helpers.ToPascalCase(field.Name)
		args = append(args, fmt.Sprintf("%s %s", argName, typeStr))
		assignments = append(assignments, fmt.Sprintf("%s: %s,", fieldName, argName))
	}

	data := map[string]any{
		"Name":        ir.Name,
		"Args":        strings.Join(args, ", "),
		"Assignments": assignments,
		"TypeParams":  e.EmitTypeParams(ir.TypeParams, true),
		"Generics":    e.EmitTypeParams(ir.TypeParams, false),
	}

	tmpl, _ := template.New("test").Parse(CreateConstructorTemplate)

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit go code", ir.Span.ToLocation())
	}

	return tpl.String(), nil
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
		"Name":       ir.Name,
		"Fields":     fields,
		"TypeParams": e.EmitTypeParams(ir.TypeParams, true),
	}

	tmpl, _ := template.New("test").Parse(ModelTemplate)

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit go code", ir.Span.ToLocation())
	}

	sb.WriteString(tpl.String())

	if ir.IsCreateConstructor {
		constructorCode, err := e.EmitCreateConstructor(ir)
		if err != nil {
			return "", err
		}

		sb.WriteString("\n")
		sb.WriteString(constructorCode)
	}

	return sb.String(), nil
}

func (e *GoEmitter) Emit(ir *generator.ProgramIR) (string, exception.IException) {
	var models strings.Builder
	var sb strings.Builder

	for _, model := range ir.Models {
		code, err := e.EmitModel(model)
		if err != nil {
			return "", err
		}
		models.WriteString(code)
	}

	data := map[string]string{
		"Models": models.String(),
	}

	tmpl, _ := template.New("test").Parse(BaseTemplate)

	var tpl bytes.Buffer
	_ = tmpl.Execute(&tpl, data)

	sb.WriteString(tpl.String())
	return sb.String(), nil
}

func NewGoEmitter() *GoEmitter {
	return &GoEmitter{}
}
