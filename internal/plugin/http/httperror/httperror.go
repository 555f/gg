package httperror

import (
	"strings"

	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/gg"
)

type ErrorField struct {
	FldName string
	FldType any
	Name    string
}

type Error struct {
	Name        string
	Title       string
	Description string
	Fields      []ErrorField
	StatusCode  int64
}

func Load(structs []*gg.Struct, errorWrapper *options.ErrorWrapper) (result []Error) {
	if errorWrapper == nil {
		return
	}
	for _, st := range structs {
		var (
			fields     []ErrorField
			statusCode int64
		)
		for _, method := range st.Named.Methods {
			for _, field := range errorWrapper.Fields {
				if method.Text == "func StatusCode() int" {
					statusCode = method.Returns[0].(int64)
					continue
				}
				if strings.Replace(method.Text, "func ", "", -1) == field.Interface {
					//var fldType any
					//if len(method.Returns) == 1 {
					//	fldType = method.Returns[0]
					//}
					fields = append(fields, ErrorField{
						FldName: field.FldName,
						FldType: field.FldType,
						Name:    field.Name,
					})
					break
				}
			}
		}

		if len(fields) > 0 {
			result = append(result, Error{
				Name:        st.Named.Name,
				Title:       st.Named.Title,
				Description: st.Named.Description,
				StatusCode:  statusCode,
				Fields:      fields,
			})
		}
	}
	return
}
