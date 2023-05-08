package generate

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/snowmerak/jetti/lib/generator"
	"github.com/snowmerak/jetti/lib/model"
)

const CONTEXT = `context`

type elementInfos struct {
	path       string
	callName   string
	memberName string
	structName string
}

func MakeContextPackage(root string, module string, elements ...string) error {
	genPath := filepath.Join(root, "gen", CONTEXT)

	if err := os.MkdirAll(genPath, os.ModePerm); err != nil {
		return err
	}

	infos := make([]elementInfos, 0, len(elements))
	for _, element := range elements {
		split := strings.Split(filepath.ToSlash(element), "/")
		if len(split) < 2 {
			continue
		}
		name := split[len(split)-2] + "." + split[len(split)-1]
		path := module + strings.TrimPrefix(element, root)
		infos = append(infos, elementInfos{
			path:       filepath.Dir(path),
			callName:   name,
			memberName: strings.ToLower(split[len(split)-1][:1]) + split[len(split)-1][1:],
			structName: split[len(split)-1],
		})
	}

	imports := make([]model.Import, 2, len(infos)+2)
	imports[0] = model.Import{
		Path: "context",
	}
	imports[1] = model.Import{
		Path: "time",
	}
	for _, info := range infos {
		imports = append(imports, model.Import{
			Path: info.path,
		})
	}

	fields := make([]model.Field, 1, len(elements)+1)
	fields[0] = model.Field{
		Name: "ctx",
		Type: "context.Context",
	}
	for _, info := range infos {
		fields = append(fields, model.Field{
			Name: info.memberName,
			Type: "*" + info.callName,
		})
	}

	methods := make([]model.Method, 5, len(elements)+5)
	methods[0] = model.Method{
		Name: "WithCancel",
		Return: []model.Field{
			{
				Type: "*Context",
			},
			{
				Type: "context.CancelFunc",
			},
		},
		Code: []string{
			"ctx, cancel := context.WithCancel($RECEIVER$.ctx)",
			"newCtx := new(Context)",
			"*newCtx = *$RECEIVER$",
			"newCtx.ctx = ctx",
			"return newCtx, cancel",
		},
	}
	methods[1] = model.Method{
		Name: "WithDeadline",
		Params: []model.Field{
			{
				Name: "deadline",
				Type: "time.Time",
			},
		},
		Return: []model.Field{
			{
				Type: "*Context",
			},
			{
				Type: "context.CancelFunc",
			},
		},
		Code: []string{
			"ctx, cancel := context.WithDeadline($RECEIVER$.ctx, deadline)",
			"newCtx := new(Context)",
			"*newCtx = *$RECEIVER$",
			"newCtx.ctx = ctx",
			"return newCtx, cancel",
		},
	}
	methods[2] = model.Method{
		Name: "WithTimeout",
		Params: []model.Field{
			{
				Name: "timeout",
				Type: "time.Duration",
			},
		},
		Return: []model.Field{
			{
				Type: "*Context",
			},
			{
				Type: "context.CancelFunc",
			},
		},
		Code: []string{
			"ctx, cancel := context.WithTimeout($RECEIVER$.ctx, timeout)",
			"newCtx := new(Context)",
			"*newCtx = *$RECEIVER$",
			"newCtx.ctx = ctx",
			"return newCtx, cancel",
		},
	}
	methods[3] = model.Method{
		Name: "Done",
		Return: []model.Field{
			{
				Type: "<-chan struct{}",
			},
		},
		Code: []string{
			"return $RECEIVER$.ctx.Done()",
		},
	}
	methods[4] = model.Method{
		Name: "RawContext",
		Return: []model.Field{
			{
				Type: "context.Context",
			},
		},
		Code: []string{
			"return $RECEIVER$.ctx",
		},
	}
	for _, info := range infos {
		methods = append(methods, model.Method{
			Name: "Get" + info.structName,
			Return: []model.Field{
				{
					Type: "*" + info.callName,
				},
			},
			Code: []string{
				"return $RECEIVER$." + info.memberName,
			},
		})
	}

	newFuncParams := make([]model.Field, 0, len(infos))
	for _, info := range infos {
		newFuncParams = append(newFuncParams, model.Field{
			Name: info.memberName,
			Type: "*" + info.callName,
		})
	}
	newFuncCodes := make([]string, 0, len(infos)+2)
	newFuncCodes = append(newFuncCodes, "newCtx := new(Context)")
	newFuncCodes = append(newFuncCodes, "newCtx.ctx = context.Background()")
	for _, info := range infos {
		newFuncCodes = append(newFuncCodes, "newCtx."+info.memberName+" = "+info.memberName)
	}
	newFuncCodes = append(newFuncCodes, "return newCtx")
	newFunc := model.Function{
		Name:   "New",
		Params: newFuncParams,
		Return: []model.Field{
			{
				Type: "*Context",
			},
		},
		Code: newFuncCodes,
	}

	pkg := &model.Package{
		Name:    CONTEXT,
		Imports: imports,
		Structs: []model.Struct{
			{
				Name:    "Context",
				Fields:  fields,
				Methods: methods,
			},
		},
		Functions: []model.Function{
			newFunc,
		},
	}

	data, err := generator.GenerateFile(pkg)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(genPath, "context.go"))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}
