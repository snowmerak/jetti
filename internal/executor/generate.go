package executor

import (
	"github.com/snowmerak/jetti/v2/internal/cache"
	"github.com/snowmerak/jetti/v2/internal/executor/check"
	"github.com/snowmerak/jetti/v2/internal/executor/generate"
	"github.com/snowmerak/jetti/v2/lib/parser"
	"log"
	"os"
	"path/filepath"
)

func Generate(root string) error {
	moduleName, err := check.GetModuleName(root)
	if err != nil {
		return err
	}
	_ = moduleName

	cachePath := filepath.Join(root, ".jetti-cache")
	if err := os.MkdirAll(cachePath, os.ModePerm); err != nil {
		return err
	}

	cc, err := cache.NewCache(cachePath)
	if err != nil {
		return err
	}

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		relativePath := filepath.ToSlash(path)[len(filepath.ToSlash(root))+1:]
		prevGenTime, ok := cc.Get(relativePath)
		if ok && prevGenTime >= info.ModTime().Unix() {
			log.Printf("skip: %s", relativePath)
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
			log.Printf("generate bean: %s", relativePath)
		case ".json":
			if err := generate.ConvertJson(path); err != nil {
				return err
			}
			log.Printf("generate json: %s", relativePath)
		case ".yml":
			fallthrough
		case ".yaml":
			if err := generate.ConvertYaml(path); err != nil {
				return err
			}
			log.Printf("generate yaml: %s", relativePath)
		case ".proto":
			if err := generate.BuildProtobuf(root, path); err != nil {
				return err
			}
			log.Printf("generate protobuf/grpc: %s", relativePath)
		case ".fbs":
			if err := generate.BuildFlatbuffers(root, path); err != nil {
				return err
			}
			log.Printf("generate flatbuffers: %s", relativePath)
		default:
			return nil
		}

		if err := cc.Set(relativePath, info.ModTime().Unix()); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
