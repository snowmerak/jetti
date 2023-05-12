package check

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"strings"
)

type Pool struct {
	Type     int
	Alias    string
	TypeName string
	PoolKind int
}

func HasPool(pkg *model.Package) ([]Pool, error) {
	ps := []Pool(nil)

	for _, st := range pkg.Structs {
		if strings.Contains(st.Doc, "jetti:pool") {
			split := strings.Split(st.Doc, "\n")
			for _, s := range split {
				if strings.Contains(s, "jetti:pool") {
					list := strings.Split(strings.TrimPrefix(s, "jetti:pool"), " ")
					for _, p := range list {
						if p == "" {
							continue
						}
						sp := strings.Split(p, ":")
						if len(sp) < 2 {
							continue
						}
						pk := SyncPool
						if sp[0] == "chan" {
							pk = ChannelPool
						}
						ps = append(ps, Pool{
							Type:     TypeStruct,
							Alias:    sp[1],
							TypeName: st.Name,
							PoolKind: pk,
						})
					}
				}
			}
		}
	}

	for _, it := range pkg.Interfaces {
		if strings.Contains(it.Doc, "jetti:pool") {
			split := strings.Split(it.Doc, "\n")
			for _, s := range split {
				if strings.Contains(s, "jetti:pool") {
					list := strings.Split(strings.TrimPrefix(s, "jetti:pool"), " ")
					for _, p := range list {
						if p == "" {
							continue
						}
						sp := strings.Split(p, ":")
						if len(sp) < 2 {
							continue
						}
						pk := SyncPool
						if sp[0] == "chan" {
							pk = ChannelPool
						}
						ps = append(ps, Pool{
							Type:     TypeInterface,
							Alias:    sp[1],
							TypeName: it.Name,
							PoolKind: pk,
						})
					}
				}
			}
		}
	}

	return ps, nil
}
