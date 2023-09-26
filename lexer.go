package tugle

type TReservedToken string
type TPunctuationToken string
type TTokenType uint

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
	SemicolonToken    TPunctuationToken = ";"
	AsteriksToken     TPunctuationToken = "*"
	CommaToken        TPunctuationToken = ","
	LeftBracketToken  TPunctuationToken = "("
	RightBracketToken TPunctuationToken = ")"
)

const (
	ReservedType TTokenType = iota
	PunctuationType
	IdentifierType
	StringType
	NumericType
)

type TToken struct {
	value string
	ttype TTokenType
	loc   TTokenLocation
}
