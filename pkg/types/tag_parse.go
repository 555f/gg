package types

import (
	"go/token"
	"regexp"
	"strconv"
	"strings"

	"github.com/555f/gg/pkg/errors"
)

var paramRegex = regexp.MustCompile("([^/]+)=([^/]+)")

type Tag struct {
	Key      string
	Value    string
	Options  []string
	Params   map[string]string
	Position token.Position
}

func (a *Tag) Param(name string) (val string, ok bool) {
	val, ok = a.Params[name]
	return
}

func (a *Tag) HasOption(opt string) bool {
	for _, option := range a.Options {
		if option == opt {
			return true
		}
	}
	return false
}

type Tags []*Tag

func (ts *Tags) GetSlice(key string) (tags []*Tag) {
	for _, tag := range *ts {
		if tag.Key == key {
			tags = append(tags, tag)
		}
	}
	return
}

func (ts *Tags) Get(key string) (*Tag, bool) {
	for _, tag := range *ts {
		if tag.Key == key {
			return tag, true
		}
	}
	return nil, false
}

func parseTags(comments Comments) (Tags, error) {
	var tags []*Tag
	for _, comment := range comments {
		val := strings.TrimSpace(comment.Value)
		for val != "" {
			i := 0
			for i < len(val) && val[i] == ' ' {
				i++
			}
			val = val[i:]
			if val == "" {
				break
			}
			if val[0] != '@' {
				break
			}
			val = val[1:]
			i = 0
			for i < len(val) && val[i] > ' ' && val[i] != ':' && val[i] != '"' && val[i] != 0x7f {
				i++
			}
			if i == 0 {
				return nil, errors.Error("bad syntax for tag key", comment.Position)
			}
			var (
				key     = val[:i]
				options []string
				value   string
				params  map[string]string
			)
			if i+1 >= len(val) || val[i] != ':' {
				val = val[i:]
			} else {
				if val[i+1] != '"' {
					return nil, errors.Error("bad syntax for tag value", comment.Position)
				}

				val = val[i+1:]

				i = 1
				for i < len(val) && val[i] != '"' {
					i++
				}
				if i >= len(val) {
					return nil, errors.Error("bad syntax for tag value", comment.Position)
				}

				lvalue := val[:i+1]
				val = val[i+1:]

				uqValue, err := strconv.Unquote(lvalue)
				if err != nil {
					return nil, errors.Error("bad syntax for tag value", comment.Position)
				}

				res := strings.Split(uqValue, ",")
				value = res[0]
				params = map[string]string{}

				for _, option := range res[1:] {
					if paramRegex.MatchString(option) {
						matches := paramRegex.FindStringSubmatch(option)
						if len(matches[0]) != len(option) {
							return nil, errors.Error("bad syntax for tag param", comment.Position)
						}
						params[matches[1]] = matches[2]
					} else {
						options = append(options, option)
					}
				}
				if len(options) == 0 {
					options = nil
				}
			}
			tags = append(tags, &Tag{
				Key:      key,
				Value:    value,
				Options:  options,
				Params:   params,
				Position: comment.Position,
			})
		}
	}
	return tags, nil
}
