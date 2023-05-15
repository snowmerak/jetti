package check

import (
	"strings"

	"github.com/snowmerak/jetti/v2/lib/model"
)

type Bean struct {
	Type    int
	Name    string
	Aliases []string
}

func HasBean(pkg *model.Package) ([]Bean, error) {
	beans := []Bean(nil)

	for _, st := range pkg.Structs {
		if strings.Contains(st.Doc, "jetti:bean") {
			bean := Bean{
				Type: TypeStruct,
				Name: st.Name,
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

	for _, it := range pkg.Interfaces {
		if strings.Contains(it.Doc, "jetti:bean") {
			bean := Bean{
				Type: TypeInterface,
				Name: it.Name,
			}
			split := strings.Split(it.Doc, "\n")
			for _, line := range split {
				if strings.Contains(line, "jetti:bean") {
					bean.Aliases = append(bean.Aliases, strings.Split(strings.TrimSpace(strings.TrimPrefix(line, "jetti:bean")), " ")...)
					beans = append(beans, bean)
					break
				}
			}
		}
	}

	for _, ali := range pkg.Aliases {
		if strings.Contains(ali.Doc, "jetti:bean") {
			bean := Bean{
				Type: TypeAlias,
				Name: ali.Name,
			}
			split := strings.Split(ali.Doc, "\n")
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
