package control

import "strings"

const (
	// tabWidth is expected width of tab character in spaces.
	tabWidth = 8
	// width is expected screen with. After reaching width chars, the line
	// will be wrapped.
	width = 80
)

func format(s string, indent int, width int) string {
	chars := width - indent*tabWidth

	var sb strings.Builder
	for i := 0; i < len(s); i += chars {
		if i != 0 {
			sb.WriteByte('\n')
		}

		for j := 0; j < indent; j++ {
			sb.WriteByte('\t')
		}

		str := s[i*chars:]
		if len(str) > chars {
			str = str[:chars]
		}
		sb.WriteString(str)
	}

	return sb.String()
}
