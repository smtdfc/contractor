package parser

import (
	"github.com/smtdfc/contractor/exception"
)

type Parser struct {
	File    string
	Tokens  TokenList
	Current *Token
	Index   int
}

func (p *Parser) ParseTypeDecl() (*TypeDeclNode, exception.IException) {
	if p.Current == nil {
		return nil, exception.NewSyntaxException("Excepted type name", p.Tokens[len(p.Tokens)-1].Loc)
	}

	var start, end *Location
	start = p.Current.Loc
	typeNode := TypeDeclNode{
		Generics: make([]*TypeDeclNode, 0),
	}

	if p.Current.MatchType(TT_LSQUARE) {
		p.Next()
		if p.Current == nil || !p.Current.MatchType(TT_RSQUARE) {
			if p.Current == nil {
				return nil, exception.NewSyntaxException("Expected ']' in array type", p.Tokens[len(p.Tokens)-1].Loc)
			}
			return nil, exception.NewSyntaxException("Expected ']' in array type", p.Current.Loc)
		}

		p.Next()
		elemType, err := p.ParseTypeDecl()
		if err != nil {
			return nil, err
		}

		typeNode.Name = &IdentNode{Value: "Array", Loc: start.Copy()}
		typeNode.Generics = append(typeNode.Generics, elemType)
		typeNode.Loc = NewLocation(start.File, start.Start, elemType.Loc.End)
		return &typeNode, nil
	}

	if p.Current.MatchType(TT_IDENT) {
		typeNode.Name = &IdentNode{Value: p.Current.Value, Loc: p.Current.Loc.Copy()}
		end = p.Current.Loc
		p.Next()
	} else {
		return nil, exception.NewSyntaxException(
			"Excepted type name",
			p.Current.Loc,
		)
	}

	if p.Current.Match(TT_OP, "<") {
		p.Next()

		if p.Current.Match(TT_OP, ">") {
			return nil, exception.NewSyntaxException(
				"Expected type argument",
				p.Current.Loc,
			)
		}

		for !p.Current.MatchType(TT_EOF) {
			genericType, err := p.ParseTypeDecl()
			if err != nil {
				return nil, err
			}

			typeNode.Generics = append(typeNode.Generics, genericType)

			if p.Current.MatchType(TT_COMMA) {
				p.Next()
				if p.Current.Match(TT_OP, ">") {
					return nil, exception.NewSyntaxException(
						"Trailing comma in generic type is not allowed",
						p.Current.Loc,
					)
				}
				continue
			}

			if p.Current.Match(TT_OP, ">") {
				end = p.Current.Loc
				p.Next()
				break
			}

			return nil, exception.NewSyntaxException(
				"Expected ',' or '>' in generic type",
				p.Current.Loc,
			)
		}

		if end == nil {
			return nil, exception.NewSyntaxException(
				"Unterminated generic type, expected '>'",
				p.Current.Loc,
			)
		}
	}

	typeNode.Loc = NewLocation(start.File, start.Start, end.End)
	return &typeNode, nil
}

func (p *Parser) ParseModelFieldDecl() (*ModelFieldDeclNode, exception.IException) {
	var start, end *Location
	start = p.Current.Loc
	field := ModelFieldDeclNode{}

	if p.Current.MatchType(TT_IDENT) {
		field.Name = &IdentNode{Value: p.Current.Value, Loc: p.Current.Loc.Copy()}
		p.Next()
	} else {
		return nil, exception.NewSyntaxException(
			"Excepted field name",
			p.Current.Loc,
		)
	}

	if p.Current.MatchType(TT_QUES) {
		field.Optional = true
		p.Next()
	}

	if p.Current.MatchType(TT_COLON) {
		p.Next()
	} else {
		return nil, exception.NewSyntaxException(
			"Excepted field type",
			p.Current.Loc,
		)
	}

	typeNode, err := p.ParseTypeDecl()
	if err != nil {
		return nil, err
	}

	field.Type = typeNode
	end = typeNode.Loc
	field.Loc = NewLocation(start.File, start.Start, end.End)
	return &field, nil
}

