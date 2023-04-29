package executor

import (
	"bytes"
	"fmt"
	"github.com/snowmerak/jetti/internal/strcase"
	"github.com/twpayne/go-jsonstruct/v2"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func ModelNew(path string) {
	content := ""

	switch filepath.Ext(path) {
	case ".json":
		content = "{}"
	case ".yaml":
		content = "a: b"
	default:
		panic("unsupported file type, only support json/yaml")
	}

	if err := os.MkdirAll("template/model/"+filepath.Dir(path), os.ModePerm); err != nil {
		panic(err)
	}

	targetFilePath := "template/model/" + path
	if err := os.WriteFile(targetFilePath, []byte(content), os.ModePerm); err != nil {
		panic(err)
	}
}

func ModelJson(path string) {
	filename := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	typeName := strcase.SnakeToPascal(filename)
	packageName := filepath.Base(filepath.Dir(path))

	if runtime.GOOS == "windows" {
		path = strings.ReplaceAll(path, "\\", "/")
	}

	targetFilePath := "template/model/" + path

	generator := jsonstruct.NewGenerator(
		jsonstruct.WithIntType("int64"),
		jsonstruct.WithPackageName(packageName),
		jsonstruct.WithTypeName(typeName),
		jsonstruct.WithOmitEmpty(jsonstruct.OmitEmptyAuto),
		jsonstruct.WithImports("os", "github.com/goccy/go-json"),
		jsonstruct.WithGoFormat(true))

	if err := generator.ObserveJSONFile(targetFilePath); err != nil {
		panic(err)
	}

	code, err := generator.Generate()
	if err != nil {
		panic(err)
	}

	buffer := bytes.NewBuffer(code)
	buffer.WriteString("\n")
	buffer.WriteString(fmt.Sprintf("func (m *%s) ToJson() ([]byte, error) {\n", typeName))
	buffer.WriteString("\treturn json.Marshal(m)\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString(fmt.Sprintf("func FromBytes(data []byte) (*%s, error) {\n", typeName))
	buffer.WriteString(fmt.Sprintf("\tm := new(%s)\n", typeName))
	buffer.WriteString("\terr := json.Unmarshal(data, m)\n")
	buffer.WriteString("\treturn m, err\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString(fmt.Sprintf("func FromFile(path string) (*%s, error) {\n", typeName))
	buffer.WriteString("\tdata, err := os.ReadFile(path)\n")
	buffer.WriteString("\tif err != nil {\n")
	buffer.WriteString("\t\treturn nil, err\n")
	buffer.WriteString("\t}\n")
	buffer.WriteString(fmt.Sprintf("\treturn FromBytes(data)\n"))
	buffer.WriteString("}\n\n")

	if err := os.MkdirAll("gen/model/"+filepath.Dir(path), os.ModePerm); err != nil {
		panic(err)
	}

	saveFilePath := "gen/model/" + path + ".go"
	if err := os.WriteFile(saveFilePath, code, os.ModePerm); err != nil {
		panic(err)
	}

	if err := goGet("github.com/goccy/go-json"); err != nil {
		panic(err)
	}
}

func ModelYaml(path string) {
	filename := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	typeName := strcase.SnakeToPascal(filename)
	packageName := filepath.Base(filepath.Dir(path))

	if runtime.GOOS == "windows" {
		path = strings.ReplaceAll(path, "\\", "/")
	}

	targetFilePath := "template/model/" + path

	generator := jsonstruct.NewGenerator(
		jsonstruct.WithIntType("int64"),
		jsonstruct.WithPackageName(packageName),
		jsonstruct.WithTypeName(typeName),
		jsonstruct.WithOmitEmpty(jsonstruct.OmitEmptyAuto),
		jsonstruct.WithImports("os", "github.com/goccy/go-yaml"),
		jsonstruct.WithGoFormat(true))

	if err := generator.ObserveYAMLFile(targetFilePath); err != nil {
		panic(err)
	}

	code, err := generator.Generate()
	if err != nil {
		panic(err)
	}

	buffer := bytes.NewBuffer(code)
	buffer.WriteString("\n")
	buffer.WriteString("func (y *" + typeName + ") ToYaml() ([]byte, error) {\n")
	buffer.WriteString("\treturn yaml.Marshal(y)\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString(fmt.Sprintf("func FromBytes(data []byte) (*%s, error) {\n", typeName))
	buffer.WriteString(fmt.Sprintf("\tr := new(%s)\n", typeName))
	buffer.WriteString("\terr := yaml.Unmarshal(data, r)\n")
	buffer.WriteString("\tif err != nil {\n")
	buffer.WriteString("\t\treturn nil, err\n")
	buffer.WriteString("\t}\n")
	buffer.WriteString("\treturn r, nil\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString(fmt.Sprintf("func FromFile(path string) (*%s, error) {\n", typeName))
	buffer.WriteString("\tdata, err := os.ReadFile(path)\n")
	buffer.WriteString("\tif err != nil {\n")
	buffer.WriteString("\t\treturn nil, err\n")
	buffer.WriteString("\t}\n")
	buffer.WriteString(fmt.Sprintf("\treturn FromBytes(data)\n"))
	buffer.WriteString("}\n\n")

	if err := os.MkdirAll("gen/model/"+filepath.Dir(path), os.ModePerm); err != nil {
		panic(err)
	}

	saveFilePath := "gen/model/" + path + ".go"
	if err := os.WriteFile(saveFilePath, buffer.Bytes(), os.ModePerm); err != nil {
		panic(err)
	}

	if err := goGet("github.com/goccy/go-yaml"); err != nil {
		panic(err)
	}
}
