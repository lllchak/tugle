package parser

import (
	"fmt"
	"pkg/ast"
	"pkg/lexer"
)

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
	ttype lexer.ETokenType,
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

	types := []lexer.ETokenType{lexer.IdentifierType, lexer.NumericType, lexer.StringType}

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

func parseColumnMeta(
	tokens []*lexer.TToken,
	inputCursor uint,
	delimeter lexer.TToken,
) (*[]*ast.TColumnMeta, uint, bool) {
	curr := inputCursor

	columnsMeta := []*ast.TColumnMeta{}

	for {
		if curr >= uint(len(tokens)) {
			return nil, inputCursor, false
		}

		if currToken := tokens[curr]; delimeter.Equal(currToken) {
			break
		}

		if len(columnsMeta) > 0 {
			var ok bool
			_, curr, ok = parseToken(tokens, curr, *lexer.CommaToken.AsToken())
			if !ok {
				logInfo(tokens, curr, "Expected comma")
				return nil, inputCursor, false
			}
		}

		columnName, currCursor, ok := parseTokenType(tokens, curr, lexer.IdentifierType)
		if !ok {
			logInfo(tokens, curr, "Expected column name")
			return nil, inputCursor, false
		}
		curr = currCursor

		columnType, currCursor, ok := parseTokenType(tokens, curr, lexer.ReservedType)
		if !ok {
			logInfo(tokens, curr, "Expected column datatype definition")
			return nil, inputCursor, false
		}
		curr = currCursor

		columnsMeta = append(
			columnsMeta,
			&ast.TColumnMeta{Name: *columnName, Datatype: *columnType},
		)
	}

	return &columnsMeta, curr, true
}

func parseCreateTableStatement(
	tokens []*lexer.TToken,
	inputCursor uint,
	delimeter lexer.TToken,
) (*ast.TCreateTableStatement, uint, bool) {
	curr := inputCursor
	ok := false

	_, curr, ok = parseToken(tokens, curr, *lexer.CreateToken.AsToken())
	if !ok {
		return nil, inputCursor, false
	}

	_, curr, ok = parseToken(tokens, curr, *lexer.TableToken.AsToken())
	if !ok {
		return nil, inputCursor, false
	}

	tableName, currCursor, ok := parseTokenType(tokens, curr, lexer.IdentifierType)
	if !ok {
		return nil, inputCursor, false
	}
	curr = currCursor

	_, curr, ok = parseToken(tokens, curr, *lexer.LeftParenthToken.AsToken())
	if !ok {
		return nil, inputCursor, false
	}

	columnsDesc, currCursor, ok := parseColumnMeta(tokens, curr, *lexer.RightParenthToken.AsToken())
	if !ok {
		return nil, inputCursor, false
	}
	curr = currCursor

	_, curr, ok = parseToken(tokens, curr, *lexer.RightParenthToken.AsToken())
	if !ok {
		return nil, inputCursor, false
	}

	return &ast.TCreateTableStatement{
		TableName: *tableName,
		Columns:   columnsDesc,
	}, curr, true
}

func parseSelectStatement(
	tokens []*lexer.TToken,
	inputCursor uint,
	delimeter lexer.TToken,
) (*ast.TSelectStatement, uint, bool) {
	curr := inputCursor
	ok := false

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

		resStatement.From = *from
		curr = currCursor
	}

	return &resStatement, curr, true
}

func parseInsertStatement(
	tokens []*lexer.TToken,
	inputCursor uint,
	_ lexer.TToken,
) (*ast.TInsertStatement, uint, bool) {
	curr := inputCursor
	ok := false

	_, curr, ok = parseToken(tokens, curr, *lexer.InsertToken.AsToken())
	if !ok {
		return nil, inputCursor, ok
	}

	_, curr, ok = parseToken(tokens, curr, *lexer.IntoToken.AsToken())
	if !ok {
		return nil, inputCursor, ok
	}

	tableName, currCursor, ok := parseTokenType(tokens, curr, lexer.IdentifierType)
	if !ok {
		logInfo(tokens, curr, "Expected table name")
		return nil, inputCursor, ok
	}
	curr = currCursor

	_, curr, ok = parseToken(tokens, curr, *lexer.ValuesToken.AsToken())
	if !ok {
		logInfo(tokens, curr, "Expected VALUES statement")
		return nil, inputCursor, ok
	}

	_, curr, ok = parseToken(tokens, curr, *lexer.LeftParenthToken.AsToken())
	if !ok {
		logInfo(tokens, curr, "Expected expressions group opening (maybe you forgot opening parenthesis)")
		return nil, inputCursor, ok
	}

	values, currCursor, ok := parseExpressions(tokens, curr, []lexer.TToken{*lexer.RightParenthToken.AsToken()})
	if !ok {
		logInfo(tokens, curr, "Expected values")
		return nil, inputCursor, ok
	}
	curr = currCursor

	_, curr, ok = parseToken(tokens, curr, *lexer.RightParenthToken.AsToken())
	if !ok {
		logInfo(tokens, curr, "Expression way never closed")
		return nil, inputCursor, ok
	}

	return &ast.TInsertStatement{
		Table:  *tableName,
		Values: values,
	}, curr, ok
}

func parseStatement(
	tokens []*lexer.TToken,
	inputCursor uint,
	delimeter lexer.TToken,
) (*ast.TStatement, uint, bool) {
	curr := inputCursor
	semicolonToken := lexer.SemicolonToken.AsToken()

	if selectStatement, currCursor, ok := parseSelectStatement(tokens, curr, *semicolonToken); ok {
		return &ast.TStatement{
			Select: selectStatement,
			Type:   ast.SelectType,
		}, currCursor, ok
	}

	if insertStatemt, currCursor, ok := parseInsertStatement(tokens, curr, *semicolonToken); ok {
		return &ast.TStatement{
			Insert: insertStatemt,
			Type:   ast.InsertType,
		}, currCursor, ok
	}

	if createTableStatement, currCursor, ok := parseCreateTableStatement(tokens, curr, *semicolonToken); ok {
		return &ast.TStatement{
			CreateTable: createTableStatement,
			Type:        ast.CreateTableType,
		}, currCursor, ok
	}

	return nil, inputCursor, false
}
