package lexer

import "strings"

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
