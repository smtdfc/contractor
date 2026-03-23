package parser

import (
	"fmt"

	"github.com/smtdfc/contractor/exception"
)

type Symbol struct {
	Name      string
	IsBuiltIn bool
	Generics  []*TypeVarNode
}

func NewSymbol(name string, isBuiltIn bool) *Symbol {
	return &Symbol{Name: name, IsBuiltIn: isBuiltIn, Generics: make([]*TypeVarNode, 0)}
}

type Context struct {
	Parent  *Context
	Symbols map[string]*Symbol
}

func (c *Context) Add(sym *Symbol) {
	c.Symbols[sym.Name] = sym
}

func (c *Context) Find(sym *Symbol) bool {
	if c.Symbols[sym.Name] != nil {
		return true
	}

	if c.Parent != nil {
		return c.Parent.Find(sym)
	}

	return false
}

func (c *Context) GetByName(name string) *Symbol {
	if c.Symbols[name] != nil {
		return c.Symbols[name]
	}

	if c.Parent != nil {
		return c.Parent.GetByName(name)
	}

	return nil
}
func NewContext(parent *Context) *Context {
	return &Context{
		Parent:  parent,
		Symbols: make(map[string]*Symbol),
	}
}

type TypeChecker struct {
	Context *Context
}

func (c *TypeChecker) FindAllSymbol(prog *ProgramNode) exception.IException {
	for _, node := range prog.Body {
		switch v := node.(type) {
		case *ModelDeclNode:
			if v.Name == nil {
				return exception.NewTypeException("Model name is missing", v.Loc)
			}

			sym := NewSymbol(v.Name.Value, false)
			sym.Generics = v.Generics
			if c.Context.Find(sym) {
				return exception.NewTypeException(fmt.Sprintf("Name '%s' is already defined", sym.Name), v.Name.Loc)
			}

			c.Context.Add(sym)
		}
	}

	return nil
}

func (c *TypeChecker) CheckType(node *TypeDeclNode) exception.IException {
	sym := c.Context.GetByName(node.Name.Value)

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

func (c *TypeChecker) CheckModelFieldType(node *ModelFieldDeclNode) exception.IException {
	err := c.CheckType(node.Type)
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
		sym := NewSymbol(typeVar.Name.Value, false)
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

	c.Context = ctx_

	return nil
}

func (c *TypeChecker) Check(ast *ProgramNode) exception.IException {
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
		}
	}

	return nil
}

func NewTypeChecker() *TypeChecker {
	ctx := NewContext(nil)
	ctx.Add(NewSymbol("String", true))
	ctx.Add(NewSymbol("Int", true))
	ctx.Add(NewSymbol("Float", true))
	ctx.Add(NewSymbol("Boo;", true))
	ctx.Add(NewSymbol("Null", true))
	ctx.Add(NewSymbol("Any", true))
	i := &TypeChecker{
		Context: ctx,
	}

	return i
}
