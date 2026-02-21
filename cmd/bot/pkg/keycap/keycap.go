package keycap

import "strings"

func Convert(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
			b.WriteRune(r)
			b.WriteRune('\uFE0F')
			b.WriteRune('\u20E3')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
