package lexer

type void struct{}

var nothing void

type TReservedToken string
type TSymbolToken string
type ETokenType uint

type TTokenLocation struct {
	Line   uint
	Column uint
}

const (
	SelectToken TReservedToken = "select"
	FromToken   TReservedToken = "from"
	CreateToken TReservedToken = "create"
	TableToken  TReservedToken = "table"
	AsToken     TReservedToken = "as"
	InsertToken TReservedToken = "insert"
	IntoToken   TReservedToken = "into"
	ValuesToken TReservedToken = "values"
	IntToken    TReservedToken = "int"
	TextToken   TReservedToken = "text"
)

const (
	SemicolonToken    TSymbolToken = ";"
	AsteriksToken     TSymbolToken = "*"
	CommaToken        TSymbolToken = ","
	LeftParenthToken  TSymbolToken = "("
	RightParenthToken TSymbolToken = ")"
)

const (
	ReservedType ETokenType = iota
	SymbolType
	IdentifierType
	StringType
	NumericType
)

type TToken struct {
	Value string
	Type  ETokenType
	Loc   TTokenLocation
}

type TCursor struct {
	CurrPos uint
	Loc     TTokenLocation
}
