package executor

import (
	"github.com/snowmerak/jetti/internal/executor/generate"
	"os"
	"path/filepath"
)

func Generate(root string) error {
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		case ".go":
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
