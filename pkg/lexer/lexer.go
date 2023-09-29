package lexer

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

func CheckReservedToken(source string, inputCursor TCursor) (*TToken, TCursor, bool) {
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

func CheckNumeric(source string, inputCursor TCursor) (*TToken, TCursor, bool) {
	curr := inputCursor

	isFloat := false
	isExponent := false

	for ; curr.CurrPos < uint(len(source)); curr.CurrPos++ {
		currChar := source[curr.CurrPos]
		curr.Loc.Column++

		isDigit := currChar >= '0' && currChar <= '9'
		isPeriod := currChar == '.'
		isExpMarker := currChar == 'e'

		if curr.CurrPos == inputCursor.CurrPos {
			if !isDigit && !isPeriod {
				return nil, inputCursor, false
			}

			isFloat = isPeriod
			continue
		}

		if isPeriod {
			if isFloat {
				return nil, inputCursor, false
			}

			isFloat = true
			continue
		}

		if isExpMarker {
			if isExponent {
				return nil, inputCursor, false
			}

			isFloat = true
			isExponent = true

			if curr.CurrPos == uint(len(source)-1) {
				return nil, inputCursor, false
			}

			nextChar := source[curr.CurrPos+1]
			if nextChar == '-' || nextChar == '+' {
				curr.CurrPos++
				curr.Loc.Column++
			}

			continue
		}

		if !isDigit {
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
	}, curr, true
}

type apply func(string, TCursor) (*TToken, TCursor, bool)
