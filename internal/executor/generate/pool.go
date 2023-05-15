package generate

import (
	"github.com/snowmerak/jetti/v2/internal/executor/check"
	"github.com/snowmerak/jetti/v2/lib/generator"
	"github.com/snowmerak/jetti/v2/lib/model"
	"os"
	"path/filepath"
	"strings"
)

func Pool(path string, pools []check.Pool) error {
	dir := filepath.Dir(path)
	packageName := filepath.Base(dir)

	for _, pool := range pools {
		alias := strings.ToUpper(pool.Alias[:1]) + pool.Alias[1:]
		lowerAlias := strings.ToLower(alias)
		filePath := filepath.Join(dir, lowerAlias+".pool.go")

		typ := pool.TypeName
		switch pool.Type {
		case check.TypeStruct:
			fallthrough
		case check.TypeAlias:
			typ = "*" + pool.TypeName
		}
		pkg := (*model.Package)(nil)
		switch pool.PoolKind {
		case check.SyncPool:
			pkg = &model.Package{
				Name: packageName,
				Imports: []model.Import{
					{
						Path: "sync",
					},
					{
						Path: "errors",
					},
					{
						Path: "runtime",
					},
				},
				GlobalVariables: []model.GlobalVariable{
					{
						Name:  "err" + alias + "CannotGet",
						Type:  "error",
						Value: "errors.New(\"cannot get " + lowerAlias + "\")",
					},
				},
				Structs: []model.Struct{
					{
						Name: alias + "Pool",
						Fields: []model.Field{
							{
								Name: "pool",
								Type: "*sync.Pool",
							},
						},
						Methods: []model.Method{
							{
								Name: "Get",
								Return: []model.Field{
									{
										Type: typ,
									},
									{
										Type: "error",
									},
								},
								Code: []string{
									"v := $RECEIVER$.pool.Get()",
									"if v == nil {",
									"\treturn nil, err" + alias + "CannotGet",
									"}",
									"return v.(" + typ + "), nil",
								},
							},
							{
								Name: "GetWithFinalizer",
								Return: []model.Field{
									{
										Type: typ,
									},
									{
										Type: "error",
									},
								},
								Code: []string{
									"v := $RECEIVER$.pool.Get()",
									"if v == nil {",
									"\treturn nil, err" + alias + "CannotGet",
									"}",
									"runtime.SetFinalizer(v, func(v interface{}) {",
									"\t$RECEIVER$.pool.Put(v)",
									"})",
									"return v.(" + typ + "), nil",
								},
							},
							{
								Name: "Put",
								Params: []model.Field{
									{
										Name: "v",
										Type: typ,
									},
								},
								Code: []string{
									"$RECEIVER$.pool.Put(v)",
								},
							},
						},
					},
				},
				Functions: []model.Function{
					{
						Name: "New" + alias + "Pool",
						Return: []model.Field{
							{
								Type: alias + "Pool",
							},
						},
						Code: []string{
							"return " + alias + "Pool{",
							"\tpool: &sync.Pool{",
							"\t\tNew: func() interface{} {",
							"\t\t\treturn new(" + pool.TypeName + ")",
							"\t\t},",
							"\t},",
							"}",
						},
					},
					{
						Name: "Is" + alias + "CannotGetErr",
						Params: []model.Field{
							{
								Name: "err",
								Type: "error",
							},
						},
						Return: []model.Field{
							{
								Type: "bool",
							},
						},
						Code: []string{
							"return errors.Is(err, err" + alias + "CannotGet)",
						},
					},
				},
			}
		case check.ChannelPool:
			pkg = &model.Package{
				Name: packageName,
				Imports: []model.Import{
					{
						Path: "runtime",
					},
					{
						Path: "time",
					},
				},
				Structs: []model.Struct{
					{
						Name: alias + "Pool",
						Fields: []model.Field{
							{
								Name: "pool",
								Type: "chan " + typ,
							},
							{
								Name: "timeout",
								Type: "time.Duration",
							},
						},
						Methods: []model.Method{
							{
								Name: "Get",
								Return: []model.Field{
									{
										Type: typ,
									},
								},
								Code: []string{
									"after := time.After($RECEIVER$.timeout)",
									"select {",
									"case v := <-$RECEIVER$.pool:",
									"\treturn v",
									"case <-after:",
									"\treturn new(" + pool.TypeName + ")",
									"}",
								},
							},
							{
								Name: "GetWithFinalizer",
								Return: []model.Field{
									{
										Type: typ,
									},
								},
								Code: []string{
									"after := time.After($RECEIVER$.timeout)",
									"resp := (" + typ + ")(nil)",
									"select {",
									"case v := <-$RECEIVER$.pool:",
									"\tresp = v",
									"case <-after:",
									"\tresp = new(" + pool.TypeName + ")",
									"}",
									"runtime.SetFinalizer(resp, func(v interface{}) {",
									"\t$RECEIVER$.pool <- v.(" + typ + ")",
									"})",
									"return resp",
								},
							},
							{
								Name: "Put",
								Params: []model.Field{
									{
										Name: "v",
										Type: typ,
									},
								},
								Code: []string{
									"select {",
									"case $RECEIVER$.pool <- v:",
									"default:",
									"}",
								},
							},
						},
					},
				},
				Functions: []model.Function{
					{
						Name: "New" + alias + "Pool",
						Params: []model.Field{
							{
								Name: "size",
								Type: "int",
							},
							{
								Name: "timeout",
								Type: "time.Duration",
							},
						},
						Return: []model.Field{
							{
								Type: alias + "Pool",
							},
						},
						Code: []string{
							"pool := make(chan " + typ + ", size)",
							"return " + alias + "Pool{",
							"\tpool: pool,",
							"\ttimeout: timeout,",
							"}",
						},
					},
				},
			}
		default:
			continue
		}

		data, err := generator.GenerateFile(pkg)
		if err != nil {
			return err
		}

		f, err := os.Create(filePath)
		if err != nil {
			return err
		}

		if _, err := f.Write(data); err != nil {
			return err
		}
	}

	return nil
}
