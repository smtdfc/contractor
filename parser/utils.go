package parser

import (
	"fmt"
	"strings"
)

func PrintTokenList(tokens TokenList) {
	for _, token := range tokens {
		fmt.Println(token)
	}
}

func IsAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func IsAlphaNumeric(r rune) bool {
	return IsAlpha(r) || IsDigit(r)
}

func FormatTypeDecl(node *TypeDeclNode) string {
	if node == nil {
		return "<nil>"
	}

	if node.Name == nil {
		return "<unnamed-type>"
	}

	if len(node.Generics) == 0 {
		return node.Name.Value
	}

	args := make([]string, 0, len(node.Generics))
	for _, generic := range node.Generics {
		args = append(args, FormatTypeDecl(generic))
	}

	return fmt.Sprintf("%s<%s>", node.Name.Value, strings.Join(args, ", "))
}

func FormatTypeVars(vars []*TypeVarNode) string {
	if len(vars) == 0 {
		return ""
	}

	items := make([]string, 0, len(vars))
	for _, item := range vars {
		if item != nil && item.Name != nil {
			items = append(items, item.Name.Value)
		}
	}

	return fmt.Sprintf("<%s>", strings.Join(items, ", "))
}

func FormatValueNode(value ASTValueNode) string {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {
	case *StringValueNode:
		return fmt.Sprintf("\"%s\"", v.Value)
	case *NumberValueNode:
		return v.Value
	case *BooleanValueNode:
		return v.Values
	case *NullValueNode:
		return "null"
	case *ArrayValueNode:
		items := make([]string, 0, len(v.Values))
		for _, item := range v.Values {
			items = append(items, FormatValueNode(item))
		}
		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	default:
		return value.GetType()
	}
}

func FormatAnnotation(annotation *AnnotationNode) string {
	if annotation == nil {
		return "@<nil>"
	}

	name := "<unnamed-annotation>"
	if annotation.Name != nil {
		name = annotation.Name.Value
	}

	if len(annotation.Args) == 0 {
		return fmt.Sprintf("@%s", name)
	}

	args := make([]string, 0, len(annotation.Args))
	for _, arg := range annotation.Args {
		args = append(args, FormatValueNode(arg))
	}

	return fmt.Sprintf("@%s(%s)", name, strings.Join(args, ", "))
}

func PrintAnnotations(annotations []*AnnotationNode, indent int) {
	if len(annotations) == 0 {
		return
	}

	tab := strings.Repeat("  ", indent)
	for _, annotation := range annotations {
		fmt.Printf("%s├── Annotation: %s\n", tab, FormatAnnotation(annotation))
	}
}

func PrintAST(node any, indent int) {
	if node == nil {
		return
	}
	tab := strings.Repeat("  ", indent)

	switch v := node.(type) {
	case *ProgramNode:
		fmt.Printf("%s[Program]\n", tab)
		for _, stmt := range v.Body {
			PrintAST(stmt, indent+1)
		}

	case *ModelDeclNode:
		modelName := "<unnamed-model>"
		if v.Name != nil {
			modelName = v.Name.Value
		}
		fmt.Printf("%s├── Model: %s%s\n", tab, modelName, FormatTypeVars(v.Generics))
		PrintAnnotations(v.Annotations, indent+1)
		for _, field := range v.Fields {
			PrintAST(field, indent+2)
		}

	case *ModelFieldDeclNode:
		opt := ""
		if v.Optional {
			opt = " (Optional)"
		}
		fieldName := "<unnamed-field>"
		if v.Name != nil {
			fieldName = v.Name.Value
		}
		fmt.Printf("%s├── Field: %s%s\n", tab, fieldName, opt)
		PrintAnnotations(v.Annotations, indent+2)
		PrintAST(v.Type, indent+3)

	case *TypeDeclNode:
		fmt.Printf("%s└── Type: %s\n", tab, FormatTypeDecl(v))

	case *TypeVarNode:
		typeVarName := "<unnamed-typevar>"
		if v.Name != nil {
			typeVarName = v.Name.Value
		}
		fmt.Printf("%s└── TypeVar: %s\n", tab, typeVarName)
	}
}
