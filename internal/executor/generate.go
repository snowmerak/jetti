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

var subDirectories = []string{"lib", "internal", "model"}

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

	beanUpdated := false
	jsonGenerated := false
	yamlGenerated := false

	for _, subDir := range subDirectories {
		if err := filepath.Walk(filepath.Join(root, subDir), func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return nil
			}

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

			errFaces := make([]check.Getter, 0)

			log.Println("check:", relativePath)
			switch filepath.Ext(path) {
			case ".go":
				pkg, err := parser.ParseFile(path)
				if err != nil {
					return err
				}

				requests, err := check.HasBean(pkg, generate.RequestDirective)
				if err != nil {
					return err
				}

				if len(requests) > 0 {
					if err := generate.RequestScopeData(path, requests); err != nil {
						return err
					}
					log.Printf("generate bean: %s", relativePath)
				}

				beans, err := check.HasBean(pkg, generate.BeanDirective)
				if err != nil {
					return err
				}

				if len(beans) > 0 {
					beanUpdated = true
					if err := generate.Bean(moduleName, path, beans); err != nil {
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

				getter, err := check.HasGetter(pkg)
				if err != nil {
					return err
				}

				if len(getter.StructMap) > 0 {
					errFaces = append(errFaces, getter)
					if err := generate.Getter(path, getter); err != nil {
						return err
					}
					log.Printf("generate getter: %s", relativePath)
				}
			case ".json":
				if err := generate.ConvertJson(path); err != nil {
					return err
				}
				jsonGenerated = true
				log.Printf("generate json: %s", relativePath)
			case ".yml":
				fallthrough
			case ".yaml":
				if err := generate.ConvertYaml(path); err != nil {
					return err
				}
				yamlGenerated = true
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

			if len(errFaces) > 0 {
				if err := generate.ErrFace(root, errFaces); err != nil {
					return err
				}
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

	if beanUpdated {
		if err := generate.BeanContainer(root); err != nil {
			return err
		}
		beanUpdated = false
	}

	if jsonGenerated {
		if err := generate.JsonSerializable(root); err != nil {
			return err
		}
		jsonGenerated = false
	}

	if yamlGenerated {
		if err := generate.YamlSerializable(root); err != nil {
			return err
		}
		yamlGenerated = false
	}

	return nil
}
