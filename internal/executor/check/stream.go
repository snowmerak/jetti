package check

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"strings"
)

type Stream struct {
	PackageName       string
	Imports           []model.Import
	FunctionSignature []model.Function
}

func HasStream(pkg *model.Package) []Stream {
	streams := make([]Stream, 0)

	for _, s := range pkg.Structs {
		if strings.Contains(s.Doc, "jetti:stream") {
			for _, f := range s.Fields {
				if f.FuncType == nil {
					continue
				}

				streams = append(streams, Stream{
					PackageName: pkg.Name,
					Imports:     pkg.Imports,
				})

				streams[len(streams)-1].FunctionSignature = append(streams[len(streams)-1].FunctionSignature, model.Function{
					Name:   f.Name,
					Params: f.FuncType.Params,
					Return: f.FuncType.Return,
				})
			}
		}
	}

	return streams
}
