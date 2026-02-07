package main

import (
	"fmt"
	"strings"
)

func PrintAST(node Node, indent int) {
	prefix := strings.Repeat("  ", indent)

	switch n := node.(type) {
	case *AST:
		fmt.Printf("%s[AST] (%v)\n", prefix, n.Loc)
		for _, stmt := range n.Statements {
			PrintAST(stmt, indent+1)
		}

	case *ModelStatementNode:
		fmt.Printf("%s[Model] Name: %s\n", prefix, n.Name)
		if n.Annotations != nil {
			PrintAST(n.Annotations, indent+1)
		}

		for _, field := range n.Fields {
			PrintAST(field, indent+1)
		}

	case *AnnotationChainNode:
		fmt.Printf("%s[AnnotationChain]\n", prefix)
		for _, anno := range n.List {
			PrintAST(anno, indent+1)
		}

	case *AnnotationNode:
		fmt.Printf("%s@%s\n", prefix, n.Name)

	case *ModelFieldNode:
		fmt.Printf("%s[Field] Name: %s\n", prefix, n.Name)
		if n.Annotations != nil {
			PrintAST(n.Annotations, indent+1)
		}
		PrintAST(n.Type, indent+1)

	case *TypeDeclarationNode:
		fmt.Printf("%s[Type] %s\n", prefix, n.Name)
		if n.Generic != nil {
			fmt.Printf("%s  [Generic]\n", prefix)
			PrintAST(n.Generic, indent+2)
		}

	default:
		fmt.Printf("%s[Unknown Node Type]\n", prefix)
	}
}
