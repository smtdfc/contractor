package golang

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/smtdfc/contractor/exception"
	"github.com/smtdfc/contractor/generator"
)

func hashModelName(name string) string {
	hasher := md5.New()
	hasher.Write([]byte(name))
	return fmt.Sprintf("%x", hasher.Sum(nil))[:8]
}

type GoEmitter struct{}

var typeMap = map[string]string{
	"Int":    "int",
	"Float":  "float64",
	"String": "string",
	"Bool":   "bool",
	"Null":   "interface{}",
	"Any":    "interface{}",
	"Array":  "[]",
}

func (t *GoEmitter) EmitTypeName(ir *generator.TypeIR) (string, *exception.EmitException) {
	var typeName strings.Builder

	if ir.Kind == generator.TypeKindBuiltin {
		if ir.Name == "Array" {
			typeName.WriteString("[]")
			if len(ir.Generics) > 0 {
				innerType, err := t.EmitTypeName(ir.Generics[0])
				if err != nil {
					return "", err
				}
				typeName.WriteString(innerType)
			} else {
				typeName.WriteString("interface{}")
			}
			return typeName.String(), nil
		}

		goType, ok := typeMap[ir.Name]
		if !ok {
			goType = "interface{}"
		}
		typeName.WriteString(goType)
	}

	if ir.Kind == generator.TypeKindModel || ir.Kind == generator.TypeKindEnum || ir.Kind == generator.TypeKindGeneric {
		typeName.WriteString(ir.Name)
	}

	if len(ir.Generics) > 0 && ir.Kind != generator.TypeKindBuiltin {
		typeName.WriteString("[")
		genericTypes := []string{}

		for _, generic := range ir.Generics {
			goGenericType, err := t.EmitTypeName(generic)
			if err != nil {
				return "", err
			}
			genericTypes = append(genericTypes, goGenericType)
		}

		typeName.WriteString(strings.Join(genericTypes, ", "))
		typeName.WriteString("]")
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
	NameTitle            string
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
	FieldTitle    string
	IsOptional    bool
	FieldType     string
	Rules         []ModelFieldValidateRuleGenerateData
}

type RestGenerateData struct {
	Name   string
	Path   string
	Method string
}

type ErrorGenerateData struct {
	Name     string
	Code     string
	Message  string
	Scope    string
	Status   string
	HasScope bool
}

type EnumMemberGenerateData struct {
	Key    string
	Value  string
	IsLast bool
}

type EnumGenerateData struct {
	Name    string
	Members []EnumMemberGenerateData
}

type EventGenerateData struct {
	MetadataConst string
	EventNameLit  string
	PayloadAlias  string
	PayloadType   string
}

func titleCase(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func (t *GoEmitter) EmitModel(tmpl *template.Template, ir *generator.ModelIR) (string, *exception.EmitException) {
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
		typeName := ""

		if field.Type.Kind == generator.TypeKindBuiltin && field.Type.Name == "Array" {
			isFieldArrayType = true
			if len(field.Type.Generics) > 0 {
				innerType := field.Type.Generics[0]
				if innerType.Kind == generator.TypeKindModel {
					arrayModelTypeName = innerType.Name
					isFieldModelType = true
					typeName = innerType.Name
				}
				if innerType.Kind == generator.TypeKindGeneric {
					arrayGenericTypeName = innerType.Name
				}
			}
		}
		if field.Type.Kind == generator.TypeKindModel {
			isFieldModelType = true
			modelTypeName = field.Type.Name
			typeName = field.Type.Name
		}
		if field.Type.Kind == generator.TypeKindGeneric {
			isFieldGenericType = true
			genericTypeName = field.Type.Name
		}

		fields = append(fields, ModelFieldGenerateData{
			Name:                 field.Name,
			NameTitle:            titleCase(field.Name),
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
			IsModel:       isFieldModelType,
			IsArray:       isFieldArrayType,
			IsGenericType: isFieldGenericType,
			TypeName:      typeName,
			FieldTitle:    titleCase(field.Name),
			IsOptional:    field.IsOptional,
			FieldType:     fieldType,
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

func (t *GoEmitter) EmitError(tmpl *template.Template, ir *generator.ErrorIR) (string, *exception.EmitException) {
	var wr bytes.Buffer

	data := ErrorGenerateData{
		Name:     ir.Name,
		Code:     quoteLiteral(ir.Code, ir.Name),
		Message:  strconv.Quote(ir.Message),
		Scope:    quoteLiteral(ir.Scope, ir.Name),
		Status:   renderStatusLiteral(ir.Status),
		HasScope: ir.Scope != nil && strings.TrimSpace(*ir.Scope) != "",
	}

	err := tmpl.ExecuteTemplate(&wr, "error.tmpl", data)
	if err != nil {
		fmt.Println(err)
		return "", exception.NewEmitException("Error", ir.Span.ToLocation())
	}

	return wr.String(), nil
}

func (t *GoEmitter) EmitEnum(tmpl *template.Template, ir *generator.EnumIR) (string, *exception.EmitException) {
	var wr bytes.Buffer

	members := make([]EnumMemberGenerateData, 0, len(ir.Members))
	for i, member := range ir.Members {
		members = append(members, EnumMemberGenerateData{
			Key:    member,
			Value:  member,
			IsLast: i == len(ir.Members)-1,
		})
	}

	data := EnumGenerateData{
		Name:    ir.Name,
		Members: members,
	}

	err := tmpl.ExecuteTemplate(&wr, "enum.tmpl", data)
	if err != nil {
		fmt.Println(err)
		return "", exception.NewEmitException("Error", ir.Span.ToLocation())
	}

	return wr.String(), nil
}

func (t *GoEmitter) EmitEvent(tmpl *template.Template, ir *generator.EventIR) (string, *exception.EmitException) {
	var wr bytes.Buffer

	payloadType, err := t.EmitTypeName(ir.PayloadType)
	if err != nil {
		return "", err
	}

	data := EventGenerateData{
		MetadataConst: ir.Name,
		EventNameLit:  strconv.Quote(ir.EventName),
		PayloadAlias:  ir.Name + "Payload",
		PayloadType:   payloadType,
	}

	execErr := tmpl.ExecuteTemplate(&wr, "event.tmpl", data)
	if execErr != nil {
		fmt.Println(execErr)
		return "", exception.NewEmitException("Error", ir.Span.ToLocation())
	}

	return wr.String(), nil
}

func (t *GoEmitter) EmitRest(tmpl *template.Template, ir *generator.RestEndpointIR) (string, *exception.EmitException) {
	var wr bytes.Buffer

	data := RestGenerateData{
		Name:   ir.Name,
		Path:   strconv.Quote(ir.Path),
		Method: strconv.Quote(ir.Method),
	}

	err := tmpl.ExecuteTemplate(&wr, "rest.tmpl", data)
	if err != nil {
		fmt.Println(err)
		return "", exception.NewEmitException("Error", ir.Span.ToLocation())
	}

	return wr.String(), nil
}

func (t *GoEmitter) Emit(ir *generator.ProgramIR) (string, exception.IException) {
	var sb strings.Builder
	tmpl, err := template.ParseFS(templateFiles, "templates/*.tmpl")
	if err != nil {
		return "", exception.NewEmitException(err.Error(), nil)
	}

	sb.WriteString("// Code generated by contractor. DO NOT EDIT.\n")
	sb.WriteString("package generated\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"github.com/smtdfc/contractor/lib/golang/core\"\n")
	sb.WriteString(")\n\n")

	for _, model := range ir.Models {
		code, err := t.EmitModel(tmpl, model)
		if err != nil {
			return "", err
		}

		sb.WriteString(code)
		sb.WriteString("\n")
	}

	for _, errorIR := range ir.Errors {
		code, err := t.EmitError(tmpl, errorIR)
		if err != nil {
			return "", err
		}

		sb.WriteString(code)
		sb.WriteString("\n")
	}

	for _, enumItem := range ir.Enums {
		code, err := t.EmitEnum(tmpl, enumItem)
		if err != nil {
			return "", err
		}

		sb.WriteString(code)
		sb.WriteString("\n")
	}

	for _, eventItem := range ir.Events {
		code, err := t.EmitEvent(tmpl, eventItem)
		if err != nil {
			return "", err
		}

		sb.WriteString(code)
		sb.WriteString("\n")
	}

	for _, rest := range ir.Rests {
		code, err := t.EmitRest(tmpl, rest)
		if err != nil {
			return "", err
		}

		sb.WriteString(code)
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func NewGoEmitter() *GoEmitter {
	return &GoEmitter{}
}

func quoteLiteral(value *string, fallback string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return strconv.Quote(fallback)
	}

	return strconv.Quote(*value)
}

func renderStatusLiteral(value *string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return ""
	}

	trimmed := strings.TrimSpace(*value)
	if _, err := strconv.Atoi(trimmed); err == nil {
		return trimmed
	}

	return ""
}

func emitValueLiteral(value *generator.ValueIR) string {
	if value == nil {
		return "nil"
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
		return "nil"
	case "Array":
		rawValues, ok := value.Value.([]*generator.ValueIR)
		if !ok {
			return "[]interface{}{}"
		}

		items := make([]string, 0, len(rawValues))
		for _, item := range rawValues {
			items = append(items, emitValueLiteral(item))
		}

		return "[]interface{}{" + strings.Join(items, ", ") + "}"
	default:
		if value.Value == nil {
			return "nil"
		}
		return fmt.Sprint(value.Value)
	}
}
