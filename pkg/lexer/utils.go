package lexer

import (
	"regexp"
	"strings"
)

func matchBestOption(
	source string,
	inputCursor TCursor,
	tokenOptions []string,
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

			if option == string(currentValue) {
				irrelevantLocs[i] = nothing
				if len(option) > len(res) {
					res = option
				}
				continue
			}

			samePrefix := string(currentValue) == option[:curr.CurrPos-inputCursor.CurrPos]
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

func checkDelimeted(source string, inputCursor TCursor, delimeter byte) (*TToken, TCursor, bool) {
	curr := inputCursor

	if len(source) == 0 {
		return nil, inputCursor, false
	}

	if source[curr.CurrPos] != delimeter {
		return nil, inputCursor, false
	}

	curr.CurrPos++
	curr.Loc.Column++

	var resMatch []byte

	for ; curr.CurrPos < uint(len(source)); curr.CurrPos++ {
		currChar := source[curr.CurrPos]

		if currChar == delimeter {
			if curr.CurrPos+1 >= uint(len(source)) || source[curr.CurrPos+1] != delimeter {
				curr.CurrPos++
				curr.Loc.Column++
				return &TToken{Value: string(resMatch), Type: StringType, Loc: inputCursor.Loc}, curr, true
			} else {
				resMatch = append(resMatch, currChar)
				curr.CurrPos++
				curr.Loc.Column++
			}
		}

		resMatch = append(resMatch, currChar)
		curr.Loc.Column++
	}

	return nil, inputCursor, false
}

func getStringRerp(options interface{}) []string {
	var res []string

	switch opts := options.(type) {
	case []TReservedToken:
		for _, token := range opts {
			res = append(res, string(token))
		}
	case []TSymbolToken:
		for _, token := range opts {
			res = append(res, string(token))
		}
	}

	return res
}

func matchRegex(char []byte, pattern string) bool {
	regexp := regexp.MustCompile(pattern)

	return regexp.Match(char)
}

func isNumeric(char byte) bool {
	return char >= '0' && char <= '9'
}

func isLetter(char byte) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')
}

func (token *TToken) Equal(other *TToken) bool {
	return token.Value == other.Value && token.Type == other.Type
}

func (reservedToken TReservedToken) AsToken() *TToken {
	return &TToken{
		Value: string(reservedToken),
		Type:  ReservedType,
	}
}

func (symbolToken TSymbolToken) AsToken() *TToken {
	return &TToken{
		Value: string(symbolToken),
		Type:  StringType,
	}
}
