package control

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name          string
		str           string
		indent        int
		width         int
		expectedLines []string
	}{
		{
			name:          "single_line",
			str:           "abc abc abc abc",
			width:         16,
			expectedLines: []string{"abc abc abc abc"},
		},
		{
			name:          "single_indented_line",
			str:           "abc abc abc abc",
			indent:        2,
			width:         32,
			expectedLines: []string{"\t\tabc abc abc abc"},
		},
		{
			name:   "line_split",
			str:    "abc abc abc abc",
			indent: 1,
			width:  16,
			expectedLines: []string{
				"\tabc abc",
				"\tabc abc",
			},
		},
		{
			name:   "too_long_word",
			str:    "superlongword",
			indent: 1,
			width:  16,
			expectedLines: []string{
				"\tsuperlon",
				"\tgword",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			s := format(tt.str, tt.indent, tt.width)

			lines := strings.Split(s, "\n")
			r.Greater(len(lines), 1)
			r.Equal("", lines[len(lines)-1])

			lines = lines[:len(lines)-1]
			r.Equal(tt.expectedLines, lines)
		})
	}
}
