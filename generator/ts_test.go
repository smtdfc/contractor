package generator

import (
	"strings"
	"testing"

	"github.com/smtdfc/contractor/parser"
)

func TestTypescriptGenerator_Generate(t *testing.T) {
	mockAST := &parser.AST{
		Statements: []parser.Node{
			&parser.ModelStatementNode{
				Name: "User",
				TypeVar: &parser.TypeVarNode{Name: "T"},
				Annotations: &parser.AnnotationChainNode{
					List: []*parser.AnnotationNode{
						{Name: "Data"},
						{Name: "CreateConstructor"},
					},
				},
				Fields: []*parser.ModelFieldNode{
					{
						Name: "id",
						Type: &parser.TypeDeclarationNode{Name: "Number"},
						Annotations: &parser.AnnotationChainNode{
							List: []*parser.AnnotationNode{
								{Name: "Private"},
							},
						},
					},
					{
						Name: "email",
						Type: &parser.TypeDeclarationNode{Name: "String"},
						Annotations: &parser.AnnotationChainNode{
							List: []*parser.AnnotationNode{
								{
									Name: "IsEmail",
									Args: []parser.Node{
										&parser.LiteralNode{Value: "Invalid email format"},
									},
								},
							},
						},
					},
					{
						Name: "extra",
						// Using the generic type T
						Type: &parser.TypeDeclarationNode{Name: "T"},
						Annotations: &parser.AnnotationChainNode{
							List: []*parser.AnnotationNode{
								{Name: "Optional"},
							},
						},
					},
				},
			},
		},
	}

	gen := NewTypescriptGenerator()
	output, err := gen.Generate(mockAST)

	if err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	expectedElements := []string{
		"import {ContractorRuntime} from 'contractor';",
		"export class User<T>",
		"private id: number;",
		"public email: string;",
		"public extra?: T;",
		"constructor(",
		"public getEmail(): string",
		"public setEmail(v: string): void",
		"public static validate(obj: any)",
		"ContractorRuntime.Validators.IsEmail(obj.email)",
		"Invalid email format",
	}

	for _, element := range expectedElements {
		if !strings.Contains(output, element) {
			t.Errorf("Missing expected element in output: %s\nActual Output:\n%s", element, output)
		}
	}
}

func TestGenerateType_PrimitiveConversion(t *testing.T) {
	gen := NewTypescriptGenerator()

	tests := []struct {
		name     string
		node     parser.Node
		expected string
	}{
		{
			"Convert String to string",
			&parser.TypeDeclarationNode{Name: "String"},
			"string",
		},
		{
			"Convert Number to number",
			&parser.TypeDeclarationNode{Name: "Number"},
			"number",
		},
		{
			"Generic Array",
			&parser.TypeDeclarationNode{
				Name:    "Array",
				Generic: &parser.TypeDeclarationNode{Name: "String"},
			},
			"Array<string>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gen.GenerateType(tt.node)
			if err != nil {
				t.Errorf("GenerateType() error = %v", err)
				return
			}
			if got != tt.expected {
				t.Errorf("GenerateType() = %v, want %v", got, tt.expected)
			}
		})
	}
}
