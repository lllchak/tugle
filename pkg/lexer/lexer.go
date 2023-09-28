package tugle

type void struct{}

var nothing void

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
	Value string
	Type  TTokenType
	Loc   TTokenLocation
}

type TCursor struct {
	CurrPos uint
	Loc     TTokenLocation
}

func (token *TToken) equal(other *TToken) bool {
	return token.Value == other.Value && token.Type == other.Type
}

func checkReservedToken(source string, inputCursor TCursor) (*TToken, TCursor, bool) {
	curr := inputCursor

	reservedTokens := []TReservedToken{
		SelectToken,
		FromToken,
		CreateToken,
		TableToken,
		AsToken,
		InsertToken,
		IntoToken,
		ValuesToken,
		IntToken,
		TextToken,
	}

	match := matchBestOption(source, inputCursor, reservedTokens)
	if uint(len(match)) == 0 {
		return nil, inputCursor, false
	}

	matchLen := uint(len(match))

	curr.CurrPos = inputCursor.CurrPos + matchLen
	curr.Loc.Column = inputCursor.Loc.Column + matchLen

	return &TToken{Value: match, Type: ReservedType, Loc: inputCursor.Loc}, curr, matchLen > 0
}

func checkNumeric(source string, inputCursor TCursor) (*TToken, TCursor, bool) {
	curr := inputCursor

	hasMantissa := false
	hasExponent := false

	for ; curr.CurrPos < uint(len(source)); curr.CurrPos++ {
		currChar := source[curr.CurrPos]

		if !checkStart(currChar, inputCursor, curr, &hasMantissa) {
			return nil, inputCursor, false
		} else if !checkMantissa(currChar, inputCursor, curr, &hasMantissa) {
			return nil, inputCursor, false
		} else if !checkExponential(source, inputCursor, curr, &hasMantissa, &hasExponent) {
			return nil, inputCursor, false
		}

		if !(currChar >= '0' && currChar <= '9') {
			break
		}
	}

	if curr.CurrPos == inputCursor.CurrPos {
		return nil, inputCursor, false
	}

	return &TToken{
		Value: source[inputCursor.CurrPos:curr.CurrPos],
		Type:  NumericType,
		Loc:   inputCursor.Loc,
	}, curr, false
}

type apply func(string, TCursor) (*TToken, TCursor, bool)