func (p *Parser) ParseModelDecl() (*ModelDeclNode, exception.IException) {
	var start *Location
	start = p.Current.Loc
	model := ModelDeclNode{
		Fields:      make([]*ModelFieldDeclNode, 0),
		Generics:    make([]*TypeVarNode, 0),
		Annotations: make([]*AnnotationNode, 0),
	}

	p.Next()

	if p.Current.MatchType(TT_IDENT) {
		model.Name = &IdentNode{Value: p.Current.Value, Loc: p.Current.Loc.Copy()}
		p.Next()
	} else {
		return nil, exception.NewSyntaxException("Expected identifier for model name", p.Current.Loc)
	}

	if p.Current.Match(TT_OP, "<") {
		generics := []*TypeVarNode{}
		p.Next()
		for !p.Current.Match(TT_OP, ">") && !p.Current.MatchType(TT_EOF) {
			if p.Current.MatchType(TT_IDENT) {
				generics = append(generics, &TypeVarNode{
					Name: &IdentNode{Value: p.Current.Value, Loc: p.Current.Loc.Copy()},
					Loc:  p.Current.Loc.Copy(),
				})
				p.Next()
			} else {
				return nil, exception.NewSyntaxException("Expected type parameter name", p.Current.Loc)
			}

			if p.Current.MatchType(TT_COMMA) {
				p.Next()
				if p.Current.Match(TT_OP, ">") {
					return nil, exception.NewSyntaxException("Trailing comma in generics not allowed", p.Current.Loc)
				}
			} else if !p.Current.Match(TT_OP, ">") {
				return nil, exception.NewSyntaxException("Expected ',' or '>' in generic list", p.Current.Loc)
			}
		}

		if !p.Current.Match(TT_OP, ">") {
			return nil, exception.NewSyntaxException("Unterminated generic list, expected '>'", p.Current.Loc)
		}
		p.Next()
		model.Generics = generics
	}

	p.SkipNewLine()

	if !p.Current.MatchType(TT_LBRACE) {
		return nil, exception.NewSyntaxException("Expected {", p.Current.Loc)
	}
	p.Next()
	p.SkipNewLine()

	for p.Current != nil && !p.Current.MatchType(TT_RBRACE) && !p.Current.MatchType(TT_EOF) {
		if p.Current.MatchType(TT_NEWLINE) {
			p.SkipNewLine()
			continue
		}

		fieldAnnotations := make([]*AnnotationNode, 0)
		for p.Current != nil && p.Current.MatchType(TT_DECORATOR) {
			annotation, err := p.ParseAnnotation()
			if err != nil {
				return nil, err
			}
			fieldAnnotations = append(fieldAnnotations, annotation)
			p.SkipNewLine()
		}

		if p.Current == nil || p.Current.MatchType(TT_RBRACE) || p.Current.MatchType(TT_EOF) {
			if len(fieldAnnotations) > 0 {
				return nil, exception.NewSyntaxException("Annotation must be followed by a field", fieldAnnotations[len(fieldAnnotations)-1].Loc)
			}
			break
		}

		if p.Current.MatchType(TT_IDENT) {
			field, err := p.ParseModelFieldDecl()
			if err != nil {
				return nil, err
			}
			field.Annotations = append(field.Annotations, fieldAnnotations...)
			model.Fields = append(model.Fields, field)
		} else {
			return nil, exception.NewSyntaxException("Unexpected token inside model body", p.Current.Loc)
		}

		p.SkipNewLine()
	}

	if p.Current != nil && p.Current.MatchType(TT_RBRACE) {
		model.Loc = NewLocation(start.File, start.Start, p.Current.Loc.End)
		p.Next()
	} else {
		return nil, exception.NewSyntaxException("Expected } at the end of model", p.Current.Loc)
	}

	return &model, nil
}

