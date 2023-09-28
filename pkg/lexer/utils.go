package tugle

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

func checkStart(
	currChar byte,
	inputCursor TCursor,
	curr TCursor,
	hasMantissa *bool,
) bool {
	isDigit := currChar >= '0' && currChar <= '9'
	isMantissa := currChar == '.'

	if curr.CurrPos == inputCursor.CurrPos {
		if !isDigit || !isMantissa {
			return false
		}
		*hasMantissa = isMantissa
	}

	return true
}

func checkMantissa(
	currChar byte,
	inputCursor TCursor,
	curr TCursor,
	hasMantissa *bool,
) bool {
	isMantissa := currChar == '.'

	if isMantissa {
		if *hasMantissa {
			return false
		}
		*hasMantissa = isMantissa
	}

	return true
}

func checkExponential(
	source string,
	inputCursor TCursor,
	curr TCursor,
	hasMantissa *bool,
	hasExponent *bool,
) bool {
	isExponential := source[curr.CurrPos] == 'e'

	if isExponential {
		if *hasExponent || curr.CurrPos == uint(len(source)-1) {
			return false
		}

		*hasMantissa = true
		*hasExponent = true

		nextChar := source[curr.CurrPos+1]
		if nextChar == '-' || nextChar == '+' {
			curr.CurrPos++
			curr.Loc.Column++
		}
	}

	return true
}
