package main

import (
	"unicode"

	core "dappco.re/go/core"
)

func splitFields(s string) []string {
	return splitFunc(s, unicode.IsSpace)
}

func splitSnakeKebab(s string) []string {
	return splitFunc(s, func(r rune) bool {
		return r == '_' || r == '-'
	})
}

func splitFunc(s string, isDelimiter func(rune) bool) []string {
	var parts []string
	buf := core.NewBuilder()

	flush := func() {
		if buf.Len() == 0 {
			return
		}
		parts = append(parts, buf.String())
		buf.Reset()
	}

	for _, r := range s {
		if isDelimiter(r) {
			flush()
			continue
		}
		buf.WriteRune(r)
	}
	flush()

	return parts
}
