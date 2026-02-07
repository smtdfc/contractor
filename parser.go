package main

// "fmt"

type Parser struct{}

func (p *Parser) CreateNodeLocation(start *Token, end *Token) *NodeLocation {
	return NewNodeLocation(
		start.Loc,
		end.Loc,
	)
}

func (p *Parser) ExpectToken(scanner *TokenScanner, t TokenType, value string, msg string) (*Token, BaseError) {
	if !scanner.Current.Match(t, value) {
		return nil, NewSyntaxError(
			msg,
			scanner.GetErrorLocation(),
		)
	}

	curr := scanner.Current.Copy()
	scanner.Next()
	return curr, nil
}

func (p *Parser) SkipNewlines(scanner *TokenScanner) {
	for scanner.Current != nil && scanner.Current.HasType(TOKEN_NEWLINE) {
		scanner.Next()
	}
}

func (p *Parser) ExpectTokenTypeAfterNewlines(scanner *TokenScanner, t TokenType, msg string) (*Token, BaseError) {
	p.SkipNewlines(scanner)
	return p.ExpectTokenType(scanner, t, msg)
}

func (p *Parser) ExpectTokenType(scanner *TokenScanner, t TokenType, msg string) (*Token, BaseError) {
	if !scanner.Current.HasType(t) {
		return nil, NewSyntaxError(
			msg,
			scanner.GetErrorLocation(),
		)
	}

	curr := scanner.Current.Copy()
	scanner.Next()
	return curr, nil
}

func (p *Parser) ParseLiteral(scanner *TokenScanner) (Node, BaseError) {
	curr := scanner.Current
	switch curr.Type {
	case TOKEN_NUMBER, TOKEN_STRING, TOKEN_BOOL, TOKEN_NULL:
		node := NewLiteralNode(curr.Type, curr.Value, p.CreateNodeLocation(curr, curr))
		scanner.Next()
		return node, nil
	default:
		return nil, NewSyntaxError(
			"Only numbers, strings, booleans, or null are allowed here",
			scanner.GetErrorLocation(),
		)
	}
}

func (p *Parser) ParseAnnotation(scanner *TokenScanner) (*AnnotationNode, BaseError) {
	startToken := scanner.Current.Copy()
	scanner.Next()
	args := ListNode{}
	name, err := p.ExpectTokenType(
		scanner,
		TOKEN_IDENTIFIER,
		"Expected annotation name",
	)
	if err != nil {
		return nil, err
	}

	if scanner.Current.HasType(TOKEN_LEFT_PAREN) {
		scanner.Next()

		for !scanner.Current.HasType(TOKEN_RIGHT_PAREN) {
			arg, err := p.ParseLiteral(scanner)
			if err != nil {
				return nil, err
			}
			args = append(args, arg)

			if scanner.Current.HasType(TOKEN_COMMA) {
				scanner.Next()

				if scanner.Current.HasType(TOKEN_RIGHT_PAREN) {
					break
				}
			} else if !scanner.Current.HasType(TOKEN_RIGHT_PAREN) {
				return nil, NewSyntaxError(
					"Expected ',' or ')' after annotation argument",
					scanner.GetErrorLocation(),
				)
			}
		}

		_, err = p.ExpectTokenType(scanner, TOKEN_RIGHT_PAREN, "Expected ')' after arguments")
		if err != nil {
			return nil, err
		}
	}

	if scanner.Current.HasType(TOKEN_NEWLINE) {
		scanner.Next()
	} else if scanner.Current.HasType(TOKEN_EOF) {
		// ignore
	} else {
		return nil, NewSyntaxError(
			"Invalid syntax",
			scanner.GetErrorLocation(),
		)
	}

	loc := p.CreateNodeLocation(startToken, name)

	return NewAnnotationNode(
		name.Value,
		args,
		loc,
	), nil
}

func (p *Parser) ParseAnnotationChain(scanner *TokenScanner) (*AnnotationChainNode, BaseError) {
	var annotations ListAnnotation
	var endToken *Token
	startToken := scanner.Current.Copy()

	first, err := p.ParseAnnotation(scanner)
	if err != nil {
		return nil, err
	}
	endToken = scanner.Current.Copy()

	annotations = append(annotations, first)
	for scanner.Current != nil && !scanner.Current.HasType(TOKEN_EOF) {
		if scanner.Current.HasType(TOKEN_ANNOTATION) {
			anno, err := p.ParseAnnotation(scanner)
			if err != nil {
				return nil, err
			}
			annotations = append(annotations, anno)
			endToken = scanner.Current.Copy()
			continue
		} else if scanner.Current.HasType(TOKEN_NEWLINE) {
			scanner.Next()
			continue
		} else {
			break
		}
	}

	loc := p.CreateNodeLocation(startToken, endToken)
	return NewAnnotationChainNode(
		annotations,
		loc,
	), nil
}

func (p *Parser) ParseType(scanner *TokenScanner) (*TypeDeclarationNode, BaseError) {
	var endToken *Token
	startToken, err := p.ExpectTokenType(scanner, TOKEN_IDENTIFIER, "Expect type name")
	if err != nil {
		return nil, err
	}

	endToken = startToken
	var generic Node
	if scanner.Current.Match(TOKEN_OPERATOR, "<") {
		scanner.Next()
		generic, err = p.ParseType(scanner)
		if err != nil {
			return nil, err
		}

		endToken, err = p.ExpectToken(scanner, TOKEN_OPERATOR, ">", "Expected '>'")
		if err != nil {
			return nil, err
		}
	}

	return NewTypeDeclarationNode(
		startToken.Value,
		generic,
		p.CreateNodeLocation(startToken, endToken),
	), nil
}

