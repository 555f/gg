package types

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

type Module struct {
	ID       string
	Version  string
	Path     string
	Dir      string
	Indirect bool
	Main     bool
}

func (m *Module) ParseImportPath(s string) (pkgPath, name string, err error) {
	u, err := url.Parse("//" + s)
	if err != nil || u.Path == "" {
		return "", "", fmt.Errorf("invalid import path: %w", err)
	}
	name = path.Ext(u.Path)
	pkgPath = strings.Replace(u.Path, name, "", -1)
	if name == "" {
		return "", "", fmt.Errorf("invalid import path: %s, example ~/pkg/foo/ContextKey", s)
	}
	name = name[1:]
	if strings.HasPrefix(u.Host, "~") {
		pkgPath = m.Path + pkgPath
	} else {
		pkgPath = u.Host + pkgPath
	}
	return
}
