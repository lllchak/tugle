package backend

type EColumnType uint

const (
	TextType EColumnType = iota
	IntType
)

type ICell interface {
	AsText() string
	AsInt() int64
}

type TResultColumn struct {
	Type EColumnType
	Name string
}

type TResult struct {
	Columns []TResultColumn
	Rows    [][]ICell
}
