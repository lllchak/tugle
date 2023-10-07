package backend

import "pkg/ast"

type Backend interface {
	CreateTable(*ast.TCreateTableStatement) error
	Insert(*ast.TInsertStatement) error
	Select(*ast.TSelectStatement) (*TResult, error)
}
