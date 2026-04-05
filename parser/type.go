package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/smtdfc/contractor/exception"
)

type Symbol interface {
	GetName() string
	GetGenerics() []*TypeVarNode
	GetKind() string
	BuiltIn() bool
}

type TypeSymbol struct {
	Name      string
	IsBuiltIn bool
	Generics  []*TypeVarNode
}

type RestSymbol struct {
	Name string
}

type ErrorSymbol struct {
	Name string
}

func (s *TypeSymbol) GetName() string {
	return s.Name
}

func (s *TypeSymbol) GetGenerics() []*TypeVarNode {
	if s == nil {
		return nil
	}
	return s.Generics
}

func (s *TypeSymbol) GetKind() string {
	return "type"
}

func (s *TypeSymbol) BuiltIn() bool {
	if s == nil {
		return false
	}
	return s.IsBuiltIn
}

func NewTypeSymbol(name string, isBuiltIn bool) *TypeSymbol {
	return &TypeSymbol{Name: name, IsBuiltIn: isBuiltIn, Generics: make([]*TypeVarNode, 0)}
}

func (s *RestSymbol) GetName() string {
	return s.Name
}

func (s *RestSymbol) GetGenerics() []*TypeVarNode {
	return nil
}

func (s *RestSymbol) GetKind() string {
	return "rest"
}

func (s *RestSymbol) BuiltIn() bool {
	return false
}

func NewRestSymbol(name string) *RestSymbol {
	return &RestSymbol{Name: name}
}

func (s *ErrorSymbol) GetName() string {
	return s.Name
}

func (s *ErrorSymbol) GetGenerics() []*TypeVarNode {
	return nil
}

func (s *ErrorSymbol) GetKind() string {
	return "error"
}

func (s *ErrorSymbol) BuiltIn() bool {
	return false
}

func NewErrorSymbol(name string) *ErrorSymbol {
	return &ErrorSymbol{Name: name}
}

type AnnotationSymbol struct {
	Name      string
	IsBuiltIn bool
	Generics  []*TypeVarNode
	Args      map[string]*TypeDeclNode
	ArgOrder  []string
}

func (s *AnnotationSymbol) GetName() string {
	return s.Name
}

func (s *AnnotationSymbol) GetGenerics() []*TypeVarNode {
	if s == nil {
		return nil
	}
	return s.Generics
}

func (s *AnnotationSymbol) GetKind() string {
	return "annotation"
}

func (s *AnnotationSymbol) BuiltIn() bool {
	if s == nil {
		return false
	}
	return s.IsBuiltIn
}

func NewAnnotationSymbol(name string, isBuiltIn bool) *AnnotationSymbol {
	return &AnnotationSymbol{
		Name:      name,
		IsBuiltIn: isBuiltIn,
		Generics:  make([]*TypeVarNode, 0),
		Args:      make(map[string]*TypeDeclNode),
		ArgOrder:  make([]string, 0),
	}
}

type Context struct {
	Parent  *Context
	Symbols map[string]Symbol
}

func (c *Context) Add(sym Symbol) {
	c.Symbols[sym.GetName()] = sym
}

func (c *Context) Find(sym Symbol) bool {
	if sym == nil {
		return false
	}

	if c.Symbols[sym.GetName()] != nil {
		return true
	}

	if c.Parent != nil {
		return c.Parent.Find(sym)
	}

	return false
}

func (c *Context) GetByName(name string) *Symbol {
	if c.Symbols[name] != nil {
		sym := c.Symbols[name]
		return &sym
	}

	if c.Parent != nil {
		return c.Parent.GetByName(name)
	}

	return nil
}

func (c *Context) GetTypeByName(name string) *TypeSymbol {
	sym := c.GetByName(name)
	if sym == nil {
		return nil
	}

	t, ok := (*sym).(*TypeSymbol)
	if !ok {
		return nil
	}

	return t
}

