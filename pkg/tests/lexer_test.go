package main

import (
	"strings"
	"testing"

	"pkg/lexer"

	"github.com/stretchr/testify/assert"
)

func TestLexer_CheckReservedKeyword(testing *testing.T) {
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
		tok, _, ok := lexer.CheckReservedToken(test.value, lexer.TCursor{})
		assert.Equal(testing, test.keyword, ok, test.value)
		if ok {
			test.value = strings.TrimSpace(test.value)
			assert.Equal(testing, strings.ToLower(test.value), tok.Value, test.value)
		}
	}
}

func TestLexer_CheckNumeric(t *testing.T) {
	tests := []struct {
		number bool
		value  string
	}{
		{
			number: true,
			value:  "105",
		},
		{
			number: true,
			value:  "105 ",
		},
		{
			number: true,
			value:  "123.",
		},
		{
			number: true,
			value:  "123.145",
		},
		{
			number: true,
			value:  "1e5",
		},
		{
			number: true,
			value:  "1.e21",
		},
		{
			number: true,
			value:  "1.1e2",
		},
		{
			number: true,
			value:  "1.1e-2",
		},
		{
			number: true,
			value:  "1.1e+2",
		},
		{
			number: true,
			value:  "1e-1",
		},
		{
			number: true,
			value:  ".1",
		},
		{
			number: true,
			value:  "4.",
		},
		{
			number: false,
			value:  "e4",
		},
		{
			number: false,
			value:  "1..",
		},
		{
			number: false,
			value:  "1ee4",
		},
		{
			number: false,
			value:  " 1",
		},
	}

	for _, test := range tests {
		tok, _, ok := lexer.CheckNumeric(test.value, lexer.TCursor{})
		assert.Equal(t, test.number, ok, test.value)
		if ok {
			assert.Equal(t, strings.TrimSpace(test.value), tok.Value, test.value)
		}
	}
}

func TestLexer_CheckSymbol(t *testing.T) {
	tests := []struct {
		symbol bool
		value  string
	}{
		{
			symbol: true,
			value:  "* ",
		},
		{
			symbol: true,
			value:  ";",
		},
		{
			symbol: true,
			value:  "(",
		},
		{
			symbol: true,
			value:  ")",
		},
		{
			symbol: true,
			value:  " ",
		},
		{
			symbol: true,
			value:  "\n",
		},
		{
			symbol: true,
			value:  "\t",
		},
		{
			symbol: false,
			value:  "",
		},
		{
			symbol: false,
			value:  "=",
		},
	}

	for _, test := range tests {
		tok, _, ok := lexer.CheckSymbol(test.value, lexer.TCursor{})
		assert.Equal(t, test.symbol, ok, test.value)
		if ok {
			test.value = strings.TrimSpace(test.value)
			if len(test.value) > 0 {
				assert.Equal(t, test.value, tok.Value, test.value)
			}
		}
	}
}

func TestToken_CheckIdentifier(t *testing.T) {
	tests := []struct {
		identifier bool
		input      string
		value      string
	}{
		{
			identifier: true,
			input:      "a",
			value:      "a",
		},
		{
			identifier: true,
			input:      "abc",
			value:      "abc",
		},
		{
			identifier: true,
			input:      "abc ",
			value:      "abc",
		},
		{
			identifier: true,
			input:      `" abc "`,
			value:      ` abc `,
		},
		{
			identifier: true,
			input:      "a9$",
			value:      "a9$",
		},
		{
			identifier: true,
			input:      "userName",
			value:      "username",
		},
		{
			identifier: true,
			input:      `"userName"`,
			value:      "userName",
		},
		{
			identifier: false,
			input:      `"`,
		},
		{
			identifier: false,
			input:      "_sadsfa",
		},
		{
			identifier: false,
			input:      "9sadsfa",
		},
		{
			identifier: false,
			input:      " abc",
		},
	}

	for _, test := range tests {
		tok, _, ok := lexer.CheckIdentifier(test.input, lexer.TCursor{})
		assert.Equal(t, test.identifier, ok, test.input)
		if ok {
			assert.Equal(t, test.value, tok.Value, test.input)
		}
	}
}

func TestToken_CheckString(t *testing.T) {
	tests := []struct {
		string bool
		value  string
	}{
		{
			string: false,
			value:  "a",
		},
		{
			string: true,
			value:  "'abc'",
		},
		{
			string: true,
			value:  "'a b'",
		},
		{
			string: true,
			value:  "'a' ",
		},
		{
			string: true,
			value:  "'a '' b'",
		},
		{
			string: false,
			value:  "'",
		},
		{
			string: false,
			value:  "",
		},
		{
			string: false,
			value:  " 'foo'",
		},
	}

	for _, test := range tests {
		tok, _, ok := lexer.CheckString(test.value, lexer.TCursor{})
		assert.Equal(t, test.string, ok, test.value)
		if ok {
			test.value = strings.TrimSpace(test.value)
			assert.Equal(t, test.value[1:len(test.value)-1], tok.Value, test.value)
		}
	}
}
