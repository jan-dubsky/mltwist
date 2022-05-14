package consoleui

import "strings"

const (
	// tabWidth is expected width of tab character in spaces.
	tabWidth = 8
	// width is expected screen with. After reaching width chars, the line
	// will be wrapped.
	width = 80
)

// findWordDelimSpace finds the latest word-delimiting space in the string and
// returns its index.
//
// This function will always return value in between 1 and len(s)-1, or -1 if no
// word-delimiting space is in the string.
//
// Value of 0 is never returned as a space at the beginning of a string is not a
// word-delimiting space.
func findWordDelimSpace(s string) int {
	for i := len(s) - 1; i > 0; i-- {
		if s[i] == ' ' {
			return i
		}
	}

	return -1
}

func format(s string, indent int, width int) string {
	chars := width - indent*tabWidth

	var sb strings.Builder
	for i := 0; i < len(s); sb.WriteByte('\n') {
		for j := 0; j < indent; j++ {
			sb.WriteByte('\t')
		}

		str := s[i:]
		if len(str) > chars {
			idx := findWordDelimSpace(str[:chars+1])
			if idx == -1 {
				idx = chars
			}
			str = str[:idx]
		}

		sb.WriteString(str)
		i += len(str)

		// Skip trailing spaces after line splitting.
		for ; i < len(s) && s[i] == ' '; i++ {
		}
	}

	return sb.String()
}
