package tugle

import "strings"

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

func matchBestOption(
	source string,
	inputCursor TCursor,
	tokenOptions []TReservedToken,
) string {
	var currentValue []byte
	var res string
	irrelevantLocs := make(map[int]void)

	curr := inputCursor

	for curr.CurrPos < uint(len(source)) {
		currentValue = append(currentValue, strings.ToLower(string(source[curr.CurrPos]))...)
		curr.CurrPos++

	res:
		for i, option := range tokenOptions {
			if _, ok := irrelevantLocs[i]; ok {
				continue res
			}

			optionSV := string(option)
			currentValueSV := string(currentValue)

			if optionSV == currentValueSV {
				irrelevantLocs[i] = nothing
				if len(optionSV) > len(res) {
					res = string(option)
				}
				continue
			}

			samePrefix := currentValueSV == optionSV[:curr.CurrPos-inputCursor.CurrPos]
			tooLong := len(currentValue) > len(option)
			if tooLong || !samePrefix {
				irrelevantLocs[i] = nothing
			}
		}

		if len(irrelevantLocs) == len(tokenOptions) {
			break
		}
	}

	return res
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

		isDigit := currChar >= '0' && currChar <= '9'
		isMantissa := currChar == '.'
		isExponential := currChar == 'e'

		if curr.CurrPos == inputCursor.CurrPos {
			if !isDigit || !isMantissa {
				return nil, inputCursor, false
			}
			hasMantissa = isMantissa

			continue
		}

		if isMantissa {
			if hasMantissa {
				return nil, inputCursor, false
			}
			hasMantissa = true

			continue
		}

		if isExponential {
			if hasExponent {
				return nil, inputCursor, false
			}

		}
	}

	return &TToken{Value: "", Type: ReservedType, Loc: inputCursor.Loc}, curr, false
}

type apply func(string, TCursor) (*TToken, TCursor, bool)
