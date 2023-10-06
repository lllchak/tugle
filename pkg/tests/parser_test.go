package main

import (
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
								Loc:   lexer.TTokenLocation{Column: 21, Line: 0},
								Type:  lexer.IdentifierType,
								Value: "sketchy name",
							},
						},
					},
				},
			},
		},
		{
			source: `SELECT id, name FROM "sketchy name";`,
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
								Loc:   lexer.TTokenLocation{Column: 21, Line: 0},
								Type:  lexer.IdentifierType,
								Value: "sketchy name",
							},
						},
					},
				},
			},
		},
		{
			source: `SELECT 1, 2, 3;`,
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
										Type:  lexer.NumericType,
										Value: "1",
									},
								},
								{
									Type: ast.LiteralType,
									Literal: &lexer.TToken{
										Loc:   lexer.TTokenLocation{Column: 11, Line: 0},
										Type:  lexer.NumericType,
										Value: "2",
									},
								},
								{
									Type: ast.LiteralType,
									Literal: &lexer.TToken{
										Loc:   lexer.TTokenLocation{Column: 15, Line: 0},
										Type:  lexer.NumericType,
										Value: "3",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			source: "INSERT INTO users VALUES (105, string)",
			ast: &ast.TSyntaxTree{
				Statements: []*ast.TStatement{
					{
						Type: ast.InsertType,
						Insert: &ast.TInsertStatement{
							Table: lexer.TToken{
								Loc:   lexer.TTokenLocation{Column: 12, Line: 0},
								Type:  lexer.IdentifierType,
								Value: "users",
							},
							Values: &[]*ast.TExpression{
								{
									Literal: &lexer.TToken{
										Loc:   lexer.TTokenLocation{Column: 26, Line: 0},
										Type:  lexer.NumericType,
										Value: "105",
									}, Type: ast.LiteralType,
								},
								{
									Literal: &lexer.TToken{
										Loc:   lexer.TTokenLocation{Column: 32, Line: 0},
										Type:  lexer.IdentifierType,
										Value: "string",
									}, Type: ast.LiteralType,
								},
							},
						},
					},
				},
			},
		},
		{
			source: "INSERT INTO users VALUES (105, 'string')",
			ast: &ast.TSyntaxTree{
				Statements: []*ast.TStatement{
					{
						Type: ast.InsertType,
						Insert: &ast.TInsertStatement{
							Table: lexer.TToken{
								Loc:   lexer.TTokenLocation{Column: 12, Line: 0},
								Type:  lexer.IdentifierType,
								Value: "users",
							},
							Values: &[]*ast.TExpression{
								{
									Literal: &lexer.TToken{
										Loc:   lexer.TTokenLocation{Column: 26, Line: 0},
										Type:  lexer.NumericType,
										Value: "105",
									}, Type: ast.LiteralType,
								},
								{
									Literal: &lexer.TToken{
										Loc:   lexer.TTokenLocation{Column: 32, Line: 0},
										Type:  lexer.StringType,
										Value: "string",
									}, Type: ast.LiteralType,
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		ast, err := parser.Parse(test.source)
		assert.Nil(t, err, test.source)
		assert.Equal(t, test.ast, ast, test.source)
	}
}
