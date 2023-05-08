package main

import (
	"fmt"
	"github.com/snowmerak/jetti/v2/lib/generator"
	"github.com/snowmerak/jetti/v2/lib/model"
)

func main() {
	pkg := model.Package{
		Name: "main",
		Imports: []model.Import{
			{
				Path: "fmt",
			},
		},
		Structs: []model.Struct{
			{
				Name: "Person",
				Fields: []model.Field{
					{
						Name: "Name",
						Type: "string",
					},
					{
						Name: "Age",
						Type: "int",
					},
				},
			},
		},
		Functions: []model.Function{
			{
				Name: "main",
				Code: []string{
					"p := Person{Name: \"John\", Age: 20}",
					"p.Name = \"Jack\"",
					"p.Age = 30",
					"fmt.Println(p)",
				},
			},
		},
	}

	rs, err := generator.GenerateFile(&pkg)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(rs))
}
