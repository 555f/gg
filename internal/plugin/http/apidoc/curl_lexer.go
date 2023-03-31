package apidoc

import (
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
)

var _ = lexers.Register(chroma.MustNewLazyLexer(
	&chroma.Config{
		Name:            "Curl",
		Aliases:         []string{"curl"},
		Filenames:       []string{"*.curl"},
		CaseInsensitive: true,
	},
	curlRules,
))

func curlRules() chroma.Rules {
	return chroma.Rules{
		"root": {
			{`\s+`, chroma.Text, nil},
			{`(curl)((?:\s*\\?\s*))`, chroma.Keyword, nil},
			{`(\-[A-Z])((?:\s*\\?\s*))`, chroma.NameProperty, nil},
			{chroma.Words(``, ``, `POST`, `PUT`, `GET`, `DELETE`, `PATCH`, `OPTIONS`), chroma.KeywordType, nil},
			{`(https|http)?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`, chroma.String, nil},
		},
	}
}