func (p *Parser) ParseEnumDecl() (*EnumDeclNode, exception.IException) {
	if p.Current == nil || !p.Current.Match(TT_IDENT, "enum") {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected 'enum'", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected 'enum'", p.Current.Loc)
	}

	start := p.Current.Loc
	node := &EnumDeclNode{Members: make([]*IdentNode, 0)}
	p.Next()

	if p.Current == nil || !p.Current.MatchType(TT_IDENT) {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected identifier for enum name", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected identifier for enum name", p.Current.Loc)
	}

	node.Name = &IdentNode{Value: p.Current.Value, Loc: p.Current.Loc.Copy()}
	p.Next()
	p.SkipNewLine()

	if p.Current == nil || !p.Current.MatchType(TT_LBRACE) {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected {", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected {", p.Current.Loc)
	}

	p.Next()
	p.SkipNewLine()

	for p.Current != nil && !p.Current.MatchType(TT_RBRACE) && !p.Current.MatchType(TT_EOF) {
		if p.Current.MatchType(TT_NEWLINE) {
			p.SkipNewLine()
			continue
		}

		if !p.Current.MatchType(TT_IDENT) {
			return nil, exception.NewSyntaxException("Expected enum member name", p.Current.Loc)
		}

		node.Members = append(node.Members, &IdentNode{Value: p.Current.Value, Loc: p.Current.Loc.Copy()})
		p.Next()

		if p.Current != nil && p.Current.MatchType(TT_COMMA) {
			p.Next()
		}

		p.SkipNewLine()
	}

	if p.Current == nil || !p.Current.MatchType(TT_RBRACE) {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected } at the end of enum declaration", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected } at the end of enum declaration", p.Current.Loc)
	}

	node.Loc = NewLocation(start.File, start.Start, p.Current.Loc.End)
	p.Next()

	return node, nil
}

