package main

import (
	"fmt"
	"pkg/ast"
	"pkg/lexer"
	"pkg/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		source string
		ast    *ast.TSyntaxTree
	}{
		{
			source: `SELECT id, name FROM "sketchy name"`,
			ast: &ast.TSyntaxTree{
				Statements: []*ast.TStatement{
					{
						Type: ast.SelectType,
						Select: &ast.TSelectStatement{
							Rules: []*ast.TExpression{
								{
									Type: ast.LiteralType,
									Literal: &lexer.TToken{
										Loc:   lexer.TTokenLocation{Column: 7, Line: 0},
										Type:  lexer.IdentifierType,
										Value: "id",
									},
								},
								{
									Type: ast.LiteralType,
									Literal: &lexer.TToken{
										Loc:   lexer.TTokenLocation{Column: 11, Line: 0},
										Type:  lexer.IdentifierType,
										Value: "name",
									},
								},
							},
							From: lexer.TToken{
								Loc:   lexer.TTokenLocation{Column: 33, Line: 0},
								Type:  lexer.IdentifierType,
								Value: "sketchy name",
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		fmt.Println("(Parser) Testing: ", test.source)
		ast, err := parser.Parse(test.source)
		assert.Nil(t, err, test.source)
		assert.Equal(t, test.ast, ast, test.source)
	}
}
