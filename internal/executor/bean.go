package executor

import (
	"bytes"
	"github.com/snowmerak/jetti/internal/finder"
	"github.com/snowmerak/jetti/internal/model"
	"go/format"
	"os"
	"path/filepath"
	"strings"
)

type BeanModel struct {
	PackagePath string
	Models      model.Structs
}

func Bean() {
	moduleName, err := finder.FindModuleName()
	if err != nil {
		panic(err)
	}

	const direction = "bean"

	models := []BeanModel(nil)

	if err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		found := finder.FindStructName(f, direction)

		if len(found.StructNames) == 0 {
			return nil
		}

		base := filepath.Dir(path)
		base = strings.ReplaceAll(base, "\\", "/")

		models = append(models, BeanModel{
			PackagePath: moduleName + "/" + base,
			Models:      found,
		})

		return nil
	}); err != nil {
		panic(err)
	}

	bs, err := GenerateBean(models...)
	if err != nil {
		panic(err)
	}

	rs, err := format.Source(bs)
	if err != nil {
		panic(err)
	}

	beanFolder := filepath.Join(generated, "bean")
	if err := os.MkdirAll(beanFolder, os.ModePerm); err != nil {
		panic(err)
	}

	f, err := os.Create(filepath.Join(beanFolder, "bean.go"))
	if err != nil {
		panic(err)
	}

	if _, err := f.Write(rs); err != nil {
		panic(err)
	}

	if err := f.Close(); err != nil {
		panic(err)
	}
}

func GenerateBean(models ...BeanModel) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	buffer.Write([]byte("package bean\n\n"))
	buffer.Write([]byte("import (\n"))
	for _, m := range models {
		buffer.Write([]byte("\t\"" + m.PackagePath + "\"\n"))
	}
	buffer.Write([]byte(")\n\n"))
	buffer.Write([]byte("type Bean struct {\n"))
	for _, m := range models {
		for _, structName := range m.Models.StructNames {
			buffer.Write([]byte("\t" + strings.ToLower(structName) + " *" + m.Models.PackageName + "." + structName + "\n"))
		}
	}
	buffer.Write([]byte("}\n\n"))

	buffer.Write([]byte("type Builder struct {\n"))
	buffer.Write([]byte("\tbean *Bean\n"))
	buffer.Write([]byte("}\n\n"))

	buffer.Write([]byte("func New() *Builder {\n"))
	buffer.Write([]byte("\treturn &Builder{}\n"))
	buffer.Write([]byte("}\n\n"))

	buffer.Write([]byte("func (b *Builder) Build() *Bean {\n"))
	buffer.Write([]byte("\treturn b.bean\n"))
	buffer.Write([]byte("}\n\n"))

	for _, m := range models {
		for _, structName := range m.Models.StructNames {
			buffer.Write([]byte("func (b *Builder) Add" + structName + "(" + strings.ToLower(structName) + " *" + m.Models.PackageName + "." + structName + ") *Builder {\n"))
			buffer.Write([]byte("\tb.bean." + strings.ToLower(structName) + " = " + strings.ToLower(structName) + "\n"))
			buffer.Write([]byte("\treturn b\n"))
			buffer.Write([]byte("}\n\n"))
		}
	}

	for _, m := range models {
		for _, structName := range m.Models.StructNames {
			buffer.Write([]byte("func (b *Bean) " + structName + "() *" + m.Models.PackageName + "." + structName + " {\n"))
			buffer.Write([]byte("\treturn b." + strings.ToLower(structName) + "\n"))
			buffer.Write([]byte("}\n\n"))
		}
	}

	return buffer.Bytes(), nil
}
