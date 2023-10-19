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
			{Pattern: `\s+`, Type: chroma.Text, Mutator: nil},
			{Pattern: `(curl)((?:\s*\\?\s*))`, Type: chroma.Keyword, Mutator: nil},
			{Pattern: `(\-[A-Z])((?:\s*\\?\s*))`, Type: chroma.NameProperty, Mutator: nil},
			{Pattern: chroma.Words(``, ``, `POST`, `PUT`, `GET`, `DELETE`, `PATCH`, `OPTIONS`), Type: chroma.KeywordType, Mutator: nil},
			{Pattern: `(https|http)?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`, Type: chroma.String, Mutator: nil},
		},
	}
}
