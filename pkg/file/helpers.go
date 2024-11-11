package file

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func guessAlias(path string) string {
	alias := path

	if strings.HasSuffix(alias, "/") {
		// training slashes are usually tolerated, so we can get rid of one if
		// it exists
		alias = alias[:len(alias)-1]
	}

	if strings.Contains(alias, "/") {
		// if the path contains a "/", use the last part
		alias = alias[strings.LastIndex(alias, "/")+1:]
	}

	// alias should be lower case
	alias = strings.ToLower(alias)

	// alias should now only contain alphanumerics
	importsRegex := regexp.MustCompile(`[^a-z0-9]`)
	alias = importsRegex.ReplaceAllString(alias, "")

	// can't have a first digit, per Go identifier rules, so just skip them
	for firstRune, runeLen := utf8.DecodeRuneInString(alias); unicode.IsDigit(firstRune); firstRune, runeLen = utf8.DecodeRuneInString(alias) {
		alias = alias[runeLen:]
	}

	// If path part was all digits, we may be left with an empty string. In this case use "pkg" as the alias.
	if alias == "" {
		alias = "pkg"
	}

	return alias
}
