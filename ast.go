package main

type Node interface {
	GetType() string
	GetLocation() *NodeLocation
}

type ListNode []Node

type LiteralNode struct {
	Type        TokenType
	Value       string
	Annotations *AnnotationChainNode
	Loc         *NodeLocation
}

func (n *LiteralNode) GetType() string {
	return "Literial"
}

func (n *LiteralNode) GetLocation() *NodeLocation {
	return n.Loc
}

func NewLiteralNode(t TokenType, value string, loc *NodeLocation) *LiteralNode {
	return &LiteralNode{
		Type:  t,
		Value: value,
		Loc:   loc,
	}
}

type ModelStatementNode struct {
	Name        string
	Fields      []*ModelFieldNode
	TypeVar     *TypeVarNode
	Annotations *AnnotationChainNode
	Loc         *NodeLocation
}

func (n *ModelStatementNode) GetType() string {
	return "ModelStatement"
}

func (n *ModelStatementNode) GetLocation() *NodeLocation {
	return n.Loc
}

func NewModelStatementNode(name string, typeVar *TypeVarNode, anno *AnnotationChainNode, fields []*ModelFieldNode, loc *NodeLocation) *ModelStatementNode {
	return &ModelStatementNode{Name: name, TypeVar: typeVar, Annotations: anno, Fields: fields, Loc: loc}
}

type AnnotationNode struct {
	Name string
	Args ListNode
	Loc  *NodeLocation
}

func (n *AnnotationNode) GetType() string {
	return "Annotation"
}

func (n *AnnotationNode) GetLocation() *NodeLocation {
	return n.Loc
}

func NewAnnotationNode(name string, args ListNode, loc *NodeLocation) *AnnotationNode {
	return &AnnotationNode{Name: name, Args: args, Loc: loc}
}

type ListAnnotation []*AnnotationNode

type AnnotationChainNode struct {
	List ListAnnotation
	Loc  *NodeLocation
}

func (n *AnnotationChainNode) GetType() string {
	return "AnnotationChain"
}

func (n *AnnotationChainNode) GetLocation() *NodeLocation {
	return n.Loc
}

func NewAnnotationChainNode(list ListAnnotation, loc *NodeLocation) *AnnotationChainNode {
	return &AnnotationChainNode{List: list, Loc: loc}
}

type TypeDeclarationNode struct {
	Name    string
	Generic Node
	Loc     *NodeLocation
}

func (n *TypeDeclarationNode) GetType() string {
	return "TypeDeclaration"
}

func (n *TypeDeclarationNode) GetLocation() *NodeLocation {
	return n.Loc
}

func NewTypeDeclarationNode(name string, generic Node, loc *NodeLocation) *TypeDeclarationNode {
	return &TypeDeclarationNode{Name: name, Generic: generic, Loc: loc}
}

type ModelFieldNode struct {
	Name        string
	Type        *TypeDeclarationNode
	Annotations *AnnotationChainNode
	Loc         *NodeLocation
}

func (n *ModelFieldNode) GetType() string {
	return "ModelField"
}

func (n *ModelFieldNode) GetLocation() *NodeLocation {
	return n.Loc
}

func NewModelFieldNode(name string, annotations *AnnotationChainNode, t *TypeDeclarationNode, loc *NodeLocation) *ModelFieldNode {
	return &ModelFieldNode{Name: name, Annotations: annotations, Type: t, Loc: loc}
}

type TypeVarNode struct {
	Name string
	Loc  *NodeLocation
}

func (n *TypeVarNode) GetType() string {
	return "TypeVar"
}

func (n *TypeVarNode) GetLocation() *NodeLocation {
	return n.Loc
}

func NewTypeVarNode(name string, loc *NodeLocation) *TypeVarNode {
	return &TypeVarNode{Name: name, Loc: loc}
}

type AST struct {
	Statements ListNode
	Loc        *NodeLocation
}

func (n *AST) GetType() string {
	return "AST"
}

func (n *AST) GetLocation() *NodeLocation {
	return n.Loc
}

func NewAST(stats ListNode, loc *NodeLocation) *AST {
	return &AST{
		Statements: stats,
		Loc:        loc,
	}
}
