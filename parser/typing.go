package parser

import (
	"reflect"
	"slices"
	"strconv"
)

type Annotation struct {
	TypeName []string
	Args     int
}

var SupportedFieldValidatorAnnotations = map[string]Annotation{
	"IsInt":     {[]string{"Any"}, 1},
	"IsFloat":   {[]string{"Any"}, 1},
	"IsBoolean": {[]string{"Any"}, 1},
	"IsString":  {[]string{"Any"}, 1},
	"IsArray":   {[]string{"Any"}, 1},

	"IsEmail":       {[]string{"String"}, 1},
	"IsNotEmpty":    {[]string{"String"}, 1},
	"IsPhoneNumber": {[]string{"String"}, 1},
	"IsDateString":  {[]string{"String"}, 1},
	"IsUUID":        {[]string{"String"}, 1},
	"IsUrl":         {[]string{"String"}, 1},

	"Max":       {[]string{"Number"}, 2},
	"Min":       {[]string{"Number"}, 2},
	"MinLength": {[]string{"String"}, 2},
	"MaxLength": {[]string{"String"}, 2},
	"Length":    {[]string{"String"}, 2},

	"ArrayMinSize": {[]string{"Array"}, 2},
	"ArrayMaxSize": {[]string{"Array"}, 2},
	"ArrayLength":  {[]string{"Array"}, 2},
}

var SupportedFieldAnnotations = map[string]Annotation{
	"Optional": {[]string{}, 0},
}

var SupportedModelAnnotations = map[string]Annotation{
	"CreateConstructor": {[]string{}, 0},
	"Mapper":            {[]string{}, 0},
	"Data":              {[]string{}, 0},
}

type SymbolType int

const (
	PRIMITIVE_TYPE SymbolType = iota
	DEFINED_TYPE
	TYPE_VAR
)

type Symbol struct {
	Type       SymbolType
	Name       string
	Node       Node
	HasGeneric bool
}

func NewSymbol(t SymbolType, name string, node Node, hasGeneric bool) *Symbol {
	return &Symbol{
		Type:       t,
		Name:       name,
		HasGeneric: hasGeneric,
		Node:       node,
	}
}

func isTrulyNil(i interface{}) bool {
	if i == nil {
		return true
	}

	val := reflect.ValueOf(i)
	// If it's a pointer, slice, map, or func, it can be nil
	switch val.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Interface, reflect.Func:
		return val.IsNil()
	}
	return false
}

type Context struct {
	Parent  *Context
	Symbols map[string]*Symbol
}

func (c *Context) CreateChildContext() *Context {
	ctx := NewContext()
	ctx.Parent = c
	return ctx
}

func (c *Context) FindSymbol(name string, recursion bool, node Node) (*Symbol, BaseError) {
	symbol, ok := c.Symbols[name]
	if !ok {
		if recursion && c.Parent != nil {
			return c.Parent.FindSymbol(name, recursion, node)
		}
	} else {
		return symbol, nil
	}

	loc := node.GetLocation()
	return nil, NewTypeError(
		"Cannot find symbol "+name+" in scope",
		loc.GetErrorLocation(),
	)
}

func NewContext() *Context {
	return &Context{
		Symbols: make(map[string]*Symbol),
	}
}

func CreateGlobalContext() *Context {
	ctx := NewContext()
	ctx.Symbols["Number"] = NewSymbol(PRIMITIVE_TYPE, "Number", nil, false)
	ctx.Symbols["Bool"] = NewSymbol(PRIMITIVE_TYPE, "Bool", nil, false)
	ctx.Symbols["Unknown"] = NewSymbol(PRIMITIVE_TYPE, "Unknown", nil, false)
	ctx.Symbols["String"] = NewSymbol(PRIMITIVE_TYPE, "String", nil, false)
	ctx.Symbols["Array"] = NewSymbol(PRIMITIVE_TYPE, "Array", nil, true)
	ctx.Symbols["Object"] = NewSymbol(PRIMITIVE_TYPE, "Object", nil, false)
	return ctx
}

var globalCtx = CreateGlobalContext()

type TypeChecker struct{}

