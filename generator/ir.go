package generator

type IR interface {
	GetKind() string
}

type SourceSpan struct {
	File      string
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
}

type AnnotationArgIR struct {
	Kind  string
	Value any
}

type AnnotationIR struct {
	Name string
	Args []AnnotationArgIR
}

type TypeKind string

const (
	TypeKindBuiltin TypeKind = "builtin"
	TypeKindModel   TypeKind = "model"
	TypeKindGeneric TypeKind = "generic"
	TypeKindUnknown TypeKind = "unknown"
)

type ModelIR struct {
	Span                *SourceSpan
	Name                string
	TypeParams          []string
	Annotations         []AnnotationIR
	Fields              []*ModelField
	IsCreateConstructor bool
	IsCreateMapper      bool
}

func (m *ModelIR) GetKind() string {
	return "model"
}

type ModelField struct {
	Span            *SourceSpan
	Name            string
	Annotations     []AnnotationIR
	Type            *TypeIR
	IsOptional      bool
	DefaultValue    *AnnotationArgIR
	ValidationRules ModelValidationRuleList
}

type ModelValidationRuleList []*ModelValidationRule[any]

type ModelValidationRule[T any] struct {
	Name    string
	Value   T
	Message string
}

type TypeIR struct {
	Span        *SourceSpan
	Kind        TypeKind
	Name        string
	Generics    []*TypeIR
	ResolvedRef string
}

func (t *TypeIR) GetKind() string {
	return "type"
}
