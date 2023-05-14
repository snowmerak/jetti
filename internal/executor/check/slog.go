package check

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"strings"
)

type Slog struct {
	Name   string
	Fields []model.Field
}

func HasSlog(pkg *model.Package) ([]Slog, error) {
	ss := []Slog(nil)

	for _, st := range pkg.Structs {
		if strings.Contains(st.Doc, "jetti:slog") {
			ss = append(ss, Slog{
				Name:   st.Name,
				Fields: st.Fields,
			})
		}
	}

	return ss, nil
}
