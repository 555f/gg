package file

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func guessAlias(path string) string {
	alias := strings.TrimSuffix(path, "/")
	if strings.Contains(alias, "/") {
		alias = alias[strings.LastIndex(alias, "/")+1:]
	}

	alias = strings.ToLower(alias)

	importsRegex := regexp.MustCompile(`[^a-z0-9]`)
	alias = importsRegex.ReplaceAllString(alias, "")

	for firstRune, runeLen := utf8.DecodeRuneInString(alias); unicode.IsDigit(firstRune); firstRune, runeLen = utf8.DecodeRuneInString(alias) {
		alias = alias[runeLen:]
	}

	if alias == "" {
		alias = "pkg"
	}

	return alias
}
