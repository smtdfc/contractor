package generator

import (
	"strings"

	"github.com/smtdfc/contractor/exception"
	"github.com/smtdfc/contractor/parser"
)

var validatorAnnotations = map[string]struct{}{
	"Is":             {},
	"Min":            {},
	"Max":            {},
	"Length":         {},
	"MinLength":      {},
	"MaxLength":      {},
	"Range":          {},
	"Matches":        {},
	"Contains":       {},
	"StartsWith":     {},
	"EndsWith":       {},
	"In":             {},
	"IsEmail":        {},
	"IsNumber":       {},
	"IsURL":          {},
	"IsUUID":         {},
	"IsDate":         {},
	"IsDateTime":     {},
	"IsAlpha":        {},
	"IsAlnum":        {},
	"NotNull":        {},
	"IsBool":         {},
	"IsModel":        {},
	"NestedValidate": {},
}

type IRGenerator struct {
	builtinTypes map[string]struct{}
}

func NewIRGenerator() *IRGenerator {
	return &IRGenerator{
		builtinTypes: map[string]struct{}{
			"String": {},
			"Int":    {},
			"Float":  {},
			"Bool":   {},
			"Array":  {},
			"Null":   {},
			"Any":    {},
		},
	}
}

func (g *IRGenerator) GenerateProgram(ast *parser.ProgramNode) (*ProgramIR, exception.IException) {
	if ast == nil {
		return &ProgramIR{Errors: make([]*ErrorIR, 0), Models: make([]*ModelIR, 0), Rests: make([]*RestEndpointIR, 0)}, nil
	}

	modelSymbols, err := g.collectModelSymbols(ast)
	if err != nil {
		return nil, err
	}

	errors := make([]*ErrorIR, 0)
	models := make([]*ModelIR, 0)
	rests := make([]*RestEndpointIR, 0)

	for _, node := range ast.Body {
		switch v := node.(type) {
		case *parser.ModelDeclNode:
			modelIR, err := g.modelToIR(v, modelSymbols)
			if err != nil {
				return nil, err
			}

			models = append(models, modelIR)
		case *parser.ErrorDeclNode:
			errorIR, err := g.errorToIR(v)
			if err != nil {
				return nil, err
			}

			errors = append(errors, errorIR)
		case *parser.RestDeclNode:
			restIR, err := g.restToIR(v, modelSymbols)
			if err != nil {
				return nil, err
			}

			rests = append(rests, restIR)
		}
	}

	return &ProgramIR{
		Errors: errors,
		Models: models,
		Rests:  rests,
	}, nil
}

func (g *IRGenerator) collectModelSymbols(ast *parser.ProgramNode) (map[string]struct{}, exception.IException) {
	result := make(map[string]struct{})

	for _, node := range ast.Body {
		modelNode, ok := node.(*parser.ModelDeclNode)
		if !ok {
			continue
		}

		if modelNode.Name == nil {
			return nil, exception.NewTypeException("Model name is missing", modelNode.Loc)
		}

		if _, exists := result[modelNode.Name.Value]; exists {
			return nil, exception.NewTypeException("Model '"+modelNode.Name.Value+"' is already defined", modelNode.Name.Loc)
		}

		result[modelNode.Name.Value] = struct{}{}
	}

	return result, nil
}

func (g *IRGenerator) modelToIR(node *parser.ModelDeclNode, modelSymbols map[string]struct{}) (*ModelIR, exception.IException) {
	if node == nil {
		fallbackLoc := parser.NewLocation("<unknown>", parser.NewPosition(1, 1), parser.NewPosition(1, 1))
		return nil, exception.NewTypeException("Model name is missing", fallbackLoc)
	}

	if node.Name == nil {
		return nil, exception.NewTypeException("Model name is missing", node.GetLocation())
	}

	genericSymbols := make(map[string]struct{})
	typeParams := make([]string, 0, len(node.Generics))
	for _, item := range node.Generics {
		if item == nil || item.Name == nil {
			continue
		}

		typeParams = append(typeParams, item.Name.Value)
		genericSymbols[item.Name.Value] = struct{}{}
	}

	fields := make([]*ModelField, 0, len(node.Fields))
	for _, field := range node.Fields {
		if field == nil || field.Name == nil {
			continue
		}

		validators, err := g.extractValidator(field)
		if err != nil {
			return nil, err
		}

		fieldType := g.typeToIR(field.Type, modelSymbols, genericSymbols)
		fieldIR := &ModelField{
			Span:        toSourceSpan(field.Loc),
			Name:        field.Name.Value,
			Annotations: g.annotationsToIR(field.Annotations),
			Type:        fieldType,
			IsOptional:  field.Optional,
			Validators:  validators,
		}

		fields = append(fields, fieldIR)
	}

	annotations := g.annotationsToIR(node.Annotations)
	return &ModelIR{
		Span:                toSourceSpan(node.Loc),
		Name:                node.Name.Value,
		TypeParams:          typeParams,
		Annotations:         annotations,
		Fields:              fields,
		IsCreateConstructor: hasAnnotation(annotations, "CreateConstructor", "Constructor"),
		IsCreateMapper:      hasAnnotation(annotations, "CreateMapper", "Mapper", "Mapping"),
	}, nil
}