func (p *Parser) ParseRestDecl() (*RestDeclNode, exception.IException) {
	if p.Current == nil || !p.Current.Match(TT_IDENT, "rest") {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected 'rest'", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected 'rest'", p.Current.Loc)
	}

	start := p.Current.Loc
	node := &RestDeclNode{}
	p.Next()

	if p.Current == nil || !p.Current.MatchType(TT_IDENT) {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected identifier for rest endpoint name", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected identifier for rest endpoint name", p.Current.Loc)
	}

	node.Name = &IdentNode{Value: p.Current.Value, Loc: p.Current.Loc.Copy()}
	p.Next()
	p.SkipNewLine()

	if p.Current == nil || !p.Current.MatchType(TT_LBRACE) {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected {", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected {", p.Current.Loc)
	}

	p.Next()
	p.SkipNewLine()

	seen := make(map[string]bool)

	for p.Current != nil && !p.Current.MatchType(TT_RBRACE) && !p.Current.MatchType(TT_EOF) {
		if p.Current.MatchType(TT_NEWLINE) {
			p.SkipNewLine()
			continue
		}

		if !p.Current.MatchType(TT_IDENT) {
			return nil, exception.NewSyntaxException("Expected property name in rest body", p.Current.Loc)
		}

		key := p.Current.Value
		if seen[key] {
			return nil, exception.NewSyntaxException("Duplicate rest property '"+key+"'", p.Current.Loc)
		}
		seen[key] = true
		p.Next()

		if p.Current == nil || !p.Current.MatchType(TT_COLON) {
			if p.Current == nil {
				return nil, exception.NewSyntaxException("Expected ':' in rest property", p.Tokens[len(p.Tokens)-1].Loc)
			}
			return nil, exception.NewSyntaxException("Expected ':' in rest property", p.Current.Loc)
		}

		p.Next()

		switch key {
		case "method":
			v, err := p.ParseValue()
			if err != nil {
				return nil, err
			}
			node.MethodValue = v
		case "path":
			v, err := p.ParseValue()
			if err != nil {
				return nil, err
			}
			node.PathValue = v
		case "requestBody":
			t, err := p.ParseTypeDecl()
			if err != nil {
				return nil, err
			}
			node.RequestBodyType = t
		case "responseBody":
			t, err := p.ParseTypeDecl()
			if err != nil {
				return nil, err
			}
			node.ResponseBodyType = t
		case "queries":
			v, err := p.ParseValue()
			if err != nil {
				return nil, err
			}
			node.QueriesValue = v
		default:
			return nil, exception.NewSyntaxException("Unknown rest property '"+key+"'", p.Current.Loc)
		}

		if p.Current != nil && p.Current.MatchType(TT_COMMA) {
			p.Next()
		}

		p.SkipNewLine()
	}

	if p.Current == nil || !p.Current.MatchType(TT_RBRACE) {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected } at the end of rest declaration", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected } at the end of rest declaration", p.Current.Loc)
	}

	node.Loc = NewLocation(start.File, start.Start, p.Current.Loc.End)
	p.Next()

	return node, nil
}
func (p *Parser) ParseEventDecl() (*EventDeclNode, exception.IException) {
	if p.Current == nil || !p.Current.Match(TT_IDENT, "event") {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected 'event'", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected 'event'", p.Current.Loc)
	}

	start := p.Current.Loc
	node := &EventDeclNode{}
	p.Next()

	if p.Current == nil || !p.Current.MatchType(TT_IDENT) {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected identifier for event name", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected identifier for event name", p.Current.Loc)
	}

	node.Name = &IdentNode{Value: p.Current.Value, Loc: p.Current.Loc.Copy()}
	p.Next()
	p.SkipNewLine()

	if p.Current == nil || !p.Current.MatchType(TT_LBRACE) {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected {", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected {", p.Current.Loc)
	}

	p.Next()
	p.SkipNewLine()

	for p.Current != nil && !p.Current.MatchType(TT_RBRACE) && !p.Current.MatchType(TT_EOF) {
		if p.Current.MatchType(TT_NEWLINE) {
			p.SkipNewLine()
			continue
		}

		if p.Current.MatchType(TT_COMMA) {
			p.Next()
			p.SkipNewLine()
			continue
		}

		if !p.Current.MatchType(TT_IDENT) {
			return nil, exception.NewSyntaxException("Expected property name in event body", p.Current.Loc)
		}

		key := p.Current.Value
		p.Next()

		if p.Current == nil || !p.Current.MatchType(TT_COLON) {
			if p.Current == nil {
				return nil, exception.NewSyntaxException("Expected ':' in event property", p.Tokens[len(p.Tokens)-1].Loc)
			}
			return nil, exception.NewSyntaxException("Expected ':' in event property", p.Current.Loc)
		}

		p.Next()

		switch key {
		case "name":
			v, err := p.ParseValue()
			if err != nil {
				return nil, err
			}
			node.NameValue = v
		case "payload":
			t, err := p.ParseTypeDecl()
			if err != nil {
				return nil, err
			}
			node.PayloadType = t
		default:
			return nil, exception.NewSyntaxException("Unknown event property '"+key+"'", p.Current.Loc)
		}

		if p.Current != nil && p.Current.MatchType(TT_COMMA) {
			p.Next()
		}

		p.SkipNewLine()
	}

	if p.Current == nil || !p.Current.MatchType(TT_RBRACE) {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected } at the end of event declaration", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected } at the end of event declaration", p.Current.Loc)
	}

	node.Loc = NewLocation(start.File, start.Start, p.Current.Loc.End)
	p.Next()

	return node, nil
}