func (c *Context) GetAnnotationByName(name string) *AnnotationSymbol {
	sym := c.GetByName(name)
	if sym == nil {
		return nil
	}

	a, ok := (*sym).(*AnnotationSymbol)
	if !ok {
		return nil
	}

	return a
}

func NewContext(parent *Context) *Context {
	return &Context{
		Parent:  parent,
		Symbols: make(map[string]Symbol),
	}
}

type TypeChecker struct {
	Context  *Context
	Warnings []TypeWarning
}

type TypeWarning struct {
	Msg string
	Loc *Location
}

func (c *TypeChecker) AddWarning(msg string, loc *Location) {
	if c == nil {
		return
	}

	c.Warnings = append(c.Warnings, TypeWarning{Msg: msg, Loc: loc})
}

func (c *TypeChecker) GetWarnings() []TypeWarning {
	if c == nil {
		return nil
	}

	return c.Warnings
}

func (c *TypeChecker) FindAllSymbol(prog *ProgramNode) exception.IException {
	for _, node := range prog.Body {
		switch v := node.(type) {
		case *ModelDeclNode:
			if v.Name == nil {
				return exception.NewTypeException("Model name is missing", v.Loc)
			}

			sym := NewTypeSymbol(v.Name.Value, false)
			sym.Generics = v.Generics
			if c.Context.Find(sym) {
				return exception.NewTypeException(fmt.Sprintf("Name '%s' is already defined", sym.Name), v.Name.Loc)
			}

			c.Context.Add(sym)
		case *RestDeclNode:
			if v.Name == nil {
				return exception.NewTypeException("Rest name is missing", v.Loc)
			}

			sym := NewRestSymbol(v.Name.Value)
			if c.Context.Find(sym) {
				return exception.NewTypeException(fmt.Sprintf("Name '%s' is already defined", sym.Name), v.Name.Loc)
			}

			c.Context.Add(sym)
		case *ErrorDeclNode:
			if v.Name == nil {
				return exception.NewTypeException("Error name is missing", v.Loc)
			}

			sym := NewErrorSymbol(v.Name.Value)
			if c.Context.Find(sym) {
				return exception.NewTypeException(fmt.Sprintf("Name '%s' is already defined", sym.Name), v.Name.Loc)
			}

			c.Context.Add(sym)
		}
	}

	return nil
}

func (c *TypeChecker) CheckType(node *TypeDeclNode) exception.IException {
	sym := c.Context.GetTypeByName(node.Name.Value)

	if sym == nil {
		return exception.NewTypeException(fmt.Sprintf("Type '%s' is not defined", node.Name.Value), node.Name.Loc)
	}

	if len(sym.Generics) != len(node.Generics) {
		return exception.NewTypeException(
			fmt.Sprintf("Type '%s' expects %d generic argument(s), got %d", node.Name.Value, len(sym.Generics), len(node.Generics)),
			node.Loc,
		)
	}

	return nil
}

