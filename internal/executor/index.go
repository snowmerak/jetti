package executor

import (
	"github.com/snowmerak/jetti/v2/internal/cache"
	"github.com/snowmerak/jetti/v2/internal/executor/check"
	"github.com/snowmerak/jetti/v2/lib/model"
	"github.com/snowmerak/jetti/v2/lib/parser"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
)

const InterfacesKey = "interfaces"

func Index(root string) error {
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

	interfaces := make([]model.InterfaceTransferObject, 0)

	goRootPath := runtime.GOROOT()
	if err := filepath.Walk(filepath.Join(goRootPath, "src"), func(path string, info os.FileInfo, err error) error {
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

		if strings.Contains(path, "internal") {
			return nil
		}

		relativePath := strings.TrimPrefix(filepath.ToSlash(filepath.Dir(path)), filepath.ToSlash(goRootPath)+"/src/")

		switch filepath.Ext(path) {
		case ".go":
			pkg, err := parser.ParseFile(path)
			if err != nil {
				return nil
			}

			imports := make([]model.Import, 0)
			for _, imp := range pkg.Imports {
				imports = append(imports, imp)
			}

			for i := range pkg.Interfaces {
				if pkg.Interfaces[i].Name == "_" {
					continue
				}
				if len(pkg.Interfaces[i].Name) > 0 && unicode.IsLower(rune(pkg.Interfaces[i].Name[0])) {
					continue
				}
				interfaces = append(interfaces, model.InterfaceTransferObject{
					Path:      relativePath,
					Imports:   imports,
					Interface: pkg.Interfaces[i],
				})
			}
		}

		return nil
	}); err != nil {
		return err
	}

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

			switch filepath.Ext(path) {
			case ".go":
				pkg, err := parser.ParseFile(path)
				if err != nil {
					return err
				}

				imports := make([]model.Import, 0)
				for _, imp := range pkg.Imports {
					imports = append(imports, imp)
				}

				packagePath := strings.Join([]string{moduleName, filepath.ToSlash(filepath.Dir(relativePath))}, "/")
				for i := range pkg.Interfaces {
					interfaces = append(interfaces, model.InterfaceTransferObject{
						Path:      packagePath,
						Imports:   imports,
						Interface: pkg.Interfaces[i],
					})
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	{
		interfaceKeys := make([]string, len(interfaces))
		for i, ifce := range interfaces {
			interfaceKeys[i] = ifce.Path + " " + ifce.Interface.Name
		}

		if err := cc.SetInterfaceNames(InterfacesKey, interfaceKeys); err != nil {
			return err
		}

		for _, ifce := range interfaces {
			if err := cc.SetInterface(ifce.Path+" "+ifce.Interface.Name, ifce); err != nil {
				return err
			}
		}
	}

	return nil
}
