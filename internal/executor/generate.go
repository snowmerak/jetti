package executor

import (
	"github.com/snowmerak/jetti/v2/internal/cache"
	"github.com/snowmerak/jetti/v2/internal/executor/check"
	"github.com/snowmerak/jetti/v2/internal/executor/generate"
	"github.com/snowmerak/jetti/v2/lib/parser"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var subDirectories = []string{"lib", "cmd", "internal", "model"}

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

	for _, subDir := range subDirectories {
		if err := filepath.Walk(filepath.Join(root, subDir), func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if strings.HasPrefix(info.Name(), ".") {
				return nil
			}

			sp := strings.Split(filepath.ToSlash(path), "/")
			for _, s := range sp {
				if strings.HasPrefix(s, ".") {
					return nil
				}
			}

			relativePath := filepath.ToSlash(path)[len(filepath.ToSlash(root))+1:]
			prevGenTime, ok := cc.Get(relativePath)
			if ok && prevGenTime >= info.ModTime().Unix() {
				log.Printf("skip: %s", relativePath)
				return nil
			}

			log.Println("check:", relativePath)
			switch filepath.Ext(path) {
			case ".go":
				pkg, err := parser.ParseFile(path)
				if err != nil {
					return err
				}

				requests, err := check.HasBean(pkg, generate.Request)
				if err != nil {
					return err
				}

				if len(requests) > 0 {
					if err := generate.RequestScopeData(path, requests); err != nil {
						return err
					}
					log.Printf("generate bean: %s", relativePath)
				}

				parameters, err := check.HasOptionalParameter(pkg)
				if err != nil {
					return err
				}

				if len(parameters) > 0 {
					if err := generate.OptionalParameter(path, parameters); err != nil {
						return err
					}
					log.Printf("generate parameter: %s", relativePath)
				}

				pools, err := check.HasPool(pkg)
				if err != nil {
					return err
				}

				if len(pools) > 0 {
					if err := generate.Pool(path, pools); err != nil {
						return err
					}
					log.Printf("generate pool: %s", relativePath)
				}

				optionals, err := check.HasOptional(pkg)
				if err != nil {
					return err
				}

				if len(optionals) > 0 {
					if err := generate.Option(path, optionals); err != nil {
						return err
					}
					log.Printf("generate optional: %s", relativePath)
				}
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

			log.Println("update:", relativePath)
			if err := cc.Set(relativePath, info.ModTime().Unix()); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}