func (c *TypeChecker) CheckAnnotations(nodes []*AnnotationNode) exception.IException {
	for _, node := range nodes {
		err := c.CheckAnnotation(node)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *TypeChecker) CheckAnnotation(node *AnnotationNode) exception.IException {
	if node == nil || node.Name == nil {
		var loc *Location
		if node != nil {
			loc = node.Loc
		}
		return exception.NewTypeException("Annotation name is missing", loc)
	}

	sym := c.Context.GetAnnotationByName(node.Name.Value)
	if sym == nil {
		return exception.NewTypeException(fmt.Sprintf("Annotation '%s' is not defined", node.Name.Value), node.Name.Loc)
	}

	if len(node.Args) != len(sym.ArgOrder) {
		return exception.NewTypeException(
			fmt.Sprintf("Annotation '%s' expects %d argument(s), got %d", sym.Name, len(sym.ArgOrder), len(node.Args)),
			node.Loc,
		)
	}

	for i, arg := range node.Args {
		argName := sym.ArgOrder[i]
		expected := sym.Args[argName]
		if expected == nil || expected.Name == nil {
			continue
		}

		if !isValueAssignableToType(arg, expected.Name.Value) {
			return exception.NewTypeException(
				fmt.Sprintf("Annotation '%s' argument '%s' expects type '%s'", sym.Name, argName, expected.Name.Value),
				arg.GetLocation(),
			)
		}
	}

	return nil
}

func (c *TypeChecker) CheckModelFieldType(node *ModelFieldDeclNode) exception.IException {
	err := c.CheckType(node.Type)
	if err != nil {
		return err
	}

	err = c.CheckAnnotations(node.Annotations)
	if err != nil {
		return err
	}

	return nil
}

func (c *TypeChecker) CheckModelType(node *ModelDeclNode) exception.IException {
	if node.Name == nil {
		return exception.NewTypeException("Model name is missing", node.Loc)
	}

	ctx_ := c.Context
	c.Context = NewContext(c.Context)

	for _, typeVar := range node.Generics {
		sym := NewTypeSymbol(typeVar.Name.Value, false)
		if c.Context.Find(sym) {
			return exception.NewTypeException(fmt.Sprintf("Type parameter '%s' is already defined", sym.Name), typeVar.Name.Loc)
		}

		c.Context.Add(sym)
	}

	for _, field := range node.Fields {
		err := c.CheckModelFieldType(field)
		if err != nil {
			return err
		}
	}

	err := c.CheckAnnotations(node.Annotations)
	if err != nil {
		return err
	}

	c.Context = ctx_

	return nil
}

func (c *TypeChecker) CheckRestType(node *RestDeclNode) exception.IException {
	if node == nil {
		fallbackLoc := NewLocation("<unknown>", NewPosition(1, 1), NewPosition(1, 1))
		return exception.NewTypeException("Rest name is missing", fallbackLoc)
	}

	if node.Name == nil {
		return exception.NewTypeException("Rest name is missing", node.Loc)
	}

	if node.MethodValue == nil {
		return exception.NewTypeException("Rest property 'method' is required", node.Loc)
	}
	if node.PathValue == nil {
		return exception.NewTypeException("Rest property 'path' is required", node.Loc)
	}

	methodValue, ok := node.MethodValue.(*StringValueNode)
	if !ok {
		return exception.NewTypeException("Rest property 'method' must be a string literal", node.MethodValue.GetLocation())
	}

	if !isHTTPMethod(methodValue.Value) {
		return exception.NewTypeException("Rest property 'method' must be one of GET, POST, PUT, PATCH, DELETE", methodValue.Loc)
	}

	httpMethod := strings.ToUpper(methodValue.Value)

	if _, ok := node.PathValue.(*StringValueNode); !ok {
		return exception.NewTypeException("Rest property 'path' must be a string literal", node.PathValue.GetLocation())
	}

	if node.RequestBodyType != nil {
		if err := c.CheckType(node.RequestBodyType); err != nil {
			return err
		}

		if !c.isUserDefinedType(node.RequestBodyType) && !isNullType(node.RequestBodyType) {
			return exception.NewTypeException("Rest property 'requestBody' must be a user-defined type", node.RequestBodyType.Loc)
		}

		if httpMethod == "GET" && !isNullType(node.RequestBodyType) {
			c.AddWarning("GET endpoint usually does not need requestBody", node.RequestBodyType.Loc)
		}
	}

	if node.ResponseBodyType != nil {
		if err := c.CheckType(node.ResponseBodyType); err != nil {
			return err
		}

		if !c.isUserDefinedType(node.ResponseBodyType) && !isNullType(node.ResponseBodyType) {
			return exception.NewTypeException("Rest property 'responseBody' must be a user-defined type", node.ResponseBodyType.Loc)
		}
	}

	if node.QueriesValue != nil {
		queriesValue, ok := node.QueriesValue.(*ArrayValueNode)
		if !ok {
			return exception.NewTypeException("Rest property 'queries' must be an array literal", node.QueriesValue.GetLocation())
		}

		for _, item := range queriesValue.Values {
			if _, ok := item.(*StringValueNode); !ok {
				return exception.NewTypeException("Rest property 'queries' must be string[]", item.GetLocation())
			}
		}
	}

	return nil
}

func (c *TypeChecker) CheckErrorType(node *ErrorDeclNode) exception.IException {
	if node == nil {
		fallbackLoc := NewLocation("<unknown>", NewPosition(1, 1), NewPosition(1, 1))
		return exception.NewTypeException("Error name is missing", fallbackLoc)
	}

	if node.Name == nil {
		return exception.NewTypeException("Error name is missing", node.Loc)
	}

	if node.MessageValue == nil {
		return exception.NewTypeException("Error property 'message' is required", node.Loc)
	}

	if _, ok := node.MessageValue.(*StringValueNode); !ok {
		return exception.NewTypeException("Error property 'message' must be a string literal", node.MessageValue.GetLocation())
	}

	if node.CodeValue != nil {
		if _, ok := node.CodeValue.(*StringValueNode); !ok {
			return exception.NewTypeException("Error property 'code' must be a string literal", node.CodeValue.GetLocation())
		}
	}

	if node.ScopeValue != nil {
		if _, ok := node.ScopeValue.(*StringValueNode); !ok {
			return exception.NewTypeException("Error property 'scope' must be a string literal", node.ScopeValue.GetLocation())
		}
	}

	return nil
}

func (c *TypeChecker) Check(ast *ProgramNode) exception.IException {
	c.Warnings = c.Warnings[:0]

	if err := c.FindAllSymbol(ast); err != nil {
		return err
	}
	for _, node := range ast.Body {
		switch v := node.(type) {
		case *ModelDeclNode:
			err := c.CheckModelType(v)
			if err != nil {
				return err
			}
		case *RestDeclNode:
			err := c.CheckRestType(v)
			if err != nil {
				return err
			}
		case *ErrorDeclNode:
			err := c.CheckErrorType(v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewTypeChecker() *TypeChecker {
	ctx := NewContext(nil)
	ctx.Add(NewTypeSymbol("String", true))
	ctx.Add(NewTypeSymbol("Int", true))
	ctx.Add(NewTypeSymbol("Float", true))
	ctx.Add(NewTypeSymbol("Bool", true))
	ctx.Add(NewTypeSymbol("Null", true))
	ctx.Add(NewTypeSymbol("Any", true))

	arrayType := NewTypeSymbol("Array", true)
	arrayType.Generics = append(arrayType.Generics, &TypeVarNode{Name: &IdentNode{Value: "T"}})
	ctx.Add(arrayType)

	createConstructor := NewAnnotationSymbol("CreateConstructor", true)
	ctx.Add(createConstructor)

	mapper := NewAnnotationSymbol("Mapper", true)
	ctx.Add(mapper)

	is := NewAnnotationSymbol("Is", true)
	is.Generics = append(is.Generics, &TypeVarNode{Name: &IdentNode{Value: "T"}})
	is.ArgOrder = append(is.ArgOrder, "constraint", "message")
	is.Args["constraint"] = newTypeRef("Any")
	is.Args["message"] = newTypeRef("String")
	ctx.Add(is)

	min := NewAnnotationSymbol("Min", true)
	min.ArgOrder = append(min.ArgOrder, "value", "message")
	min.Args["value"] = newTypeRef("Float")
	min.Args["message"] = newTypeRef("String")
	ctx.Add(min)

	max := NewAnnotationSymbol("Max", true)
	max.ArgOrder = append(max.ArgOrder, "value", "message")
	max.Args["value"] = newTypeRef("Float")
	max.Args["message"] = newTypeRef("String")
	ctx.Add(max)

	length := NewAnnotationSymbol("Length", true)
	length.ArgOrder = append(length.ArgOrder, "value", "message")
	length.Args["value"] = newTypeRef("Int")
	length.Args["message"] = newTypeRef("String")
	ctx.Add(length)

	minLength := NewAnnotationSymbol("MinLength", true)
	minLength.ArgOrder = append(minLength.ArgOrder, "value", "message")
	minLength.Args["value"] = newTypeRef("Int")
	minLength.Args["message"] = newTypeRef("String")
	ctx.Add(minLength)

	maxLength := NewAnnotationSymbol("MaxLength", true)
	maxLength.ArgOrder = append(maxLength.ArgOrder, "value", "message")
	maxLength.Args["value"] = newTypeRef("Int")
	maxLength.Args["message"] = newTypeRef("String")
	ctx.Add(maxLength)

	rangeValue := NewAnnotationSymbol("Range", true)
	rangeValue.ArgOrder = append(rangeValue.ArgOrder, "min", "max", "message")
	rangeValue.Args["min"] = newTypeRef("Float")
	rangeValue.Args["max"] = newTypeRef("Float")
	rangeValue.Args["message"] = newTypeRef("String")
	ctx.Add(rangeValue)

	matches := NewAnnotationSymbol("Matches", true)
	matches.ArgOrder = append(matches.ArgOrder, "pattern", "message")
	matches.Args["pattern"] = newTypeRef("String")
	matches.Args["message"] = newTypeRef("String")
	ctx.Add(matches)

	contains := NewAnnotationSymbol("Contains", true)
	contains.ArgOrder = append(contains.ArgOrder, "value", "message")
	contains.Args["value"] = newTypeRef("String")
	contains.Args["message"] = newTypeRef("String")
	ctx.Add(contains)

	startsWith := NewAnnotationSymbol("StartsWith", true)
	startsWith.ArgOrder = append(startsWith.ArgOrder, "value", "message")
	startsWith.Args["value"] = newTypeRef("String")
	startsWith.Args["message"] = newTypeRef("String")
	ctx.Add(startsWith)

	endsWith := NewAnnotationSymbol("EndsWith", true)
	endsWith.ArgOrder = append(endsWith.ArgOrder, "value", "message")
	endsWith.Args["value"] = newTypeRef("String")
	endsWith.Args["message"] = newTypeRef("String")
	ctx.Add(endsWith)

	in := NewAnnotationSymbol("In", true)
	in.ArgOrder = append(in.ArgOrder, "values", "message")
	in.Args["values"] = newTypeRef("Array")
	in.Args["message"] = newTypeRef("String")
	ctx.Add(in)

	isEmail := NewAnnotationSymbol("IsEmail", true)
	isEmail.ArgOrder = append(isEmail.ArgOrder, "message")
	isEmail.Args["message"] = newTypeRef("String")
	ctx.Add(isEmail)

	isNumber := NewAnnotationSymbol("IsNumber", true)
	isNumber.ArgOrder = append(isNumber.ArgOrder, "message")
	isNumber.Args["message"] = newTypeRef("String")
	ctx.Add(isNumber)

	isURL := NewAnnotationSymbol("IsURL", true)
	isURL.ArgOrder = append(isURL.ArgOrder, "message")
	isURL.Args["message"] = newTypeRef("String")
	ctx.Add(isURL)

	isUUID := NewAnnotationSymbol("IsUUID", true)
	isUUID.ArgOrder = append(isUUID.ArgOrder, "message")
	isUUID.Args["message"] = newTypeRef("String")
	ctx.Add(isUUID)

	isDate := NewAnnotationSymbol("IsDate", true)
	isDate.ArgOrder = append(isDate.ArgOrder, "message")
	isDate.Args["message"] = newTypeRef("String")
	ctx.Add(isDate)

	isDateTime := NewAnnotationSymbol("IsDateTime", true)
	isDateTime.ArgOrder = append(isDateTime.ArgOrder, "message")
	isDateTime.Args["message"] = newTypeRef("String")
	ctx.Add(isDateTime)

	isAlpha := NewAnnotationSymbol("IsAlpha", true)
	isAlpha.ArgOrder = append(isAlpha.ArgOrder, "message")
	isAlpha.Args["message"] = newTypeRef("String")
	ctx.Add(isAlpha)

	isAlnum := NewAnnotationSymbol("IsAlnum", true)
	isAlnum.ArgOrder = append(isAlnum.ArgOrder, "message")
	isAlnum.Args["message"] = newTypeRef("String")
	ctx.Add(isAlnum)

	notNull := NewAnnotationSymbol("NotNull", true)
	notNull.ArgOrder = append(notNull.ArgOrder, "message")
	notNull.Args["message"] = newTypeRef("String")
	ctx.Add(notNull)

	isBool := NewAnnotationSymbol("IsBool", true)
	isBool.ArgOrder = append(isBool.ArgOrder, "message")
	isBool.Args["message"] = newTypeRef("String")
	ctx.Add(isBool)

	isModel := NewAnnotationSymbol("IsModel", true)
	isModel.ArgOrder = append(isModel.ArgOrder, "message")
	isModel.Args["message"] = newTypeRef("String")
	ctx.Add(isModel)

	i := &TypeChecker{
		Context:  ctx,
		Warnings: make([]TypeWarning, 0),
	}

	return i
}

func newTypeRef(name string) *TypeDeclNode {
	return &TypeDeclNode{
		Name:     &IdentNode{Value: name},
		Generics: make([]*TypeDeclNode, 0),
	}
}

func isValueAssignableToType(node ASTValueNode, expectedTypeName string) bool {
	switch expectedTypeName {
	case "Any":
		return true
	case "String":
		_, ok := node.(*StringValueNode)
		return ok
	case "Bool":
		_, ok := node.(*BooleanValueNode)
		return ok
	case "Null":
		_, ok := node.(*NullValueNode)
		return ok
	case "Array":
		_, ok := node.(*ArrayValueNode)
		return ok
	case "Int":
		n, ok := node.(*NumberValueNode)
		if !ok {
			return false
		}
		_, err := strconv.ParseInt(n.Value, 10, 64)
		return err == nil
	case "Float":
		n, ok := node.(*NumberValueNode)
		if !ok {
			return false
		}
		_, err := strconv.ParseFloat(n.Value, 64)
		return err == nil
	default:
		return false
	}
}

func isStringType(node *TypeDeclNode) bool {
	if node == nil || node.Name == nil {
		return false
	}

	if len(node.Generics) > 0 {
		return false
	}

	return node.Name.Value == "String" || node.Name.Value == "string"
}

func isArrayOfStringType(node *TypeDeclNode) bool {
	if node == nil || node.Name == nil {
		return false
	}

	if node.Name.Value != "Array" || len(node.Generics) != 1 {
		return false
	}

	return isStringType(node.Generics[0])
}

func isHTTPMethod(value string) bool {
	switch strings.ToUpper(value) {
	case "GET", "POST", "PUT", "PATCH", "DELETE":
		return true
	default:
		return false
	}
}

func isNullType(node *TypeDeclNode) bool {
	if node == nil || node.Name == nil {
		return false
	}

	if len(node.Generics) != 0 {
		return false
	}

	return node.Name.Value == "Null"
}

func (c *TypeChecker) isUserDefinedType(node *TypeDeclNode) bool {
	if node == nil || node.Name == nil {
		return false
	}

	sym := c.Context.GetTypeByName(node.Name.Value)
	if sym == nil {
		return false
	}

	return !sym.BuiltIn()
}
