package parser

import (
	"fmt"
	"pkg/ast"
	"pkg/lexer"
)

func tokenFromReservedToken(reservedToken lexer.TReservedToken) lexer.TToken {
	return lexer.TToken{
		Value: string(reservedToken),
		Type:  lexer.ReservedType,
	}
}

func tokenFromSymbol(symbol lexer.TSymbolToken) lexer.TToken {
	return lexer.TToken{
		Value: string(symbol),
		Type:  lexer.SymbolType,
	}
}

func logInfo(tokens []*lexer.TToken, cursor uint, message string) {
	var curr *lexer.TToken

	if cursor < uint(len(tokens)) {
		curr = tokens[cursor]
	} else {
		curr = tokens[cursor-1]
	}

	fmt.Printf("[%d, %d]: %s, got: %s\n", curr.Loc.Line, curr.Loc.Column, message, curr.Value)
}

func isDelimeter(candidate *lexer.TToken, delimeters *[]lexer.TToken) bool {
	for _, delimeter := range *delimeters {
		if delimeter.Equal(candidate) {
			return true
		}
	}

	return false
}

func parseTokenType(
	tokens []*lexer.TToken,
	initialCursor uint,
	ttype lexer.TTokenType,
) (*lexer.TToken, uint, bool) {
	cursor := initialCursor

	if cursor >= uint(len(tokens)) {
		return nil, initialCursor, false
	}

	if currToken := tokens[cursor]; currToken.Type == ttype {
		return currToken, cursor + 1, true
	}

	return nil, initialCursor, false
}

func parseToken(tokens []*lexer.TToken, inputCursor uint, candToken lexer.TToken) (*lexer.TToken, uint, bool) {
	curr := inputCursor

	if curr >= uint(len(tokens)) {
		return nil, inputCursor, false
	}

	if currToken := tokens[curr]; candToken.Equal(currToken) {
		return currToken, curr + 1, true
	}

	return nil, inputCursor, false
}

func parseExpression(
	tokens []*lexer.TToken,
	inputCursor uint,
	_ lexer.TToken, /* delimeter token in future */
) (*ast.TExpression, uint, bool) {
	curr := inputCursor

	types := []lexer.TTokenType{lexer.IdentifierType, lexer.NumericType, lexer.StringType}

	for _, ttype := range types {
		if currToken, currCursor, ok := parseTokenType(tokens, curr, ttype); ok {
			return &ast.TExpression{
				Literal: currToken,
				Type:    ast.LiteralType,
			}, currCursor, true
		}
	}

	return nil, inputCursor, false
}

func parseExpressions(
	tokens []*lexer.TToken,
	inputCursor uint,
	delimeters []lexer.TToken,
) (*[]*ast.TExpression, uint, bool) {
	curr := inputCursor

	expressions := []*ast.TExpression{}

outer:
	for {
		if curr >= uint(len(tokens)) {
			return nil, inputCursor, false
		}

		currToken := tokens[curr]
		if isDelimeter(currToken, &delimeters) {
			break outer
		}

		if len(expressions) > 0 {
			if _, curr, ok := parseToken(tokens, curr, *lexer.CommaToken.AsToken()); !ok {
				logInfo(tokens, curr, "Maybe you missed comma")
				return nil, inputCursor, false
			}
			curr++
		}

		expression, currCursor, ok := parseExpression(tokens, curr, *lexer.CommaToken.AsToken())
		if !ok {
			logInfo(tokens, curr, "Expected expression")
			return nil, inputCursor, false
		}

		expressions = append(expressions, expression)

		curr = currCursor
	}

	return &expressions, curr, true
}

func parseSelectStatement(
	tokens []*lexer.TToken,
	inputCursor uint,
	delimeter lexer.TToken,
) (*ast.TSelectStatement, uint, bool) {
	var ok bool
	curr := inputCursor

	_, curr, ok = parseToken(tokens, curr, *lexer.SelectToken.AsToken())
	if !ok {
		return nil, inputCursor, false
	}

	resStatement := ast.TSelectStatement{}

	fromToken := *lexer.FromToken.AsToken()
	expressions, curr, ok := parseExpressions(tokens, curr, []lexer.TToken{fromToken, delimeter})
	if !ok {
		return nil, inputCursor, false
	}

	resStatement.Rules = *expressions

	_, curr, ok = parseToken(tokens, curr, fromToken)
	if ok {
		from, currCursor, ok := parseTokenType(tokens, curr, lexer.IdentifierType)
		if !ok {
			logInfo(tokens, curr, "Expected FROM statement")
			return nil, inputCursor, false
		}

		resStatement.From = from
		curr = currCursor
	}

	return &resStatement, curr, true
}
