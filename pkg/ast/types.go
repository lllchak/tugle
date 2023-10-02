package ast

import "pkg/lexer"

type TStatementType uint
type TExpressionType uint

const (
	LiteralType TExpressionType = iota
)

const (
	SelectType TStatementType = iota
	CreateTableType
	InsertType
)

type TColumnMeta struct {
	Name     lexer.TToken
	Datatype lexer.TToken
}

type TExpression struct {
	Literal *lexer.TToken
	Type    TExpressionType
}

type TInsertStatement struct {
	Table  lexer.TToken
	Values *[]*TExpression
}

type TCreateTableStatement struct {
	TableName lexer.TToken
	Columns   *[]*TColumnMeta
}

type TSelectStatement struct {
	From  *lexer.TToken
	Rules []*TExpression
}

type Statement struct {
	Insert *TInsertStatement
	Type   TStatementType
}

type TSyntaxTree struct {
	Statements []*Statement
}
