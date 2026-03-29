package java

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

type JavaEmitter struct{}

func NewJavaEmitter() *JavaEmitter {
	return &JavaEmitter{}
}

func (e *JavaEmitter) emitTypeParams(params []string) string {
	if len(params) == 0 {
		return ""
	}

	return fmt.Sprintf("<%s>", strings.Join(params, ", "))
}

func (e *JavaEmitter) EmitBuiltinType(ir *generator.TypeIR) (string, exception.IException) {
	switch ir.Name {
	case "Int":
		return "Integer", nil
	case "Float":
		return "Double", nil
	case "Bool":
		return "Boolean", nil
	case "Any":
		return "Object", nil
	case "String":
		return "String", nil
	case "Null":
		return "Object", nil
	case "Array":
		if len(ir.Generics) == 0 {
			return "List<Object>", nil
		}

		genericType, err := e.EmitType(ir.Generics[0])
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("List<%s>", genericType), nil
	default:
		return "Object", nil
	}
}

func (e *JavaEmitter) EmitType(ir *generator.TypeIR) (string, exception.IException) {
	if ir == nil {
		return "Object", nil
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
		return "Object", nil
	}
}

func (e *JavaEmitter) EmitValue(ir *generator.ValueIR) (string, exception.IException) {
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
			return "java.util.Arrays.asList()", nil
		}

		items := make([]string, 0, len(arr))
		for _, item := range arr {
			value, err := e.EmitValue(item)
			if err != nil {
				return "", err
			}

			items = append(items, value)
		}

		return fmt.Sprintf("java.util.Arrays.asList(%s)", strings.Join(items, ", ")), nil
	default:
		return "null", nil
	}
}

func (e *JavaEmitter) EmitModelField(ir *generator.ModelField) (string, exception.IException) {
	typeStr, err := e.EmitType(ir.Type)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("public %s %s;", typeStr, ir.Name), nil
}

func (e *JavaEmitter) EmitConstructor(ir *generator.ModelIR) (string, exception.IException) {
	params := make([]string, 0, len(ir.Fields))
	assignments := make([]string, 0, len(ir.Fields))

	for _, field := range ir.Fields {
		typeStr, err := e.EmitType(field.Type)
		if err != nil {
			return "", err
		}

		params = append(params, fmt.Sprintf("%s %s", typeStr, field.Name))
		assignments = append(assignments, fmt.Sprintf("this.%s = %s;", field.Name, field.Name))
	}

	tmpl, _ := template.New("java-constructor").Parse(ConstructorTemplate)
	data := map[string]any{
		"Name":        ir.Name,
		"Params":      strings.Join(params, ", "),
		"Assignments": assignments,
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit java constructor", ir.Span.ToLocation())
	}

	return tpl.String(), nil
}

func (e *JavaEmitter) EmitFieldValidator(v *generator.FieldValidator, field *generator.ModelField) (string, exception.IException) {
	fieldRef := fmt.Sprintf("this.%s", field.Name)

	if v.Name == "IsModel" && field.Type != nil && field.Type.Kind == generator.TypeKindModel {
		validateCall := fmt.Sprintf("%s.validate();", fieldRef)
		if field.IsOptional {
			return fmt.Sprintf("if (%s != null) { %s }", fieldRef, validateCall), nil
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

	line := fmt.Sprintf("Validators.%s(%s);", v.Name, callArgs)
	if field.IsOptional && v.Name != "NotNull" {
		return fmt.Sprintf("if (%s != null) { %s }", fieldRef, line), nil
	}

	return line, nil
}

func (e *JavaEmitter) EmitValidatorMethod(ir *generator.ModelIR) (string, exception.IException) {
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

	tmpl, _ := template.New("java-validator-method").Parse(ValidatorMethodTemplate)
	data := map[string]any{
		"Lines": lines,
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit java validator", ir.Span.ToLocation())
	}

	return tpl.String(), nil
}

func (e *JavaEmitter) EmitModel(ir *generator.ModelIR) (string, exception.IException) {
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
		code, err := e.EmitConstructor(ir)
		if err != nil {
			return "", err
		}
		constructorCode = code
	}

	validatorCode, err := e.EmitValidatorMethod(ir)
	if err != nil {
		return "", err
	}

	tmpl, _ := template.New("java-model").Parse(ModelTemplate)
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
		return "", exception.NewEmitException("Error when emit java model", ir.Span.ToLocation())
	}

	return tpl.String(), nil
}

func (e *JavaEmitter) EmitRest(ir *generator.RestEndpointIR) (string, exception.IException) {
	restName := helpers.ToPascalCase(ir.Name)

	requestType := "\"any\""
	if ir.RequestBodyType != nil {
		t, err := e.EmitType(ir.RequestBodyType)
		if err != nil {
			return "", err
		}
		requestType = strconv.Quote(t)
	}

	responseType := "\"any\""
	if ir.ResponseBodyType != nil {
		t, err := e.EmitType(ir.ResponseBodyType)
		if err != nil {
			return "", err
		}
		responseType = strconv.Quote(t)
	}

	queries := make([]string, 0, len(ir.Queries))
	for _, q := range ir.Queries {
		queries = append(queries, strconv.Quote(q))
	}

	tmpl, _ := template.New("java-rest").Parse(RestTemplate)
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
		return "", exception.NewEmitException("Error when emit java rest", ir.Span.ToLocation())
	}

	return tpl.String(), nil
}

func (e *JavaEmitter) Emit(ir *generator.ProgramIR) (string, exception.IException) {
	var models strings.Builder
	var rests strings.Builder

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

	tmpl, _ := template.New("java-base").Parse(BaseTemplate)
	data := map[string]any{
		"Runtime":  strings.TrimSpace(RuntimeTemplate),
		"Models":   strings.TrimSpace(models.String()),
		"Rests":    strings.TrimSpace(rests.String()),
		"HasRests": rests.Len() > 0,
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", exception.NewEmitException("Error when emit java output", nil)
	}

	return strings.TrimSpace(tpl.String()) + "\n", nil
}
