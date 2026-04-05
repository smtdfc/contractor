package typescript

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/smtdfc/contractor/emitters"
	"github.com/smtdfc/contractor/exception"
	"github.com/smtdfc/contractor/generator"
	"github.com/smtdfc/contractor/internal/helpers"
)

type TypescriptEmitter struct{}

var _ emitters.ProgramEmitter = (*TypescriptEmitter)(nil)

func NewTypescriptEmitter() *TypescriptEmitter {
	return &TypescriptEmitter{}
}

func (e *TypescriptEmitter) EmitTypeParams(params []string) string {
	if len(params) == 0 {
		return ""
	}

	return fmt.Sprintf("<%s>", strings.Join(params, ", "))
}

func (e *TypescriptEmitter) EmitBuiltinType(ir *generator.TypeIR) (string, exception.IException) {
	switch ir.Name {
	case "Int", "Float":
		return "number", nil
	case "Bool":
		return "boolean", nil
	case "Any":
		return "any", nil
	case "String":
		return "string", nil
	case "Null":
		return "null", nil
	case "Array":
		if len(ir.Generics) == 0 {
			return "any[]", nil
		}

		genericType, err := e.EmitType(ir.Generics[0])
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s[]", genericType), nil
	default:
		return "any", nil
	}
}

func (e *TypescriptEmitter) EmitType(ir *generator.TypeIR) (string, exception.IException) {
	if ir == nil {
		return "any", nil
	}

	switch ir.Kind {
	case generator.TypeKindBuiltin:
		return e.EmitBuiltinType(ir)
	case generator.TypeKindModel:
		if len(ir.Generics) == 0 {
			return ir.Name, nil
		}

		genericTypes := make([]string, 0, len(ir.Generics))
		for _, item := range ir.Generics {
			t, err := e.EmitType(item)
			if err != nil {
				return "", err
			}

			genericTypes = append(genericTypes, t)
		}

		return fmt.Sprintf("%s<%s>", ir.Name, strings.Join(genericTypes, ", ")), nil
	case generator.TypeKindGeneric:
		return ir.Name, nil
	default:
		return "any", nil
	}
}

func (e *TypescriptEmitter) EmitValue(ir *generator.ValueIR) (string, exception.IException) {
	if ir == nil {
		return "null", nil
	}

	switch ir.Kind {
	case "Number", "Boolean":
		return fmt.Sprintf("%v", ir.Value), nil
	case "String":
		value, _ := ir.Value.(string)
		return strconv.Quote(value), nil
	case "Null":
		return "null", nil
	case "Array":
		arr, ok := ir.Value.([]*generator.ValueIR)
		if !ok {
			return "[]", nil
		}

		items := make([]string, 0, len(arr))
		for _, item := range arr {
			value, err := e.EmitValue(item)
			if err != nil {
				return "", err
			}

			items = append(items, value)
		}

		return fmt.Sprintf("[%s]", strings.Join(items, ", ")), nil
	default:
		return "null", nil
	}
}

func (e *TypescriptEmitter) EmitModelField(ir *generator.ModelField) (string, exception.IException) {
	typeStr, err := e.EmitType(ir.Type)
	if err != nil {
		return "", err
	}

	if ir.IsOptional {
		return fmt.Sprintf("%s?: %s;", ir.Name, typeStr), nil
	}

	return fmt.Sprintf("// @ts-ignore: generated field may be initialized outside constructor\n%s: %s;", ir.Name, typeStr), nil
}

