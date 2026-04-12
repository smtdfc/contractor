package csharp

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

type CSharpEmitter struct{}

var _ emitters.ProgramEmitter = (*CSharpEmitter)(nil)

func NewCSharpEmitter() *CSharpEmitter {
	return &CSharpEmitter{}
}

func (e *CSharpEmitter) isNestedModelType(typeIR *generator.TypeIR) bool {
	if typeIR == nil {
		return false
	}

	if typeIR.Kind == generator.TypeKindModel {
		return true
	}

	if typeIR.Kind == generator.TypeKindBuiltin && typeIR.Name == "Array" && len(typeIR.Generics) == 1 {
		return typeIR.Generics[0] != nil && typeIR.Generics[0].Kind == generator.TypeKindModel
	}

	return false
}

func (e *CSharpEmitter) emitTypeParams(params []string) string {
	if len(params) == 0 {
		return ""
	}

	return fmt.Sprintf("<%s>", strings.Join(params, ", "))
}

func (e *CSharpEmitter) EmitBuiltinType(ir *generator.TypeIR) (string, exception.IException) {
	switch ir.Name {
	case "Int":
		return "int", nil
	case "Float":
		return "double", nil
	case "Bool":
		return "bool", nil
	case "Any":
		return "object", nil
	case "String":
		return "string", nil
	case "Null":
		return "object?", nil
	case "Array":
		if len(ir.Generics) == 0 {
			return "List<object?>", nil
		}

		genericType, err := e.EmitType(ir.Generics[0], false)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("List<%s>", genericType), nil
	default:
		return "object", nil
	}
}

func (e *CSharpEmitter) EmitType(ir *generator.TypeIR, isOptional bool) (string, exception.IException) {
	if ir == nil {
		if isOptional {
			return "object?", nil
		}
		return "object", nil
	}

	var typeName string
	var err exception.IException

	switch ir.Kind {
	case generator.TypeKindBuiltin:
		typeName, err = e.EmitBuiltinType(ir)
		if err != nil {
			return "", err
		}
	case generator.TypeKindModel:
		if len(ir.Generics) == 0 {
			typeName = ir.Name
		} else {
			genericTypes := make([]string, 0, len(ir.Generics))
			for _, item := range ir.Generics {
				t, e2 := e.EmitType(item, false)
				if e2 != nil {
					return "", e2
				}
				genericTypes = append(genericTypes, t)
			}
			typeName = fmt.Sprintf("%s<%s>", ir.Name, strings.Join(genericTypes, ", "))
		}
	case generator.TypeKindGeneric:
		typeName = ir.Name
	default:
		typeName = "object"
	}

	if isOptional && !strings.HasSuffix(typeName, "?") {
		typeName += "?"
	}

	return typeName, nil
}

func (e *CSharpEmitter) EmitValue(ir *generator.ValueIR) (string, exception.IException) {
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
			return "Array.Empty<object?>()", nil
		}

		items := make([]string, 0, len(arr))
		for _, item := range arr {
			value, err := e.EmitValue(item)
			if err != nil {
				return "", err
			}
			items = append(items, value)
		}
		return fmt.Sprintf("new object?[] { %s }", strings.Join(items, ", ")), nil
	default:
		return "null", nil
	}
}

func (e *CSharpEmitter) EmitModelField(ir *generator.ModelField) (string, exception.IException) {
	typeStr, err := e.EmitType(ir.Type, ir.IsOptional)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("public %s %s { get; set; }", typeStr, helpers.ToPascalCase(ir.Name)), nil
}

func (e *CSharpEmitter) EmitConstructor(ir *generator.ModelIR) (string, exception.IException) {
	params := make([]string, 0, len(ir.Fields))
	assignments := make([]string, 0, len(ir.Fields))

	for _, field := range ir.Fields {
		typeStr, err := e.EmitType(field.Type, field.IsOptional)
		if err != nil {
			return "", err
		}

		paramName := helpers.ToCamelCase(field.Name)
		fieldName := helpers.ToPascalCase(field.Name)
		params = append(params, fmt.Sprintf("%s %s", typeStr, paramName))
		assignments = append(assignments, fmt.Sprintf("%s = %s;", fieldName, paramName))
	}

	tmpl, _ := template.New("cs-constructor").Parse(ConstructorTemplate)
	data := map[string]any{
		"Name":        ir.Name,
		"Params":      strings.Join(params, ", "),
		"Assignments": assignments,
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit csharp constructor", ir.Span.ToLocation())
	}

	return tpl.String(), nil
}

