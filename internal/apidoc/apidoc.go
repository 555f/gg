package apidoc

import (
	_ "embed"
	"encoding/json"
	"text/template"

	"github.com/555f/gg/pkg/file"
)

//go:embed files/index.go.tpl
var apiDocTemplate string

//go:embed files/style.css.tpl
var style string

//go:embed files/script.js.tpl
var script string

type Param struct {
	Name     string  `json:"name"`
	Title    string  `json:"title"`
	Type     string  `json:"type"`
	In       string  `json:"in"`
	Required bool    `json:"required"`
	Array    bool    `json:"array"`
	Example  string  `json:"example"`
	Params   []Param `json:"params"`
}

type Endpoint struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Path        string  `json:"path"`
	Method      string  `json:"method"`
	ContentType string  `json:"contentType"`
	Params      []Param `json:"params"`
	Results     []Param `json:"results"`
}

type Service struct {
	Name        string     `json:"name"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Endpoints   []Endpoint `json:"endpoints"`
}

type API struct {
	Title    string             `json:"title"`
	Schemas  map[string]*Schema `json:"schemas"`
	Services []Service          `json:"services"`
}

type templateValues struct {
	Style   string
	Script  string
	APIData string
}

func Gen(title string, api API) func(f *file.TxtFile) error {
	return func(f *file.TxtFile) error {
		apiData, _ := json.Marshal(api)
		t := template.Must(template.New("apidoc").Parse(apiDocTemplate))
		return t.Execute(f, &templateValues{
			APIData: string(apiData),
			Style:   style,
			Script:  script,
		})
	}
}