func (p *Parser) ParseErrorDecl() (*ErrorDeclNode, exception.IException) {
	if p.Current == nil || !p.Current.Match(TT_IDENT, "error") {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected 'error'", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected 'error'", p.Current.Loc)
	}

	start := p.Current.Loc
	node := &ErrorDeclNode{}
	p.Next()

	if p.Current == nil || !p.Current.MatchType(TT_IDENT) {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected identifier for error name", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected identifier for error name", p.Current.Loc)
	}

	node.Name = &IdentNode{Value: p.Current.Value, Loc: p.Current.Loc.Copy()}
	p.Next()
	p.SkipNewLine()

	if p.Current == nil || !p.Current.MatchType(TT_LBRACE) {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected {", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected {", p.Current.Loc)
	}

	p.Next()
	p.SkipNewLine()

	seen := make(map[string]bool)

	for p.Current != nil && !p.Current.MatchType(TT_RBRACE) && !p.Current.MatchType(TT_EOF) {
		if p.Current.MatchType(TT_NEWLINE) {
			p.SkipNewLine()
			continue
		}

		if !p.Current.MatchType(TT_IDENT) {
			return nil, exception.NewSyntaxException("Expected property name in error body", p.Current.Loc)
		}

		key := p.Current.Value
		if seen[key] {
			return nil, exception.NewSyntaxException("Duplicate error property '"+key+"'", p.Current.Loc)
		}
		seen[key] = true
		p.Next()

		if p.Current == nil || !p.Current.MatchType(TT_COLON) {
			if p.Current == nil {
				return nil, exception.NewSyntaxException("Expected ':' in error property", p.Tokens[len(p.Tokens)-1].Loc)
			}
			return nil, exception.NewSyntaxException("Expected ':' in error property", p.Current.Loc)
		}

		p.Next()

		switch key {
		case "code":
			v, err := p.ParseValue()
			if err != nil {
				return nil, err
			}
			node.CodeValue = v
		case "message":
			v, err := p.ParseValue()
			if err != nil {
				return nil, err
			}
			node.MessageValue = v
		case "scope":
			v, err := p.ParseValue()
			if err != nil {
				return nil, err
			}
			node.ScopeValue = v
		case "status":
			v, err := p.ParseValue()
			if err != nil {
				return nil, err
			}
			node.StatusValue = v
		default:
			return nil, exception.NewSyntaxException("Unknown error property '"+key+"'", p.Current.Loc)
		}

		if p.Current != nil && p.Current.MatchType(TT_COMMA) {
			p.Next()
		}

		p.SkipNewLine()
	}

	if p.Current == nil || !p.Current.MatchType(TT_RBRACE) {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected } at the end of error declaration", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected } at the end of error declaration", p.Current.Loc)
	}

	node.Loc = NewLocation(start.File, start.Start, p.Current.Loc.End)
	p.Next()

	return node, nil
}

func (p *Parser) SkipNewLine() {
	for p.Current != nil && p.Current.MatchType(TT_NEWLINE) {
		p.Next()
	}
}

