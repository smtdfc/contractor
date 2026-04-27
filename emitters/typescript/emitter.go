package typescript

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/smtdfc/contractor/exception"
	"github.com/smtdfc/contractor/generator"
)

type TypescriptEmitter struct{}

var typeMap = map[string]string{
	"Int":    "number",
	"Float":  "number",
	"String": "string",
	"Bool":   "boolean",
	"Null":   "null",
	"Any":    "any",
}

func (t *TypescriptEmitter) EmitTypeName(ir *generator.TypeIR) (string, exception.IException) {
	var typeName strings.Builder

	if ir.Kind == generator.TypeKindBuiltin {
		tsType, ok := typeMap[ir.Name]
		if !ok {
			typeName.WriteString("unknown")
		}

		typeName.WriteString(tsType)
	}

	if ir.Kind == generator.TypeKindModel {
		typeName.WriteString(ir.Name)
	}

	if ir.Kind == generator.TypeKindGeneric {
		typeName.WriteString(ir.Name)
	}

	if len(ir.Generics) > 0 {
		typeName.WriteString("<")
		genericTypes := []string{}

		for _, generic := range ir.Generics {
			tsGenericType, err := t.EmitTypeName(generic)
			if err != nil {
				return "", err
			}

			genericTypes = append(genericTypes, tsGenericType)
		}

		typeName.WriteString(strings.Join(genericTypes, ","))
		typeName.WriteString(">")
	}

	return typeName.String(), nil
}

func (t *TypescriptEmitter) EmitModel(tmpl *template.Template, ir *generator.ModelIR) (string, exception.IException) {
	var sb strings.Builder
	var data = map[string]any{
		"ModelName":         ir.Name,
		"CreateConstructor": ir.IsCreateConstructor,
		"TypeParams":        ir.TypeParams,
		"IsGeneric":         len(ir.TypeParams) > 0,
		"IsMapper":          true,
	}

	var fields = []any{}
	var fieldValidators = []any{}
	for _, field := range ir.Fields {
		fieldTypeName, err := t.EmitTypeName(field.Type)
		if err != nil {
			return "", err
		}

		fields = append(fields, map[string]any{
			"Name":          field.Name,
			"IsOptional":    field.IsOptional,
			"Type":          fieldTypeName,
			"IsModelType":   field.Type.Kind == generator.TypeKindModel,
			"ModelTypeName": field.Type.Name,
		})

		isModelType := field.Type.Kind == generator.TypeKindModel
		isArrayOfModelType := false
		modelTypeName := field.Type.Name
		if field.Type.Kind == generator.TypeKindBuiltin && field.Type.Name == "Array" && len(field.Type.Generics) == 1 {
			genericItem := field.Type.Generics[0]
			if genericItem != nil && genericItem.Kind == generator.TypeKindModel {
				isArrayOfModelType = true
				modelTypeName = genericItem.Name
			}
		}

		validators := []any{}
		for _, validator := range field.Validators {
			args := make([]string, 0, len(validator.Args))
			for _, arg := range validator.Args {
				args = append(args, emitValueLiteral(arg))
			}

			validators = append(validators, map[string]any{
				"Name":               validator.Name,
				"Args":               args,
				"IsNestedValidate":   validator.Name == "NestedValidate",
				"Field":              field.Name,
				"IsModelType":        isModelType,
				"IsArrayOfModelType": isArrayOfModelType,
				"ModelTypeName":      modelTypeName,
			})
		}

		if len(validators) > 0 {
			fieldValidators = append(fieldValidators, map[string]any{
				"Field":              field.Name,
				"IsOptional":         field.IsOptional,
				"IsModelType":        isModelType,
				"IsArrayOfModelType": isArrayOfModelType,
				"ModelTypeName":      modelTypeName,
				"Validators":         validators,
			})
		}

	}
	data["Fields"] = fields
	data["FieldValidators"] = fieldValidators

	if err := tmpl.ExecuteTemplate(&sb, "model.tmpl", data); err != nil {
		return "", exception.NewEmitException(err.Error(), nil)
	}
	return sb.String(), nil
}

func (t *TypescriptEmitter) EmitRest(tmpl *template.Template, ir *generator.RestEndpointIR) (string, exception.IException) {
	var sb strings.Builder

	requestTypeName := "unknown"
	if ir.RequestBodyType != nil {
		emitted, err := t.EmitTypeName(ir.RequestBodyType)
		if err != nil {
			return "", err
		}
		requestTypeName = emitted
	}

	responseTypeName := "unknown"
	if ir.ResponseBodyType != nil {
		emitted, err := t.EmitTypeName(ir.ResponseBodyType)
		if err != nil {
			return "", err
		}
		responseTypeName = emitted
	}

	queryLiterals := make([]string, 0, len(ir.Queries))
	for _, query := range ir.Queries {
		queryLiterals = append(queryLiterals, strconv.Quote(query))
	}

	data := map[string]any{
		"Name":         ir.Name,
		"Path":         strconv.Quote(ir.Path),
		"Method":       strconv.Quote(ir.Method),
		"Queries":      queryLiterals,
		"RequestType":  requestTypeName,
		"ResponseType": responseTypeName,
	}

	if err := tmpl.ExecuteTemplate(&sb, "rest.tmpl", data); err != nil {
		return "", exception.NewEmitException(err.Error(), nil)
	}

	return sb.String(), nil
}

func (t *TypescriptEmitter) Emit(ir *generator.ProgramIR) (string, exception.IException) {
	var sb strings.Builder
	tmpl, err := template.ParseFS(templateFiles, "templates/*.tmpl")
	if err != nil {
		return "", exception.NewEmitException(err.Error(), nil)
	}

	sb.WriteString("// @ts-nocheck\n")
	sb.WriteString("import { Validator } from \"contractor-ts\";\n\n")
	sb.WriteString("import type { RestMetadata, RestRequestBody, RestResponseBody } from \"contractor-ts\";\n\n")

	for _, model := range ir.Models {
		code, err := t.EmitModel(tmpl, model)
		if err != nil {
			return "", err
		}

		sb.WriteString(code)
	}

	for _, rest := range ir.Rests {
		code, err := t.EmitRest(tmpl, rest)
		if err != nil {
			return "", err
		}

		sb.WriteString(code)
	}

	return sb.String(), nil
}

func NewTypescriptEmitter() *TypescriptEmitter {
	return &TypescriptEmitter{}
}

func emitValueLiteral(value *generator.ValueIR) string {
	if value == nil {
		return "null"
	}

	switch value.Kind {
	case "String":
		if raw, ok := value.Value.(string); ok {
			return strconv.Quote(raw)
		}
		return "\"\""
	case "Number", "Boolean":
		if raw, ok := value.Value.(string); ok {
			return raw
		}
		return fmt.Sprint(value.Value)
	case "Null":
		return "null"
	case "Array":
		rawValues, ok := value.Value.([]*generator.ValueIR)
		if !ok {
			return "[]"
		}

		items := make([]string, 0, len(rawValues))
		for _, item := range rawValues {
			items = append(items, emitValueLiteral(item))
		}

		return "[" + strings.Join(items, ", ") + "]"
	default:
		if value.Value == nil {
			return "null"
		}
		return fmt.Sprint(value.Value)
	}
}
