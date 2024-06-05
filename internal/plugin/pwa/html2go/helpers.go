package html2go

import (
	"strings"
)

func trimQuote(s string) string {
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func extractKeyName(attrName, prefix string) (key string) {
	if i := strings.Index(attrName, prefix); i != -1 {
		key = attrName[i+len(prefix):]
		return
	}
	return
}