func (p *Parser) ParseValue() (ASTValueNode, exception.IException) {
	if p.Current == nil {
		return nil, exception.NewSyntaxException("Expected value", p.Tokens[len(p.Tokens)-1].Loc)
	}

	start := p.Current.Loc

	switch {
	case p.Current.MatchType(TT_STRING):
		node := &StringValueNode{Value: p.Current.Value, Loc: p.Current.Loc.Copy()}
		p.Next()
		return node, nil

	case p.Current.MatchType(TT_NUMBER):
		node := &NumberValueNode{Value: p.Current.Value, Loc: p.Current.Loc.Copy()}
		p.Next()
		return node, nil

	case p.Current.MatchType(TT_IDENT):
		if p.Current.Value == "true" || p.Current.Value == "false" {
			node := &BooleanValueNode{Value: p.Current.Value, Loc: p.Current.Loc.Copy()}
			p.Next()
			return node, nil
		}

		if p.Current.Value == "null" {
			node := &NullValueNode{Loc: p.Current.Loc.Copy()}
			p.Next()
			return node, nil
		}

		return nil, exception.NewSyntaxException("Expected value literal", p.Current.Loc)

	case p.Current.MatchType(TT_LSQUARE):
		p.Next()
		values := make([]ASTValueNode, 0)

		if p.Current.MatchType(TT_RSQUARE) {
			loc := NewLocation(start.File, start.Start, p.Current.Loc.End)
			p.Next()
			return &ArrayValueNode{Values: values, Loc: loc}, nil
		}

		for p.Current != nil && !p.Current.MatchType(TT_EOF) {
			value, err := p.ParseValue()
			if err != nil {
				return nil, err
			}
			values = append(values, value)

			if p.Current.MatchType(TT_COMMA) {
				p.Next()
				if p.Current.MatchType(TT_RSQUARE) {
					return nil, exception.NewSyntaxException("Trailing comma in array is not allowed", p.Current.Loc)
				}
				continue
			}

			if p.Current.MatchType(TT_RSQUARE) {
				loc := NewLocation(start.File, start.Start, p.Current.Loc.End)
				p.Next()
				return &ArrayValueNode{Values: values, Loc: loc}, nil
			}

			return nil, exception.NewSyntaxException("Expected ',' or ']' in array", p.Current.Loc)
		}

		return nil, exception.NewSyntaxException("Unterminated array, expected ']'", p.Current.Loc)

	default:
		return nil, exception.NewSyntaxException("Expected value", p.Current.Loc)
	}
}

func (p *Parser) ParseArgs() ([]ASTValueNode, exception.IException) {
	if p.Current == nil || !p.Current.MatchType(TT_LPAREN) {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected '(' to start argument list", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected '(' to start argument list", p.Current.Loc)
	}

	args := make([]ASTValueNode, 0)
	p.Next()

	if p.Current == nil {
		return nil, exception.NewSyntaxException("Unterminated argument list, expected ')'", p.Tokens[len(p.Tokens)-1].Loc)
	}

	if p.Current.MatchType(TT_RPAREN) {
		p.Next()
		return args, nil
	}

	for p.Current != nil && !p.Current.MatchType(TT_EOF) {
		value, err := p.ParseValue()
		if err != nil {
			return nil, err
		}
		args = append(args, value)

		if p.Current == nil {
			return nil, exception.NewSyntaxException("Unterminated argument list, expected ')'", p.Tokens[len(p.Tokens)-1].Loc)
		}

		if p.Current.MatchType(TT_COMMA) {
			p.Next()
			if p.Current == nil {
				return nil, exception.NewSyntaxException("Unterminated argument list, expected value after ','", p.Tokens[len(p.Tokens)-1].Loc)
			}
			if p.Current.MatchType(TT_RPAREN) {
				return nil, exception.NewSyntaxException("Trailing comma in arguments is not allowed", p.Current.Loc)
			}
			continue
		}

		if p.Current.MatchType(TT_RPAREN) {
			p.Next()
			return args, nil
		}

		return nil, exception.NewSyntaxException("Expected ',' or ')' in argument list", p.Current.Loc)
	}

	if p.Current == nil {
		return nil, exception.NewSyntaxException("Unterminated argument list, expected ')'", p.Tokens[len(p.Tokens)-1].Loc)
	}

	return nil, exception.NewSyntaxException("Unterminated argument list, expected ')'", p.Current.Loc)
}