func (p *Parser) ParseModelField(scanner *TokenScanner, annotations *AnnotationChainNode) (*ModelFieldNode, BaseError) {
	startToken := scanner.Current.Copy()

	typeNode, err := p.ParseType(scanner)
	if err != nil {
		return nil, err
	}

	name, err := p.ExpectTokenTypeAfterNewlines(
		scanner,
		TOKEN_IDENTIFIER,
		"Expected field name",
	)
	if err != nil {
		return nil, err
	}

	p.SkipNewlines(scanner)
	loc := p.CreateNodeLocation(startToken, name)
	return NewModelFieldNode(
		name.Value,
		annotations,
		typeNode,
		loc,
	), nil
}

func (p *Parser) ParseAllModelFields(scanner *TokenScanner) ([]*ModelFieldNode, BaseError) {
	var annotationChain *AnnotationChainNode
	var fields []*ModelFieldNode

	for scanner.Current != nil && !scanner.Current.HasType(TOKEN_EOF) {
		p.SkipNewlines(scanner)

		if scanner.Current.HasType(TOKEN_RIGHT_BRACE) {
			break
		}

		if scanner.Current.HasType(TOKEN_ANNOTATION) {
			node, err := p.ParseAnnotationChain(scanner)
			if err != nil {
				return nil, err
			}
			annotationChain = node
			continue
		}

		field, err := p.ParseModelField(scanner, annotationChain)
		if err != nil {
			return nil, err
		}

		fields = append(fields, field)
		annotationChain = nil
	}

	if annotationChain != nil {
		return nil, NewSyntaxError(
			"Dangling annotation: no field affected",
			annotationChain.Loc.GetErrorLocation(),
		)
	}

	return fields, nil
}

func (p *Parser) ParseModelStatement(scanner *TokenScanner, annotations *AnnotationChainNode) (*ModelStatementNode, BaseError) {
	startToken := scanner.Current.Copy()
	scanner.Next()

	name, err := p.ExpectTokenType(
		scanner,
		TOKEN_IDENTIFIER,
		"Expected model name",
	)
	if err != nil {
		return nil, err
	}

	var typeVar *TypeVarNode
	if scanner.Current.Match(TOKEN_OPERATOR, "<") {
		scanner.Next()
		n, err := p.ExpectTokenType(scanner, TOKEN_IDENTIFIER, "Expected type var")
		if err != nil {
			return nil, err
		}

		typeVar = NewTypeVarNode(n.Value, p.CreateNodeLocation(n, n))
		_, err = p.ExpectToken(scanner, TOKEN_OPERATOR, ">", "Expected '>'")
		if err != nil {
			return nil, err
		}
	}

	_, err = p.ExpectTokenTypeAfterNewlines(
		scanner,
		TOKEN_LEFT_BRACE,
		"Expected {",
	)
	if err != nil {
		return nil, err
	}

	fields, err := p.ParseAllModelFields(scanner)
	if err != nil {
		return nil, err
	}

	//end
	endToken, err := p.ExpectTokenTypeAfterNewlines(
		scanner,
		TOKEN_RIGHT_BRACE,
		"Expected }",
	)
	if err != nil {
		return nil, err
	}

	loc := p.CreateNodeLocation(startToken, endToken)
	return NewModelStatementNode(
		name.Value,
		typeVar,
		annotations,
		fields,
		loc,
	), nil
}

func (p *Parser) Parse(tokens ListToken, file string) (*AST, BaseError) {
	var statements ListNode
	var annotationChain *AnnotationChainNode
	scanner := NewTokenScanner(tokens)
	scanner.Next()
	startToken := scanner.Current.Copy()

	for scanner.Current != nil && !scanner.Current.HasType(TOKEN_EOF) {
		if scanner.Current.HasType(TOKEN_ANNOTATION) {
			node, err := p.ParseAnnotationChain(scanner)
			annotationChain = node
			if err != nil {
				return nil, err
			}
			continue
		}

		if scanner.Current.Match(TOKEN_KEYWORD, "model") {
			node, err := p.ParseModelStatement(scanner, annotationChain)
			if err != nil {
				return nil, err
			}

			annotationChain = nil
			statements = append(statements, node)
		}

		if scanner.Current != nil {
			if scanner.Current.HasType(TOKEN_NEWLINE) {
				scanner.Next()
			} else if scanner.Current.HasType(TOKEN_EOF) {
				if annotationChain != nil {
					return nil, NewSyntaxError(
						"Dangling annotation: no object affected",
						annotationChain.Loc.GetErrorLocation(),
					)
				}

				break
			} else {
				return nil, NewSyntaxError(
					"Invalid syntax",
					scanner.GetErrorLocation(),
				)
			}
		}
	}

	loc := p.CreateNodeLocation(startToken, scanner.Current)
	return NewAST(
		statements,
		loc,
	), nil
}

func NewParser() *Parser {
	return &Parser{}
}
