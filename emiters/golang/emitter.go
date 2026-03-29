package golang

import (
	"bytes"
	"fmt"
	"strconv"
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

func (e *GoEmitter) EmitModelType(ir *generator.TypeIR, isOptional bool) (string, exception.IException) {
	genericTypes := []string{}

	for _, generic := range ir.Generics {
		genericType, err := e.EmitType(generic, false)
		if err != nil {
			return "", err
		}

		genericTypes = append(genericTypes, genericType)
	}

	if len(genericTypes) > 0 {
		base := fmt.Sprintf(`%s[%s]`, ir.Name, strings.Join(genericTypes, ","))
		return fmt.Sprintf(`*%s`, base), nil
	}

	_ = isOptional
	return fmt.Sprintf(`*%s`, ir.Name), nil
}

func (e *GoEmitter) EmitType(ir *generator.TypeIR, isOptional bool) (string, exception.IException) {
	if ir == nil {
		return "any", nil
	}

	if ir.Kind == generator.TypeKindBuiltin {
		t, err := e.EmitBuildInType(ir)
		if isOptional {
			t = fmt.Sprintf("*%s", t)
		}

		return t, err
	}

	if ir.Kind == generator.TypeKindModel {
		return e.EmitModelType(ir, isOptional)
	}

	if ir.Kind == generator.TypeKindGeneric {
		t, err := e.EmitGenericType(ir)
		if isOptional {
			t = fmt.Sprintf("*%s", t)
		}

		return t, err
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

	if ir.Kind == "Boolean" {
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

		return fmt.Sprintf(`[]any{%s}`, strings.Join(elements, ",")), nil
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

	validatorCode, err := e.EmitValidator(ir)
	if err != nil {
		return "", exception.NewEmitException("Error when emit go code", ir.Span.ToLocation())
	}

	sb.WriteString(validatorCode)

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

func (e *GoEmitter) EmitValidator(ir *generator.ModelIR) (string, exception.IException) {
	var sb strings.Builder

	for _, field := range ir.Fields {
		for _, validator := range field.Validators {
			validatorCode, err := e.EmitFieldValidator(validator, field.Name, field.Type.Kind == generator.TypeKindModel, field.IsOptional, field.Type)
			if err != nil {
				return "", err
			}
			sb.WriteString(validatorCode + "\n")
		}
	}

	data := map[string]string{
		"Name":       ir.Name,
		"Validators": sb.String(),
		"TypeParams": e.EmitTypeParams(ir.TypeParams, true),
		"Generics":   e.EmitTypeParams(ir.TypeParams, false),
	}

	tmpl, _ := template.New("test").Parse(ValidatorTemplate)

	var tpl bytes.Buffer
	_ = tmpl.Execute(&tpl, data)

	return tpl.String(), nil
}

func (e *GoEmitter) EmitFieldValidator(ir *generator.FieldValidator, name string, isModelType bool, isOptional bool, typeDef *generator.TypeIR) (string, exception.IException) {

	args := []string{}
	for _, arg := range ir.Args {
		argValue, err := e.EmitValue(arg)
		if err != nil {
			return "", err
		}
		args = append(args, argValue)
	}

	if ir.Name == "IsModel" && isModelType {
		_ = isOptional
		modelValue := fmt.Sprintf("instance.%s", helpers.ToPascalCase(name))
		return fmt.Sprintf(`if err := %sValidate(%s); err != nil { return NewValidatorError(%s) }`, typeDef.Name, modelValue, args[0]), nil
	}

	return fmt.Sprintf(`if err := %s(instance.%s,%s); err != nil { return err }`, ir.Name, helpers.ToPascalCase(name), strings.Join(args, ",")), nil
}

func (e *GoEmitter) EmitRest(ir *generator.RestEndpointIR) (string, exception.IException) {
	restName := helpers.ToPascalCase(ir.Name)

	requestType := "any"
	if ir.RequestBodyType != nil {
		t, err := e.EmitType(ir.RequestBodyType, false)
		if err != nil {
			return "", err
		}
		requestType = t
	}

	responseType := "any"
	if ir.ResponseBodyType != nil {
		t, err := e.EmitType(ir.ResponseBodyType, false)
		if err != nil {
			return "", err
		}
		responseType = t
	}

	queryValues := make([]string, 0, len(ir.Queries))
	for _, q := range ir.Queries {
		queryValues = append(queryValues, strconv.Quote(q))
	}

	tmpl, _ := template.New("go-rest").Parse(RestTemplate)

	data := map[string]any{
		"Name":         restName,
		"Path":         strconv.Quote(ir.Path),
		"Method":       strconv.Quote(ir.Method),
		"Queries":      strings.Join(queryValues, ", "),
		"RequestType":  requestType,
		"ResponseType": responseType,
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit go code", ir.Span.ToLocation())
	}

	return tpl.String(), nil
}

func (e *GoEmitter) Emit(ir *generator.ProgramIR) (string, exception.IException) {
	var models strings.Builder
	var rests strings.Builder
	var sb strings.Builder

	for _, model := range ir.Models {
		code, err := e.EmitModel(model)
		if err != nil {
			return "", err
		}
		models.WriteString(code)
	}

	for _, rest := range ir.Rests {
		code, err := e.EmitRest(rest)
		if err != nil {
			return "", err
		}

		rests.WriteString(code)
		rests.WriteString("\n")
	}

	data := map[string]any{
		"Models":   models.String(),
		"Rests":    rests.String(),
		"HasRests": len(ir.Rests) > 0,
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