func (g *IRGenerator) extractValidator(node *parser.ModelFieldDeclNode) ([]*FieldValidator, exception.IException) {

	validators := []*FieldValidator{}

	for _, anno := range node.Annotations {
		if _, ok := validatorAnnotations[anno.Name.Value]; ok {
			args := []*ValueIR{}
			for _, arg := range anno.Args {
				value := g.valueToIR(arg)
				args = append(args, value)
			}

			validators = append(validators, &FieldValidator{
				Name: anno.Name.Value,
				Args: args,
			})
		}
	}

	return validators, nil
}

func (g *IRGenerator) restToIR(node *parser.RestDeclNode, modelSymbols map[string]struct{}) (*RestEndpointIR, exception.IException) {
	if node == nil {
		fallbackLoc := parser.NewLocation("<unknown>", parser.NewPosition(1, 1), parser.NewPosition(1, 1))
		return nil, exception.NewTypeException("Rest name is missing", fallbackLoc)
	}

	if node.Name == nil {
		return nil, exception.NewTypeException("Rest name is missing", node.GetLocation())
	}

	methodNode, ok := node.MethodValue.(*parser.StringValueNode)
	if !ok || methodNode == nil {
		return nil, exception.NewTypeException("Rest property 'method' must be a string literal", node.GetLocation())
	}

	pathNode, ok := node.PathValue.(*parser.StringValueNode)
	if !ok || pathNode == nil {
		return nil, exception.NewTypeException("Rest property 'path' must be a string literal", node.GetLocation())
	}

	queries, err := g.queriesToIR(node.QueriesValue)
	if err != nil {
		return nil, err
	}

	requestType := g.typeToIR(node.RequestBodyType, modelSymbols, map[string]struct{}{})
	if node.RequestBodyType == nil {
		requestType = nil
	}

	responseType := g.typeToIR(node.ResponseBodyType, modelSymbols, map[string]struct{}{})
	if node.ResponseBodyType == nil {
		responseType = nil
	}

	return &RestEndpointIR{
		Span:             toSourceSpan(node.Loc),
		Name:             node.Name.Value,
		Method:           strings.ToUpper(methodNode.Value),
		Path:             pathNode.Value,
		RequestBodyType:  requestType,
		ResponseBodyType: responseType,
		Queries:          queries,
	}, nil
}

func (g *IRGenerator) errorToIR(node *parser.ErrorDeclNode) (*ErrorIR, exception.IException) {
	if node == nil {
		fallbackLoc := parser.NewLocation("<unknown>", parser.NewPosition(1, 1), parser.NewPosition(1, 1))
		return nil, exception.NewTypeException("Error name is missing", fallbackLoc)
	}

	if node.Name == nil {
		return nil, exception.NewTypeException("Error name is missing", node.GetLocation())
	}

	message := ""
	if msg, ok := node.MessageValue.(*parser.StringValueNode); ok && msg != nil {
		message = msg.Value
	}

	return &ErrorIR{
		Span:    toSourceSpan(node.Loc),
		Name:    node.Name.Value,
		Code:    stringValueOrNil(node.CodeValue),
		Message: message,
		Scope:   stringValueOrNil(node.ScopeValue),
		Status:  statusValueOrNil(node.StatusValue),
	}, nil
}

func (g *IRGenerator) queriesToIR(node parser.ASTValueNode) ([]string, exception.IException) {
	if node == nil {
		return make([]string, 0), nil
	}

	arrayNode, ok := node.(*parser.ArrayValueNode)
	if !ok {
		return nil, exception.NewTypeException("Rest property 'queries' must be an array literal", node.GetLocation())
	}

	queries := make([]string, 0, len(arrayNode.Values))
	for _, item := range arrayNode.Values {
		strNode, ok := item.(*parser.StringValueNode)
		if !ok {
			return nil, exception.NewTypeException("Rest property 'queries' must be string[]", item.GetLocation())
		}

		queries = append(queries, strNode.Value)
	}

	return queries, nil
}