func (e *TypescriptEmitter) EmitConstructor(ir *generator.ModelIR) (string, exception.IException) {
	params := make([]string, 0, len(ir.Fields))
	assignments := make([]string, 0, len(ir.Fields))

	for _, field := range ir.Fields {
		typeStr, err := e.EmitType(field.Type)
		if err != nil {
			return "", err
		}

		if field.IsOptional {
			params = append(params, fmt.Sprintf("%s?: %s", field.Name, typeStr))
		} else {
			params = append(params, fmt.Sprintf("%s: %s", field.Name, typeStr))
		}

		assignments = append(assignments, fmt.Sprintf("this.%s = %s;", field.Name, field.Name))
	}

	tmpl, _ := template.New("ts-constructor").Parse(ConstructorTemplate)

	data := map[string]any{
		"Params":      strings.Join(params, ", "),
		"Assignments": assignments,
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit typescript constructor", ir.Span.ToLocation())
	}

	return tpl.String(), nil
}

func (e *TypescriptEmitter) EmitFieldValidator(v *generator.FieldValidator, field *generator.ModelField) (string, exception.IException) {
	fieldRef := fmt.Sprintf("this.%s", field.Name)

	if v.Name == "IsModel" && field.Type != nil && field.Type.Kind == generator.TypeKindModel {
		validateCall := fmt.Sprintf("%s.validate(%s);", field.Type.Name, fieldRef)
		if field.IsOptional {
			return fmt.Sprintf("if (%s !== undefined) { %s }", fieldRef, validateCall), nil
		}

		return validateCall, nil
	}

	args := make([]string, 0, len(v.Args))
	for _, arg := range v.Args {
		value, err := e.EmitValue(arg)
		if err != nil {
			return "", err
		}

		args = append(args, value)
	}

	callArgs := fieldRef
	if len(args) > 0 {
		callArgs = fmt.Sprintf("%s, %s", callArgs, strings.Join(args, ", "))
	}

	line := fmt.Sprintf("%s(%s);", v.Name, callArgs)
	if field.IsOptional && v.Name != "NotNull" {
		return fmt.Sprintf("if (%s !== undefined) { %s }", fieldRef, line), nil
	}

	return line, nil
}

func (e *TypescriptEmitter) EmitValidatorMethod(ir *generator.ModelIR) (string, exception.IException) {
	lines := make([]string, 0)

	for _, field := range ir.Fields {
		for _, validator := range field.Validators {
			line, err := e.EmitFieldValidator(validator, field)
			if err != nil {
				return "", err
			}

			lines = append(lines, line)
		}
	}

	tmpl, _ := template.New("ts-validator-method").Parse(ValidatorMethodTemplate)

	data := map[string]any{
		"Lines": lines,
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit typescript validator", ir.Span.ToLocation())
	}

	return tpl.String(), nil
}

func (e *TypescriptEmitter) EmitModel(ir *generator.ModelIR) (string, exception.IException) {
	fields := make([]string, 0, len(ir.Fields))
	for _, field := range ir.Fields {
		line, err := e.EmitModelField(field)
		if err != nil {
			return "", err
		}

		fields = append(fields, line)
	}

	constructorCode := ""
	if ir.IsCreateConstructor {
		c, err := e.EmitConstructor(ir)
		if err != nil {
			return "", exception.NewEmitException("Error when emit typescript constructor", ir.Span.ToLocation())
		}
		constructorCode = c
	}

	validatorCode, err := e.EmitValidatorMethod(ir)
	if err != nil {
		return "", exception.NewEmitException("Error when emit typescript validator", ir.Span.ToLocation())
	}

	tmpl, _ := template.New("ts-model").Parse(ModelTemplate)

	data := map[string]any{
		"Name":            ir.Name,
		"TypeParams":      e.EmitTypeParams(ir.TypeParams),
		"Fields":          fields,
		"ConstructorCode": constructorCode,
		"ValidatorCode":   validatorCode,
	}

	var tpl bytes.Buffer
	err_ := tmpl.Execute(&tpl, data)
	if err_ != nil {
		return "", exception.NewEmitException("Error when emit typescript model", ir.Span.ToLocation())
	}

	return tpl.String(), nil
}

