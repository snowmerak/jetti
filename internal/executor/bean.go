package executor

import (
	"github.com/snowmerak/jetti/internal/finder"
	"github.com/snowmerak/jetti/internal/generator"
	"github.com/snowmerak/jetti/internal/model"
	"go/format"
	"os"
	"path/filepath"
)

func Bean() {
	const direction = "bean"

	models := []model.Structs(nil)

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

		models = append(models, finder.FindStructName(f, direction))

		return nil
	}); err != nil {
		panic(err)
	}

	bs, err := generator.Generate(models...)
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
