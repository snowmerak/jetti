package executor

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/snowmerak/jetti/v2/internal/cache"
	"github.com/snowmerak/jetti/v2/internal/executor/check"
	"github.com/snowmerak/jetti/v2/internal/executor/generate"
	"github.com/snowmerak/jetti/v2/lib/generator"
	"github.com/snowmerak/jetti/v2/lib/model"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ImplInteractive(root string) error {
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

	ifceKeys, ok := cc.GetInterfaceNames(InterfacesKey)
	if !ok {
		return errors.New("no interfaces found")
	}

	response := struct {
		Selected    []string `survey:"selected"`
		PackageName string   `survey:"package_name"`
		StructName  string   `survey:"struct_name"`
	}{}
	if err := survey.Ask([]*survey.Question{
		{
			Name: "selected",
			Prompt: &survey.MultiSelect{
				Message:  "Select interfaces:",
				Options:  ifceKeys,
				PageSize: 8,
			},
		},
		{
			Name: "package_name",
			Prompt: &survey.Input{
				Message: "Package path:",
			},
		},
		{
			Name: "struct_name",
			Prompt: &survey.Input{
				Message: "Struct name:",
			},
		},
	}, &response); err != nil {
		return err
	}

	packageName := filepath.Base(response.PackageName)
	packagePath := strings.Replace(response.PackageName, moduleName, root, 1)

	if err := os.MkdirAll(packagePath, os.ModePerm); err != nil {
		return err
	}

	methodMap := make(map[string]model.MethodTransferObject)

	for _, ifceKey := range response.Selected {
		ifce, ok := cc.GetInterface(ifceKey)
		if !ok {
			return errors.New("interface not found: " + ifceKey)
		}

		for _, method := range ifce.Interface.Methods {
			if _, ok := methodMap[method.Name]; ok {
				log.Printf("[skip] collision: %s.%s", ifce.Interface.Name, method.Name)
				continue
			}
			methodMap[method.Name] = model.MethodTransferObject{
				Path:    ifce.Path,
				Imports: ifce.Imports,
				Method:  method,
			}
		}
	}

	st := model.Struct{
		Name: response.StructName,
	}

	structFilePath := generate.MakeGeneratedFileName(packagePath, response.StructName)
	{
		structFilePkg := &model.Package{
			Name: packageName,
			Structs: []model.Struct{
				st,
			},
		}

		data, err := generator.GenerateFile(structFilePkg)
		if err != nil {
			return err
		}

		if err := os.WriteFile(structFilePath, data, os.ModePerm); err != nil {
			return err
		}
	}

	for _, method := range methodMap {
		methodFilePath := generate.MakeGeneratedFileName(packagePath, response.StructName, method.Method.Name)
		{
			method.Method.Receiver = model.Field{
				Name: strings.ToLower(response.StructName[:1]),
				Type: response.StructName,
			}
			method.Method.Code = []string{
				"// TODO: implement this method",
				"panic(\"not implemented\")",
			}
			methodFilePkg := &model.Package{
				Name:    packageName,
				Imports: method.Imports,
				Methods: []model.Method{
					method.Method,
				},
			}

			data, err := generator.GenerateFile(methodFilePkg)
			if err != nil {
				return err
			}

			if err := os.WriteFile(methodFilePath, data, os.ModePerm); err != nil {
				return err
			}
		}

		methodTestFilePath := generate.MakeGeneratedTestFileName(packagePath, response.StructName, method.Method.Name)
		testFilePkg := &model.Package{
			Name:    packageName,
			Imports: append(method.Imports, model.Import{Path: "testing"}),
			Functions: []model.Function{
				{
					Name: "Test" + response.StructName + "_" + method.Method.Name,
					Params: []model.Field{
						{
							Name: "t",
							Type: "*testing.T",
						},
					},
				},
				{
					Name: "Benchmark" + response.StructName + "_" + method.Method.Name,
					Params: []model.Field{
						{
							Name: "b",
							Type: "*testing.B",
						},
					},
				},
				{
					Name: "Example" + response.StructName + "_" + method.Method.Name,
				},
				{
					Name: "Fuzz" + response.StructName + "_" + method.Method.Name,
					Params: []model.Field{
						{
							Name: "f",
							Type: "*testing.F",
						},
					},
				},
			},
		}

		data, err := generator.GenerateFile(testFilePkg)
		if err != nil {
			return err
		}

		if err := os.WriteFile(methodTestFilePath, data, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func ImplTargets(root string, target []string) error {
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

	_ = cc

	return nil
}