func (p *Parser) ParseAnnotation() (*AnnotationNode, exception.IException) {
	if p.Current == nil || !p.Current.MatchType(TT_DECORATOR) {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected '@' for annotation", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected '@' for annotation", p.Current.Loc)
	}

	start := p.Current.Loc
	p.Next()

	if p.Current == nil || !p.Current.MatchType(TT_IDENT) {
		if p.Current == nil {
			return nil, exception.NewSyntaxException("Expected annotation name", p.Tokens[len(p.Tokens)-1].Loc)
		}
		return nil, exception.NewSyntaxException("Expected annotation name", p.Current.Loc)
	}

	node := &AnnotationNode{
		Name: &IdentNode{Value: p.Current.Value, Loc: p.Current.Loc.Copy()},
		Args: make([]ASTValueNode, 0),
	}

	p.Next()
	if p.Current != nil && p.Current.MatchType(TT_LPAREN) {
		args, err := p.ParseArgs()
		if err != nil {
			return nil, err
		}

		node.Args = append(node.Args, args...)
	}

	end := start.End
	if len(node.Args) > 0 {
		end = node.Args[len(node.Args)-1].GetLocation().End
	}
	if p.Index-1 >= 0 {
		prev := p.Tokens[p.Index-1]
		if prev != nil && (prev.MatchType(TT_RPAREN) || prev.MatchType(TT_IDENT)) {
			end = prev.Loc.End
		}
	}

	node.Loc = NewLocation(start.File, start.Start, end)
	return node, nil
}

func (p *Parser) Parse() (*ProgramNode, exception.IException) {
	program := &ProgramNode{
		Body: make([]ASTNode, 0),
	}

	p.Next()

	for p.Current != nil && !p.Current.MatchType(TT_EOF) {
		switch {
		case p.Current.MatchType(TT_NEWLINE):
			p.SkipNewLine()

		case p.Current.MatchType(TT_DECORATOR):
			annotations := make([]*AnnotationNode, 0)
			for p.Current != nil && p.Current.MatchType(TT_DECORATOR) {
				annotation, err := p.ParseAnnotation()
				if err != nil {
					return nil, err
				}
				annotations = append(annotations, annotation)
				p.SkipNewLine()
			}

			if p.Current == nil || !p.Current.Match(TT_IDENT, "model") {
				return nil, exception.NewSyntaxException("Annotation must be followed by a model", annotations[len(annotations)-1].Loc)
			}

			n, err := p.ParseModelDecl()
			if err != nil {
				return nil, err
			}

			n.Annotations = append(n.Annotations, annotations...)
			program.Body = append(program.Body, n)

		case p.Current.Match(TT_IDENT, "model"):
			n, err := p.ParseModelDecl()
			if err != nil {
				return nil, err
			}

			program.Body = append(program.Body, n)

		case p.Current.Match(TT_IDENT, "enum"):
			n, err := p.ParseEnumDecl()
			if err != nil {
				return nil, err
			}

			program.Body = append(program.Body, n)

		case p.Current.Match(TT_IDENT, "rest"):
			n, err := p.ParseRestDecl()
			if err != nil {
				return nil, err
			}

			program.Body = append(program.Body, n)
		case p.Current.Match(TT_IDENT, "event"):
			n, err := p.ParseEventDecl()
			if err != nil {
				return nil, err
			}

			program.Body = append(program.Body, n)

		case p.Current.Match(TT_IDENT, "error"):
			n, err := p.ParseErrorDecl()
			if err != nil {
				return nil, err
			}

			program.Body = append(program.Body, n)

		case p.Current.MatchType(TT_EOF):
			p.Next()

		default:
			return nil, exception.NewSyntaxException("Unexpected token at top level", p.Current.Loc)

		}

	}

	return program, nil
}

func (p *Parser) Next() {
	if p.Index+1 < len(p.Tokens) {
		p.Index++
		p.Current = p.Tokens[p.Index]
	} else {
		p.Current = nil
	}
}

func (p *Parser) Peek() *Token {
	if p.Index+1 < len(p.Tokens) {
		return p.Tokens[p.Index+1]
	}
	return nil
}

func NewParser(file string, tokens TokenList) *Parser {
	return &Parser{File: file, Tokens: tokens, Current: nil, Index: -1}
}
