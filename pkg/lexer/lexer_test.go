package tugle

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken_lexKeyword(testing *testing.T) {
	tests := []struct {
		keyword bool
		value   string
	}{
		{
			keyword: true,
			value:   "select ",
		},
		{
			keyword: true,
			value:   "from",
		},
		{
			keyword: true,
			value:   "as",
		},
		{
			keyword: true,
			value:   "SELECT",
		},
		{
			keyword: true,
			value:   "into",
		},
		{
			keyword: false,
			value:   " into",
		},
		{
			keyword: false,
			value:   "flubbrety",
		},
	}

	for _, test := range tests {
		tok, _, ok := checkReservedToken(test.value, TCursor{})
		assert.Equal(testing, test.keyword, ok, test.value)
		if ok {
			test.value = strings.TrimSpace(test.value)
			assert.Equal(testing, strings.ToLower(test.value), tok.Value, test.value)
		}
	}
}
