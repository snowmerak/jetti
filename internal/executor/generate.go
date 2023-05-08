package executor

import (
	"github.com/snowmerak/jetti/v2/internal/executor/check"
	"github.com/snowmerak/jetti/v2/internal/executor/generate"
	"github.com/snowmerak/jetti/v2/lib/parser"
	"os"
	"path/filepath"
)

func Generate(root string) error {
	moduleName, err := check.GetModuleName(root)
	if err != nil {
		return err
	}
	_ = moduleName

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

			beans, err := check.HasBean(path, pkg)
			if err != nil {
				return err
			}

			if err := generate.Bean(path, beans); err != nil {
				return err
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

	return nil
}