func (e *TypescriptEmitter) EmitRest(ir *generator.RestEndpointIR) (string, exception.IException) {
	restName := helpers.ToPascalCase(ir.Name)

	requestType := "any"
	if ir.RequestBodyType != nil {
		t, err := e.EmitType(ir.RequestBodyType)
		if err != nil {
			return "", err
		}
		requestType = t
	}

	responseType := "any"
	if ir.ResponseBodyType != nil {
		t, err := e.EmitType(ir.ResponseBodyType)
		if err != nil {
			return "", err
		}
		responseType = t
	}

	queries := make([]string, 0, len(ir.Queries))
	for _, query := range ir.Queries {
		queries = append(queries, strconv.Quote(query))
	}

	tmpl, _ := template.New("ts-rest").Parse(RestTemplate)

	data := map[string]any{
		"RestName":     restName,
		"Path":         strconv.Quote(ir.Path),
		"Method":       strconv.Quote(ir.Method),
		"Queries":      strings.Join(queries, ", "),
		"RequestType":  requestType,
		"ResponseType": responseType,
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit typescript rest", ir.Span.ToLocation())
	}

	return tpl.String(), nil
}

func (e *TypescriptEmitter) EmitError(ir *generator.ErrorIR) (string, exception.IException) {
	data := map[string]any{
		"Name":     ir.Name,
		"Message":  strconv.Quote(ir.Message),
		"HasCode":  ir.Code != nil,
		"HasScope": ir.Scope != nil,
	}

	if ir.Code != nil {
		data["Code"] = strconv.Quote(*ir.Code)
	}
	if ir.Scope != nil {
		data["Scope"] = strconv.Quote(*ir.Scope)
	}

	tmpl, _ := template.New("ts-error").Parse(ErrorTemplate)

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit typescript error", ir.Span.ToLocation())
	}

	return tpl.String(), nil
}

func (e *TypescriptEmitter) EmitErrorMap(items []*generator.ErrorIR) string {
	if len(items) == 0 {
		return ""
	}

	seen := make(map[string]struct{})
	entries := make([]string, 0)

	for _, item := range items {
		if item.Code == nil {
			continue
		}

		key := *item.Code
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		entries = append(entries, fmt.Sprintf("  %s: () => new %s(),", strconv.Quote(key), item.Name))
	}

	if len(entries) == 0 {
		return ""
	}

	return strings.Join([]string{
		"export const ErrorMap: Record<string, () => Error> = {",
		strings.Join(entries, "\n"),
		"};",
		"",
		"export function createErrorByKey(key: string): Error | null {",
		"  const factory = ErrorMap[key];",
		"  return factory ? factory() : null;",
		"}",
	}, "\n")
}

func (e *TypescriptEmitter) Emit(ir *generator.ProgramIR) (string, exception.IException) {
	var errors strings.Builder
	var models strings.Builder
	var rests strings.Builder

	for _, item := range ir.Errors {
		code, err := e.EmitError(item)
		if err != nil {
			return "", err
		}

		errors.WriteString(code)
		errors.WriteString("\n")
	}

	for _, model := range ir.Models {
		code, err := e.EmitModel(model)
		if err != nil {
			return "", err
		}

		models.WriteString(code)
		models.WriteString("\n")
	}

	for _, rest := range ir.Rests {
		code, err := e.EmitRest(rest)
		if err != nil {
			return "", err
		}

		rests.WriteString(code)
		rests.WriteString("\n")
	}

	tmpl, _ := template.New("ts-base").Parse(BaseTemplate)

	data := map[string]any{
		"Runtime":  strings.TrimSpace(RuntimeTemplate),
		"Errors":   strings.TrimSpace(errors.String()),
		"ErrorMap": strings.TrimSpace(e.EmitErrorMap(ir.Errors)),
		"Models":   strings.TrimSpace(models.String()),
		"Rests":    strings.TrimSpace(rests.String()),
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit typescript base", nil)
	}

	return strings.TrimSpace(tpl.String()) + "\n", nil
}
