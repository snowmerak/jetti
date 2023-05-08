package check

import (
	"strings"

	"github.com/snowmerak/jetti/v2/lib/model"
)

type Bean struct {
	StructName string
	Aliases    []string
}

func HasBean(path string, pkg *model.Package) ([]Bean, error) {
	beans := []Bean(nil)

	for _, st := range pkg.Structs {
		if strings.Contains(st.Doc, "jetti:bean") {
			bean := Bean{
				StructName: st.Name,
			}
			split := strings.Split(st.Doc, "\n")
			for _, line := range split {
				if strings.Contains(line, "jetti:bean") {
					bean.Aliases = append(bean.Aliases, strings.Split(strings.TrimSpace(strings.TrimPrefix(line, "jetti:bean")), " ")...)
					beans = append(beans, bean)
					break
				}
			}
		}
	}

	return beans, nil
}
