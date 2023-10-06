package lexer

import (
	"fmt"
	"strings"
)

func CheckSymbol(source string, inputCursor TCursor) (*TToken, TCursor, bool) {
	if uint(len(source)) == 0 {
		return nil, inputCursor, false
	}

	currChar := source[inputCursor.CurrPos]
	curr := inputCursor

	curr.CurrPos++
	curr.Loc.Column++

	switch currChar {
	case '\n':
		curr.Loc.Line++
		curr.Loc.Column = 0
		fallthrough
	case '\t':
		fallthrough
	default:
		if matchRegex([]byte{currChar}, "\\s") {
			return nil, curr, true
		}
	}

	symbols := []TSymbolToken{
		SemicolonToken,
		AsteriksToken,
		CommaToken,
		LeftParenthToken,
		RightParenthToken,
	}

	match := matchBestOption(source, inputCursor, getStringRerp(symbols))
	matchLen := uint(len(match))

	if matchLen == 0 {
		return nil, inputCursor, false
	}

	curr.CurrPos = inputCursor.CurrPos + matchLen
	curr.Loc.Column = inputCursor.Loc.Column + matchLen

	return &TToken{Value: match, Type: SymbolType, Loc: inputCursor.Loc}, curr, matchLen > 0
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

	match := matchBestOption(source, inputCursor, getStringRerp(reservedTokens))
	matchLen := uint(len(match))

	if matchLen == 0 {
		return nil, inputCursor, false
	}

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

		isDigit := isNumeric(currChar)
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

func CheckIdentifier(source string, inputCursor TCursor) (*TToken, TCursor, bool) {
	if token, curr, ok := checkDelimeted(source, inputCursor, '"'); ok {
		token.Type = IdentifierType
		return token, curr, ok
	}

	curr := inputCursor
	currChar := source[curr.CurrPos]

	isAlpha := isLetter(currChar)
	if !isAlpha {
		return nil, inputCursor, false
	}
	curr.CurrPos++
	curr.Loc.Column++

	match := []byte{currChar}
	for ; curr.CurrPos < uint(len(source)); curr.CurrPos++ {
		currChar := source[curr.CurrPos]

		isAlpha := isLetter(currChar)
		isNumeric := isNumeric(currChar)

		if isAlpha || isNumeric || currChar == '$' || currChar == '_' {
			match = append(match, currChar)
			curr.Loc.Column++
			continue
		}

		break
	}

	return &TToken{
		Value: strings.ToLower(string(match)),
		Type:  IdentifierType,
		Loc:   inputCursor.Loc,
	}, curr, true
}

func CheckString(source string, inputCursor TCursor) (*TToken, TCursor, bool) {
	return checkDelimeted(source, inputCursor, '\'')
}

type apply func(string, TCursor) (*TToken, TCursor, bool)

func Tokenize(source string) ([]*TToken, error) {
	tokens := []*TToken{}
	curr := TCursor{}

Tokenize:
	for curr.CurrPos < uint(len(source)) {
		lexers := []apply{CheckReservedToken, CheckSymbol, CheckString, CheckNumeric, CheckIdentifier}

		for _, lexer := range lexers {
			if token, currCursor, ok := lexer(source, curr); ok {
				curr = currCursor

				if token != nil {
					tokens = append(tokens, token)
				}

				continue Tokenize
			}
		}

		hint := ""
		if len(tokens) > 0 {
			hint = "after " + tokens[len(tokens)-1].Value
		}
		return nil, fmt.Errorf("Unable to lex token %s, at %d:%d", hint, curr.Loc.Line, curr.Loc.Column)
	}

	return tokens, nil
}
