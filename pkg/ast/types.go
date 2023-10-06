package ast

import "pkg/lexer"

type EStatementType uint
type EExpressionType uint

const (
	LiteralType EExpressionType = iota
)

const (
	SelectType EStatementType = iota
	CreateTableType
	InsertType
)

type TColumnMeta struct {
	Name     lexer.TToken
	Datatype lexer.TToken
}

type TExpression struct {
	Literal *lexer.TToken
	Type    EExpressionType
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
	From  lexer.TToken
	Rules []*TExpression
}

type TStatement struct {
	CreateTable *TCreateTableStatement
	Select      *TSelectStatement
	Insert      *TInsertStatement
	Type        EStatementType
}

type TSyntaxTree struct {
	Statements []*TStatement
}
