package main

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

func (c *TypeChecker) CheckType(ctx *Context, t *TypeDeclarationNode) BaseError {
	if t == nil {
		return nil
	}

	symbol, err := ctx.FindSymbol(t.Name, true, t)
	if err != nil {
		return err
	}

	if !symbol.HasGeneric && t.Generic != nil {
		return NewTypeError("Symbol "+t.Name+" is not a generic type", t.Loc.GetErrorLocation())
	} else if symbol.HasGeneric && t.Generic == nil {
		return NewTypeError("Symbol "+t.Name+" requires a generic type argument", t.Loc.GetErrorLocation())
	} else if symbol.HasGeneric && t.Generic != nil {
		return c.CheckType(ctx, t.Generic.(*TypeDeclarationNode))
	}

	return nil
}

func (c *TypeChecker) CheckModelField(parentCtx *Context, node *ModelFieldNode) BaseError {
	err := c.CheckType(parentCtx, node.Type)
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

	return nil
}

func (c *TypeChecker) Check(ast *AST) BaseError {
	ctx := NewContext()
	ctx.Parent = globalCtx
	for _, node := range ast.Statements {
		switch node.(type) {
		case *ModelStatementNode:
			err := c.CheckModelStatement(ctx, node.(*ModelStatementNode))
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
