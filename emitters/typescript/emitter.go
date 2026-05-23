package typescript

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"crypto/md5"

	"github.com/smtdfc/contractor/exception"
	"github.com/smtdfc/contractor/generator"
)

func hashModelName(name string) string {
	hasher := md5.New()
	hasher.Write([]byte(name))
	return fmt.Sprintf("%x", hasher.Sum(nil))[:8]
}

type TypescriptEmitter struct{}

var typeMap = map[string]string{
	"Int":    "number",
	"Float":  "number",
	"String": "string",
	"Bool":   "boolean",
	"Null":   "null",
	"Any":    "any",
	"Array":  "Array",
}

func (t *TypescriptEmitter) EmitTypeName(ir *generator.TypeIR) (string, *exception.EmitException) {
	var typeName strings.Builder

	if ir.Kind == generator.TypeKindBuiltin {
		tsType, ok := typeMap[ir.Name]
		if !ok {
			tsType = "unknown"
		}

		typeName.WriteString(tsType)
	}

	if ir.Kind == generator.TypeKindModel {
		typeName.WriteString(ir.Name)
	}

	if ir.Kind == generator.TypeKindEnum {
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

type ModelGenerateData struct {
	Name        string
	ModelHash   string
	Description string
	IsGeneric   bool
	TypeParams  []string
	Fields      []ModelFieldGenerateData
	Validators  []ModelFieldValidatorGenerateData
}

type ModelFieldGenerateData struct {
	Name                 string
	IsOptional           bool
	Type                 string
	IsArray              bool
	IsModel              bool
	IsGeneric            bool
	ModelTypeName        string
	GenericTypeName      string
	ArrayModelTypeName   string
	ArrayGenericTypeName string
}

type ModelFieldValidateRuleGenerateData struct {
	Name string
	Args []string
}

type ModelFieldValidatorGenerateData struct {
	IsArray       bool
	IsModel       bool
	IsGenericType bool
	TypeName      string
	Field         string
	Rules         []ModelFieldValidateRuleGenerateData
}

type RestGenerateData struct {
	Name     string
	Path     string
	Method   string
	ResType  string
	ReqType  string
	ResClass string
	ReqClass string
}

func (t *TypescriptEmitter) EmitModel(tmpl *template.Template, ir *generator.ModelIR) (string, *exception.EmitException) {
	var wr bytes.Buffer

	var fields []ModelFieldGenerateData
	var validators []ModelFieldValidatorGenerateData

	for _, field := range ir.Fields {
		fieldType, err := t.EmitTypeName(field.Type)
		if err != nil {
			return "", err
		}

		isFieldArrayType := false
		isFieldModelType := false
		isFieldGenericType := false
		modelTypeName := ""
		genericTypeName := ""
		arrayModelTypeName := ""
		arrayGenericTypeName := ""

		if field.Type.Kind == generator.TypeKindBuiltin && field.Type.Name == "Array" {
			isFieldArrayType = true
			if len(field.Type.Generics) > 0 {
				innerType := field.Type.Generics[0]
				if innerType.Kind == generator.TypeKindModel {
					arrayModelTypeName = innerType.Name
				}

				if innerType.Kind == generator.TypeKindGeneric {
					arrayGenericTypeName = innerType.Name
				}
			}
		}

		if field.Type.Kind == generator.TypeKindModel {
			isFieldModelType = true
			modelTypeName = field.Type.Name
		}

		if field.Type.Kind == generator.TypeKindGeneric {
			isFieldGenericType = true
			genericTypeName = field.Type.Name
		}

		isArrayType := false
		isModelType := false
		isGenericType := false
		typeName := ""
		if field.Type.Kind == generator.TypeKindBuiltin && field.Type.Name == "Array" {
			isArrayType = true
			if len(field.Type.Generics) > 0 {
				innerType := field.Type.Generics[0]
				typeName = innerType.Name

				if innerType.Kind == generator.TypeKindModel {
					isModelType = true
				}
			}
		}

		if field.Type.Kind == generator.TypeKindModel {
			isModelType = true
			typeName = field.Type.Name
		}

		if field.Type.Kind == generator.TypeKindGeneric {
			isGenericType = true
		}

		fields = append(fields, ModelFieldGenerateData{
			Name:                 field.Name,
			Type:                 fieldType,
			IsOptional:           field.IsOptional,
			IsArray:              isFieldArrayType,
			IsModel:              isFieldModelType,
			IsGeneric:            isFieldGenericType,
			ModelTypeName:        modelTypeName,
			GenericTypeName:      genericTypeName,
			ArrayModelTypeName:   arrayModelTypeName,
			ArrayGenericTypeName: arrayGenericTypeName,
		})

		validateData := ModelFieldValidatorGenerateData{
			Field:         field.Name,
			IsModel:       isModelType,
			IsArray:       isArrayType,
			IsGenericType: isGenericType,
			TypeName:      typeName,
			Rules:         []ModelFieldValidateRuleGenerateData{},
		}

		for _, validator := range field.Validators {

			var convertedArgs []string
			for _, arg := range validator.Args {
				convertedArgs = append(convertedArgs, emitValueLiteral(arg))
			}

			validateData.Rules = append(validateData.Rules, ModelFieldValidateRuleGenerateData{
				Name: validator.Name,
				Args: convertedArgs,
			})
		}

		validators = append(validators, validateData)
	}

	data := ModelGenerateData{
		Name:       ir.Name,
		ModelHash:  hashModelName(ir.Name),
		IsGeneric:  len(ir.TypeParams) > 0,
		TypeParams: ir.TypeParams,
		Fields:     fields,
		Validators: validators,
	}

	err := tmpl.ExecuteTemplate(&wr, "model.tmpl", data)
	if err != nil {
		fmt.Println(err)
		return "", exception.NewEmitException("Error", ir.Span.ToLocation())
	}

	return wr.String(), nil
}

func (t *TypescriptEmitter) EmitRest(tmpl *template.Template, ir *generator.RestEndpointIR) (string, *exception.EmitException) {
	var wr bytes.Buffer

	ReqType := "any"
	ResType := "any"
	var err = new(exception.EmitException)
	if ir.RequestBodyType != nil {
		ReqType, err = t.EmitTypeName(ir.RequestBodyType)
		if err != nil {
			return "", err
		}
	}

	if ir.ResponseBodyType != nil {
		ResType, err = t.EmitTypeName(ir.ResponseBodyType)
		if err != nil {
			return "", err
		}
	}

	data := RestGenerateData{
		Name:     ir.Name,
		Path:     ir.Path,
		Method:   ir.Method,
		ResType:  ReqType,
		ReqType:  ResType,
		ReqClass: ReqType,
		ResClass: ResType,
	}

	err1 := tmpl.ExecuteTemplate(&wr, "rest.tmpl", data)
	if err1 != nil {
		fmt.Println(err1)
		return "", exception.NewEmitException("Error", ir.Span.ToLocation())
	}

	return wr.String(), nil
}

func (t *TypescriptEmitter) Emit(ir *generator.ProgramIR) (string, exception.IException) {
	var sb strings.Builder
	tmpl, err := template.ParseFS(templateFiles, "templates/*.tmpl")
	if err != nil {
		return "", exception.NewEmitException(err.Error(), nil)
	}

	sb.WriteString("// @ts-nocheck\n")
	sb.WriteString("import { Validator, ContractBaseError } from \"contractor-ts\";\n\n")
	sb.WriteString("import type { GeneratedErrorConstructorMap, GeneratedValidationDetails, EventMetadata, EventPayload, RestMetadata, RestRequestBody, RestResponseBody } from \"contractor-ts\";\n\n")

	// if len(ir.Errors) > 0 {
	// 	for _, errorIR := range ir.Errors {
	// 		code, err := t.EmitError(tmpl, errorIR)
	// 		if err != nil {
	// 			return "", err
	// 		}

	// 		sb.WriteString(code)
	// 	}

	// 	sb.WriteString("export const errorConstructorsByCode: GeneratedErrorConstructorMap = {\n")
	// 	for _, errorIR := range ir.Errors {
	// 		codeLiteral := quoteLiteral(errorIR.Code, errorIR.Name)
	// 		sb.WriteString("  ")
	// 		sb.WriteString(codeLiteral)
	// 		sb.WriteString(": ")
	// 		sb.WriteString(errorIR.Name)
	// 		sb.WriteString(",\n")
	// 	}
	// 	sb.WriteString("};\n\n")
	// }

	for _, model := range ir.Models {
		code, err := t.EmitModel(tmpl, model)
		if err != nil {
			return "", err
		}

		sb.WriteString(code)
	}

	// for _, enumItem := range ir.Enums {
	// 	code, err := t.EmitEnum(tmpl, enumItem)
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	sb.WriteString(code)
	// }

	// for _, eventItem := range ir.Events {
	// 	code, err := t.EmitEvent(tmpl, eventItem)
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	sb.WriteString(code)
	// }

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

func quoteLiteral(value *string, fallback string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return strconv.Quote(fallback)
	}

	return strconv.Quote(*value)
}
