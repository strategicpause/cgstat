package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type ShortenTest struct {
	str           string
	truncationLen int
	expected      string
	description   string
}

func TestShorten(t *testing.T) {
	tests := []*ShortenTest{
		{str: "aba", truncationLen: 5, expected: "aba",
			description: "Return the string if it is less than the truncation length."},
		{str: "aba", truncationLen: 3, expected: "aba",
			description: "Return the string if it is equal to the truncation length."},
		{str: "aba", truncationLen: 1, expected: DividerText,
			description: "Return the divider text if it is greater than the length."},
		{str: "ababa", truncationLen: 3, expected: DividerText,
			description: "Return the divider text if it is greater than the length."},
		{str: "abababab", truncationLen: 6, expected: "a[...]b",
			description: "Return truncated text."},
		{str: "abababa", truncationLen: 6, expected: "a[...]a",
			description: "Return truncated text."},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			shortenStr := Shorten(test.str, test.truncationLen)
			assert.Equal(t, test.expected, shortenStr)
		})
	}
}
