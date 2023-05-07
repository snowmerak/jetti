package check

import (
	"github.com/snowmerak/jetti/lib/model"
	"strings"
)

func HasBean(path string, pkg *model.Package) ([]string, error) {
	rs := make([]string, 0)

	for _, st := range pkg.Structs {
		if strings.Contains(st.Doc, "jetti:bean") {
			rs = append(rs, path+"/"+st.Name)
		}
	}

	return rs, nil
}
