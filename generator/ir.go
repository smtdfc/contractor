package generator

import "github.com/smtdfc/contractor/parser"

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

func (s *SourceSpan) ToLocation() *parser.Location {
	if s == nil {
		return nil
	}

	return parser.NewLocation(
		s.File,
		parser.NewPosition(s.StartLine, s.StartCol),
		parser.NewPosition(s.EndLine, s.EndCol),
	)
}

type ValueIR struct {
	Kind  string
	Value any
}

type AnnotationIR struct {
	Name string
	Args []*ValueIR
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
	Annotations         []*AnnotationIR
	Fields              []*ModelField
	IsCreateConstructor bool
	IsCreateMapper      bool
}

func (m *ModelIR) GetKind() string {
	return "model"
}

type ProgramIR struct {
	Errors []*ErrorIR
	Models []*ModelIR
	Rests  []*RestEndpointIR
}

func (p *ProgramIR) GetKind() string {
	return "program"
}

type ErrorIR struct {
	Span    *SourceSpan
	Name    string
	Code    *string
	Message string
	Scope   *string
}

func (e *ErrorIR) GetKind() string {
	return "error"
}

type ModelField struct {
	Span         *SourceSpan
	Name         string
	Annotations  []*AnnotationIR
	Type         *TypeIR
	IsOptional   bool
	DefaultValue *ValueIR
	Validators   []*FieldValidator
}

type FieldValidator struct {
	Name string
	Args []*ValueIR
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

type RestEndpointIR struct {
	Span             *SourceSpan
	Name             string
	Method           string
	Path             string
	RequestBodyType  *TypeIR
	ResponseBodyType *TypeIR
	Queries          []string
}

func (r *RestEndpointIR) GetKind() string {
	return "rest-endpoint"
}