func (e *CSharpEmitter) EmitFieldValidator(v *generator.FieldValidator, field *generator.ModelField) (string, exception.IException) {
	fieldName := helpers.ToPascalCase(field.Name)
	fieldRef := fieldName
	defaultMsg := "invalid model"

	if v.Name == "NestedValidate" && e.isNestedModelType(field.Type) {
		msg := strconv.Quote(defaultMsg)
		if len(v.Args) > 0 {
			value, err := e.EmitValue(v.Args[0])
			if err != nil {
				return "", err
			}

			msg = value
		}

		return fmt.Sprintf("Validators.NestedValidate(%s, %s, %s);", fieldRef, strconv.Quote(field.Name), msg), nil
	}

	if v.Name == "IsModel" && field.Type != nil && field.Type.Kind == generator.TypeKindModel {
		msg := strconv.Quote(defaultMsg)
		if len(v.Args) > 0 {
			value, err := e.EmitValue(v.Args[0])
			if err != nil {
				return "", err
			}

			msg = value
		}

		return fmt.Sprintf("Validators.ValidateModel(%s, %s);", fieldRef, msg), nil
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

	line := fmt.Sprintf("Validators.%s(%s);", v.Name, callArgs)
	if field.IsOptional && v.Name != "NotNull" {
		return fmt.Sprintf("if (%s is not null) { %s }", fieldRef, line), nil
	}

	return line, nil
}

func (e *CSharpEmitter) EmitValidatorMethod(ir *generator.ModelIR) (string, exception.IException) {
	items := make([]struct {
		Field string
		Line  string
	}, 0)
	for _, field := range ir.Fields {
		for _, validator := range field.Validators {
			line, err := e.EmitFieldValidator(validator, field)
			if err != nil {
				return "", err
			}
			items = append(items, struct {
				Field string
				Line  string
			}{Field: field.Name, Line: line})
		}
	}

	tmpl, _ := template.New("cs-validator").Parse(ValidatorMethodTemplate)
	data := map[string]any{"Items": items}
	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit csharp validator", ir.Span.ToLocation())
	}
	return tpl.String(), nil
}

func (e *CSharpEmitter) EmitModel(ir *generator.ModelIR) (string, exception.IException) {
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
			return "", err
		}
		constructorCode = c
	}

	validatorCode, err := e.EmitValidatorMethod(ir)
	if err != nil {
		return "", err
	}

	tmpl, _ := template.New("cs-model").Parse(ModelTemplate)
	data := map[string]any{
		"Name":            ir.Name,
		"TypeParams":      e.emitTypeParams(ir.TypeParams),
		"Fields":          fields,
		"ConstructorCode": constructorCode,
		"ValidatorCode":   validatorCode,
	}
	var tpl bytes.Buffer
	err_ := tmpl.Execute(&tpl, data)
	if err_ != nil {
		return "", exception.NewEmitException("Error when emit csharp model", ir.Span.ToLocation())
	}
	return tpl.String(), nil
}

func (e *CSharpEmitter) EmitRest(ir *generator.RestEndpointIR) (string, exception.IException) {
	restName := helpers.ToPascalCase(ir.Name)

	requestType := "object"
	if ir.RequestBodyType != nil {
		t, err := e.EmitType(ir.RequestBodyType, false)
		if err != nil {
			return "", err
		}
		requestType = t
	}

	responseType := "object"
	if ir.ResponseBodyType != nil {
		t, err := e.EmitType(ir.ResponseBodyType, false)
		if err != nil {
			return "", err
		}
		responseType = t
	}

	queries := make([]string, 0, len(ir.Queries))
	for _, q := range ir.Queries {
		queries = append(queries, strconv.Quote(q))
	}

	tmpl, _ := template.New("cs-rest").Parse(RestTemplate)
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
		return "", exception.NewEmitException("Error when emit csharp rest", ir.Span.ToLocation())
	}
	return tpl.String(), nil
}

func (e *CSharpEmitter) EmitError(ir *generator.ErrorIR) (string, exception.IException) {
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

	tmpl, _ := template.New("cs-error").Parse(ErrorTemplate)

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit csharp error", ir.Span.ToLocation())
	}

	return tpl.String(), nil
}

func (e *CSharpEmitter) EmitErrorMap(items []*generator.ErrorIR) string {
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
		entries = append(entries, fmt.Sprintf("        [%s] = () => new %s(),", strconv.Quote(key), item.Name))
	}

	if len(entries) == 0 {
		return ""
	}

	return strings.Join([]string{
		"    public static readonly IReadOnlyDictionary<string, Func<Exception>> ErrorMap =",
		"        new Dictionary<string, Func<Exception>>",
		"        {",
		strings.Join(entries, "\n"),
		"        };",
		"",
		"    public static Exception? CreateErrorByKey(string key)",
		"        => ErrorMap.TryGetValue(key, out var factory) ? factory() : null;",
	}, "\n")
}

func (e *CSharpEmitter) Emit(ir *generator.ProgramIR) (string, exception.IException) {
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

	tmpl, _ := template.New("cs-base").Parse(BaseTemplate)
	data := map[string]any{
		"Runtime":  strings.TrimSpace(RuntimeTemplate),
		"Errors":   strings.TrimSpace(errors.String()),
		"ErrorMap": strings.TrimSpace(e.EmitErrorMap(ir.Errors)),
		"Models":   strings.TrimSpace(models.String()),
		"Rests":    strings.TrimSpace(rests.String()),
		"HasRests": rests.Len() > 0,
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit csharp output", nil)
	}

	return strings.TrimSpace(tpl.String()) + "\n", nil
}
