package parser

type ASTNode interface {
	GetLocation() *Location
	GetType() string
}

type ProgramNode struct {
	Body []ASTNode
	Loc  *Location
}

func (n *ProgramNode) GetLocation() *Location {
	return n.Loc
}

func (n *ProgramNode) GetType() string {
	return "Program"
}

type IdentNode struct {
	Value string
	Loc   *Location
}

func (n *IdentNode) GetLocation() *Location {
	return n.Loc
}

func (n *IdentNode) GetType() string {
	return "Ident"
}

type ModelDeclNode struct {
	Name        *IdentNode
	Generics    []*TypeVarNode
	Fields      []*ModelFieldDeclNode
	Annotations []*AnnotationNode
	Loc         *Location
}

func (n *ModelDeclNode) GetLocation() *Location {
	return n.Loc
}

func (n *ModelDeclNode) GetType() string {
	return "ModelDecl"
}

type RestDeclNode struct {
	Name             *IdentNode
	MethodValue      ASTValueNode
	PathValue        ASTValueNode
	RequestBodyType  *TypeDeclNode
	ResponseBodyType *TypeDeclNode
	QueriesValue     ASTValueNode
	Loc              *Location
}

func (n *RestDeclNode) GetLocation() *Location {
	return n.Loc
}

func (n *RestDeclNode) GetType() string {
	return "RestDecl"
}

type ModelFieldDeclNode struct {
	Name        *IdentNode
	Type        *TypeDeclNode
	Optional    bool
	Annotations []*AnnotationNode
	Loc         *Location
}

func (n *ModelFieldDeclNode) GetLocation() *Location {
	return n.Loc
}

func (n *ModelFieldDeclNode) GetType() string {
	return "ModelFieldDecl"
}

type TypeVarNode struct {
	Name *IdentNode
	Loc  *Location
}

func (n *TypeVarNode) GetLocation() *Location {
	return n.Loc
}

func (n *TypeVarNode) GetType() string {
	return "TypeVar"
}

type TypeDeclNode struct {
	Name     *IdentNode
	Generics []*TypeDeclNode
	Loc      *Location
}

func (n *TypeDeclNode) GetLocation() *Location {
	return n.Loc
}

func (n *TypeDeclNode) GetType() string {
	return "TypeDecl"
}

type AnnotationNode struct {
	Name *IdentNode
	Args []ASTValueNode
	Loc  *Location
}

func (n *AnnotationNode) GetLocation() *Location {
	return n.Loc
}

func (n *AnnotationNode) GetType() string {
	return "Annotation"
}

type ASTValueNode interface {
	ASTNode
	GetKind() string
}

type StringValueNode struct {
	Value string
	Loc   *Location
}

func (n *StringValueNode) GetLocation() *Location {
	return n.Loc
}

func (n *StringValueNode) GetType() string {
	return "StringValue"
}

func (n *StringValueNode) GetKind() string {
	return "String"
}

type NumberValueNode struct {
	Value string
	Loc   *Location
}

func (n *NumberValueNode) GetLocation() *Location {
	return n.Loc
}

func (n *NumberValueNode) GetType() string {
	return "NumberValue"
}

func (n *NumberValueNode) GetKind() string {
	return "Number"
}

type ArrayValueNode struct {
	Values []ASTValueNode
	Loc    *Location
}

func (n *ArrayValueNode) GetLocation() *Location {
	return n.Loc
}

func (n *ArrayValueNode) GetType() string {
	return "ArrayValue"
}

func (n *ArrayValueNode) GetKind() string {
	return "Array"
}

type BooleanValueNode struct {
	Value string
	Loc   *Location
}

func (n *BooleanValueNode) GetLocation() *Location {
	return n.Loc
}

func (n *BooleanValueNode) GetType() string {
	return "BooleanValue"
}

func (n *BooleanValueNode) GetKind() string {
	return "Boolean"
}

type NullValueNode struct {
	Loc *Location
}

func (n *NullValueNode) GetLocation() *Location {
	return n.Loc
}

func (n *NullValueNode) GetType() string {
	return "NullValue"
}

func (n *NullValueNode) GetKind() string {
	return "Null"
}