func (g *IRGenerator) typeToIR(node *parser.TypeDeclNode, modelSymbols map[string]struct{}, genericSymbols map[string]struct{}) *TypeIR {
	if node == nil || node.Name == nil {
		return &TypeIR{Kind: TypeKindUnknown}
	}

	kind := TypeKindUnknown
	name := node.Name.Value

	if _, ok := genericSymbols[name]; ok {
		kind = TypeKindGeneric
	} else if _, ok := g.builtinTypes[name]; ok {
		kind = TypeKindBuiltin
	} else if _, ok := modelSymbols[name]; ok {
		kind = TypeKindModel
	}

	generics := make([]*TypeIR, 0, len(node.Generics))
	for _, item := range node.Generics {
		generics = append(generics, g.typeToIR(item, modelSymbols, genericSymbols))
	}

	resolvedRef := ""
	if kind == TypeKindModel {
		resolvedRef = name
	}

	return &TypeIR{
		Span:        toSourceSpan(node.Loc),
		Kind:        kind,
		Name:        name,
		Generics:    generics,
		ResolvedRef: resolvedRef,
	}
}

func (g *IRGenerator) annotationsToIR(nodes []*parser.AnnotationNode) []*AnnotationIR {
	items := make([]*AnnotationIR, 0, len(nodes))
	for _, node := range nodes {
		if node == nil || node.Name == nil {
			continue
		}

		args := make([]*ValueIR, 0, len(node.Args))
		for _, arg := range node.Args {
			args = append(args, g.valueToIR(arg))
		}

		items = append(items, &AnnotationIR{
			Name: node.Name.Value,
			Args: args,
		})
	}

	return items
}

func (g *IRGenerator) valueToIR(node parser.ASTValueNode) *ValueIR {
	if node == nil {
		return &ValueIR{Kind: "Null", Value: "null"}
	}

	switch v := node.(type) {
	case *parser.StringValueNode:
		return &ValueIR{Kind: v.GetKind(), Value: v.Value}
	case *parser.NumberValueNode:
		return &ValueIR{Kind: v.GetKind(), Value: v.Value}
	case *parser.BooleanValueNode:
		return &ValueIR{Kind: v.GetKind(), Value: v.Value}
	case *parser.NullValueNode:
		return &ValueIR{Kind: v.GetKind(), Value: "null"}
	case *parser.ArrayValueNode:
		values := make([]*ValueIR, 0, len(v.Values))
		for _, item := range v.Values {
			values = append(values, g.valueToIR(item))
		}
		return &ValueIR{Kind: v.GetKind(), Value: values}
	default:
		return &ValueIR{Kind: node.GetKind(), Value: nil}
	}
}

func stringValueOrNil(node parser.ASTValueNode) *string {
	if node == nil {
		return nil
	}

	value, ok := node.(*parser.StringValueNode)
	if !ok || value == nil {
		return nil
	}

	result := value.Value
	return &result
}

func statusValueOrNil(node parser.ASTValueNode) *string {
	if node == nil {
		return nil
	}

	switch value := node.(type) {
	case *parser.StringValueNode:
		if value == nil {
			return nil
		}
		result := value.Value
		return &result
	case *parser.NumberValueNode:
		if value == nil {
			return nil
		}
		result := value.Value
		return &result
	default:
		return nil
	}
}

func toSourceSpan(loc *parser.Location) *SourceSpan {
	if loc == nil || loc.Start == nil || loc.End == nil {
		return nil
	}

	return &SourceSpan{
		File:      loc.File,
		StartLine: loc.Start.Line,
		StartCol:  loc.Start.Col,
		EndLine:   loc.End.Line,
		EndCol:    loc.End.Col,
	}
}

func hasAnnotation(annotations []*AnnotationIR, names ...string) bool {
	if len(annotations) == 0 || len(names) == 0 {
		return false
	}

	match := make(map[string]struct{}, len(names))
	for _, name := range names {
		match[strings.ToLower(name)] = struct{}{}
	}

	for _, annotation := range annotations {
		if _, ok := match[strings.ToLower(annotation.Name)]; ok {
			return true
		}
	}

	return false
}
