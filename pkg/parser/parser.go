package parser

import (
	"errors"
	"pkg/ast"
	"pkg/lexer"
)

func Parse(source string) (*ast.TSyntaxTree, error) {
	tokens, err := lexer.Tokenize(source)
	if err != nil {
		return nil, err
	}

	semicolonToken := lexer.SemicolonToken.AsToken()
	if len(tokens) > 0 && !tokens[len(tokens)-1].Equal(semicolonToken) {
		tokens = append(tokens, semicolonToken)
	}

	syntaxTree := ast.TSyntaxTree{}
	curr := uint(0)

	for curr < uint(len(tokens)) {
		statement, currCursor, ok := parseStatement(tokens, curr, *semicolonToken)
		if !ok {
			logInfo(tokens, curr, "Expected statement")
			return nil, errors.New("Failed to parse, expected statement")
		}
		curr = currCursor

		syntaxTree.Statements = append(syntaxTree.Statements, statement)

		hasSemicolon := false
		for {
			_, curr, hasSemicolon = parseToken(tokens, curr, *semicolonToken)
			if hasSemicolon {
				break
			}
		}

		if !hasSemicolon {
			logInfo(tokens, curr, "Expected semi-colon delimiter between statements")
			return nil, errors.New("Missing semi-colon between statements")
		}
	}

	return &syntaxTree, nil
}