func (c *TypeChecker) CheckFieldAnnotation(node *ModelFieldNode) BaseError {

	if node.Annotations == nil {
		return nil
	}

	for _, anno := range node.Annotations.List {
		if validator, ok := SupportedFieldValidatorAnnotations[anno.Name]; ok {
			if len(validator.TypeName) > 0 && !slices.Contains(validator.TypeName, "Any") && !slices.Contains(validator.TypeName, node.Type.Name) {
				return NewTypeError(
					"Validator "+anno.Name+" cannot be applied to type "+node.Type.Name,
					anno.Loc.GetErrorLocation(),
				)
			}

			if len(anno.Args) < validator.Args {
				return NewTypeError(
					"Validator "+anno.Name+" requires "+strconv.Itoa(validator.Args)+" arguments",
					anno.Loc.GetErrorLocation(),
				)
			}
		} else if fieldAnno, ok := SupportedFieldAnnotations[anno.Name]; ok {
			if len(fieldAnno.TypeName) > 0 && !slices.Contains(fieldAnno.TypeName, "Any") && !slices.Contains(fieldAnno.TypeName, node.Type.Name) {
				return NewTypeError(
					"Annotation "+anno.Name+" cannot be applied to type "+node.Type.Name,
					anno.Loc.GetErrorLocation(),
				)
			}

			if len(anno.Args) < fieldAnno.Args {
				return NewTypeError(
					"Validator "+anno.Name+" requires "+strconv.Itoa(validator.Args)+" arguments",
					anno.Loc.GetErrorLocation(),
				)
			}
		} else {
			return NewTypeError(
				"Unsupported annotation "+anno.Name+" on field "+node.Name,
				anno.Loc.GetErrorLocation(),
			)
		}
	}

	return nil
}

func (c *TypeChecker) CheckModelAnnotation(node *ModelStatementNode) BaseError {

	if node.Annotations == nil {
		return nil
	}

	for _, anno := range node.Annotations.List {
		if modelAnno, ok := SupportedModelAnnotations[anno.Name]; ok {
			if len(modelAnno.TypeName) > 0 && !slices.Contains(modelAnno.TypeName, "Any") && !slices.Contains(modelAnno.TypeName, node.Name) {
				return NewTypeError(
					"Annotation "+anno.Name+" cannot be applied to model "+node.Name,
					anno.Loc.GetErrorLocation(),
				)
			}

			if len(anno.Args) < modelAnno.Args {
				return NewTypeError(
					"Validator "+anno.Name+" requires "+strconv.Itoa(modelAnno.Args)+" arguments",
					anno.Loc.GetErrorLocation(),
				)
			}
		} else {
			return NewTypeError(
				"Unsupported annotation "+anno.Name+" on model "+node.Name,
				anno.Loc.GetErrorLocation(),
			)
		}
	}

	return nil
}

func (c *TypeChecker) CheckType(ctx *Context, t *TypeDeclarationNode) BaseError {
	if t == nil {
		return nil
	}

	symbol, err := ctx.FindSymbol(t.Name, true, t)
	if err != nil {
		return err
	}
	if !symbol.HasGeneric && !isNil(t.Generic) {
		return NewTypeError("Symbol "+t.Name+" is not a generic type", t.Loc.GetErrorLocation())
	} else if symbol.HasGeneric && isNil(t.Generic) {
		return NewTypeError("Symbol "+t.Name+" requires a generic type argument", t.Loc.GetErrorLocation())
	} else if symbol.HasGeneric && !isNil(t.Generic) {
		return c.CheckType(ctx, t.Generic.(*TypeDeclarationNode))
	}

	return nil
}

func (c *TypeChecker) CheckModelField(parentCtx *Context, node *ModelFieldNode) BaseError {
	err := c.CheckType(parentCtx, node.Type)
	if err != nil {
		return err
	}

	err = c.CheckFieldAnnotation(node)
	if err != nil {
		return err
	}

	return nil
}

func (c *TypeChecker) CheckModelStatement(parentCtx *Context, node *ModelStatementNode) BaseError {
	ctx := parentCtx.CreateChildContext()

	symbol, _ := parentCtx.FindSymbol(node.Name, true, node)
	if symbol != nil {
		return NewTypeError(
			"Re-declare symbol "+node.Name+" in scope",
			node.Loc.GetErrorLocation(),
		)
	}

	parentCtx.Symbols[node.Name] = NewSymbol(DEFINED_TYPE, node.Name, node, false)

	if node.TypeVar != nil {
		symbol, _ = parentCtx.FindSymbol(node.TypeVar.Name, true, node)
		if symbol != nil {
			return NewTypeError(
				"Cannot use symbol "+node.TypeVar.Name+" as type var in scope",
				node.TypeVar.Loc.GetErrorLocation(),
			)
		}
		ctx.Symbols[node.TypeVar.Name] = NewSymbol(TYPE_VAR, node.TypeVar.Name, node.TypeVar, false)
		parentCtx.Symbols[node.Name].HasGeneric = true
	}

	for _, fieldNode := range node.Fields {
		err := c.CheckModelField(ctx, fieldNode)
		if err != nil {
			return err
		}
	}

	err := c.CheckModelAnnotation(node)
	if err != nil {
		return err
	}

	return nil
}

func (c *TypeChecker) Check(ast *AST) BaseError {
	ctx := NewContext()
	ctx.Parent = globalCtx
	for _, node := range ast.Statements {
		switch v := node.(type) {
		case *ModelStatementNode:
			err := c.CheckModelStatement(ctx, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewTypeChecker() *TypeChecker {
	return &TypeChecker{}
}
