package executor

import (
	"fmt"
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

	fmt.Printf("%+v\n", models)

	bs, err := generator.Generate(models...)
	if err != nil {
		panic(err)
	}

	rs, err := format.Source(bs)
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll("bean", os.ModePerm); err != nil {
		panic(err)
	}

	f, err := os.Create("src/bean/bean.go")
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
