package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createNodeLoc() *NodeLocation {
	return &NodeLocation{
		Start: NewTokenLocation(
			NewPosition(1, 1),
			NewPosition(1, 4),
			"test",
		),

		End: NewTokenLocation(
			NewPosition(1, 1),
			NewPosition(1, 4),
			"test",
		),
	}
}

// Helper to create a basic TypeDeclarationNode for testing
func createTypeNode(name string, generic *TypeDeclarationNode) *TypeDeclarationNode {
	return &TypeDeclarationNode{
		Name:    name,
		Generic: generic,
		Loc:     createNodeLoc(),
	}
}

func TestTypeChecker_CheckType(t *testing.T) {
	tc := NewTypeChecker()
	ctx := CreateGlobalContext()

	// Scenario 1: Valid primitive type
	t1 := createTypeNode("Number", nil)
	assert.NoError(t, tc.CheckType(ctx, t1))

	// Scenario 2: Valid generic type (Array<String>)
	t2 := createTypeNode("Array", createTypeNode("String", nil))
	assert.NoError(t, tc.CheckType(ctx, t2))

	// Scenario 3: Error - Generic provided for non-generic type (String<Number>)
	t3 := createTypeNode("String", createTypeNode("Number", nil))
	err := tc.CheckType(ctx, t3)
	assert.Error(t, err)
	assert.Contains(t, err.GetMessage(), "is not a generic type")

	// Scenario 4: Error - Missing generic for generic type (Array)
	t4 := createTypeNode("Array", nil)
	err = tc.CheckType(ctx, t4)
	assert.Error(t, err)
	assert.Contains(t, err.GetMessage(), "requires a generic type argument")
}

func TestTypeChecker_SymbolScope(t *testing.T) {
	parent := NewContext()
	parent.Symbols["Base"] = NewSymbol(DEFINED_TYPE, "Base", nil, false)

	child := parent.CreateChildContext()

	// Scenario: Find symbol in parent context (recursion = true)
	sym, err := child.FindSymbol("Base", true, &ModelStatementNode{})
	assert.NoError(t, err)
	assert.NotNil(t, sym)
	assert.Equal(t, "Base", sym.Name)

	// Scenario: Symbol not found
	_, err = child.FindSymbol("Undefined", true, &ModelStatementNode{
		Loc: createNodeLoc(),
	})
	assert.Error(t, err)
}

func TestTypeChecker_CheckModelStatement(t *testing.T) {
	tc := NewTypeChecker()
	global := CreateGlobalContext()
	ctx := NewContext()
	ctx.Parent = global

	// Scenario 1: Valid Model definition
	// model User { String name }
	model := &ModelStatementNode{
		Name: "User",
		Fields: []*ModelFieldNode{
			{
				Name: "name",
				Type: createTypeNode("String", nil),
			},
		},
		Loc: createNodeLoc(),
	}
	assert.NoError(t, tc.CheckModelStatement(ctx, model))

	// Scenario 2: Error - Re-declaring the same model
	err := tc.CheckModelStatement(ctx, model)
	assert.Error(t, err)
	assert.Contains(t, err.GetMessage(), "Re-declare symbol User")

	// Scenario 3: Valid Generic Model
	// model Box<T> { T content }
	typeVar := &TypeVarNode{Name: "T", Loc: createNodeLoc()}
	genericModel := &ModelStatementNode{
		Name:    "Box",
		TypeVar: typeVar,
		Fields: []*ModelFieldNode{
			{
				Name: "content",
				Type: createTypeNode("T", nil),
			},
		},
		Loc: createNodeLoc(),
	}
	assert.NoError(t, tc.CheckModelStatement(ctx, genericModel))

	// Verify "Box" is now marked as HasGeneric in context
	sym, _ := ctx.FindSymbol("Box", false, genericModel)
	assert.True(t, sym.HasGeneric)
}
