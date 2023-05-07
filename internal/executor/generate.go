package executor

import (
	"github.com/snowmerak/jetti/internal/executor/check"
	"github.com/snowmerak/jetti/internal/executor/generate"
	"github.com/snowmerak/jetti/lib/parser"
	"os"
	"path/filepath"
)

func Generate(root string) error {
	moduleName, err := check.GetModuleName(root)
	if err != nil {
		return err
	}

	beanStructs := make([]string, 0)

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		case ".go":
			pkg, err := parser.ParseFile(path)
			if err != nil {
				return err
			}

			beans, err := check.HasBean(pkg)
			if err != nil {
				return err
			}

			for _, bean := range beans {
				beanStructs = append(beanStructs, path+"/"+bean.Name)
			}
		case ".json":
			if err := generate.ConvertJson(path); err != nil {
				return err
			}
		case ".yml":
			fallthrough
		case ".yaml":
			if err := generate.ConvertYaml(path); err != nil {
				return err
			}
		case ".proto":
			if err := generate.BuildProtobuf(root, path); err != nil {
				return err
			}
		case ".fbs":
			if err := generate.BuildFlatbuffers(root, path); err != nil {
				return err
			}
		default:
			return nil
		}

		return nil
	}); err != nil {
		return err
	}

	if err := generate.MakeContextPackage(root, moduleName, beanStructs...); err != nil {
		return err
	}

	return nil
}
