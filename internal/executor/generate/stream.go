package generate

import (
	"fmt"
	"github.com/snowmerak/jetti/v2/internal/executor/check"
	"github.com/snowmerak/jetti/v2/lib/generator"
	"github.com/snowmerak/jetti/v2/lib/model"
	"os"
	"path/filepath"
	"strings"
)

func Stream(path string, streams []check.Stream) error {
	dir := filepath.Dir(path)

	for _, stream := range streams {
		fileName := MakeGeneratedFileName(dir, strings.ToLower(stream.StructName), "stream")

		pkg := &model.Package{
			Name:    stream.PackageName,
			Imports: stream.Imports,
		}

		funcSeq := 0
		paramSeq := 0
		returnSeq := 0

		makeParam := func(funcSeq int, paramSeq int) string {
			return fmt.Sprintf("r%03d%03d", funcSeq, paramSeq)
		}
		makeReturn := func(funcSeq int, returnSeq int) string {
			return fmt.Sprintf("r%03d%03d", funcSeq, returnSeq)
		}

		params := make([]model.Field, len(stream.FunctionSignature[0].Params))
		for i, param := range stream.FunctionSignature[0].Params {
			params[i] = model.Field{
				Name: makeParam(funcSeq, paramSeq),
				Type: param.Type,
			}
			paramSeq++
		}
		paramSeq = 0
		funcSeq++

		lastIndex := len(stream.FunctionSignature) - 1
		returns := make([]model.Field, len(stream.FunctionSignature[lastIndex].Return))
		hasError := false
		errIdx := 0
		for i, ret := range stream.FunctionSignature[lastIndex].Return {
			returns[i] = model.Field{
				Name: "",
				Type: ret.Type,
			}
			returnSeq++
			if ret.Type == "error" {
				hasError = true
				errIdx = i
			}
		}
		returnSeq = 0

		if !hasError {
			returns = append(returns, model.Field{
				Name: "",
				Type: "error",
			})
			errIdx = len(returns) - 1
		}

		method := model.Method{
			Receiver: model.Field{
				Name: strings.ToLower(stream.StructName[:1]),
				Type: stream.StructName,
			},
			Name:   stream.StructName + "Stream",
			Params: params,
			Return: returns,
		}

		code := strings.Builder{}
		for _, fun := range stream.FunctionSignature {
			for j, ret := range fun.Return {
				code.WriteString(fmt.Sprintf("var %s %s\n", makeReturn(funcSeq, j), ret.Type))
			}

			hasInnerError := false
			innerErrorIdx := 0
			for j := range fun.Return {
				if j > 0 {
					code.WriteString(", ")
				}
				code.WriteString(fmt.Sprintf("%s", makeReturn(funcSeq, j)))

				if returns[j].Type == "error" {
					hasInnerError = true
					innerErrorIdx = j
				}
			}

			code.WriteString(fmt.Sprintf(" = $RECEIVER$.%s(", fun.Name))

			for j := range fun.Params {
				if j > 0 {
					code.WriteString(", ")
				}
				code.WriteString(makeParam(funcSeq-1, j))
			}

			code.WriteString(")\n")

			if hasInnerError {
				code.WriteString(fmt.Sprintf("if %s != nil {\n", makeReturn(funcSeq, innerErrorIdx)))
				code.WriteString("return ")
				for j := range returns {
					if j > 0 {
						code.WriteString(", ")
					}
					switch j {
					case innerErrorIdx:
						code.WriteString(fmt.Sprintf("%s", makeReturn(funcSeq, innerErrorIdx)))
					default:
						code.WriteString(fmt.Sprintf("*new(%s)", returns[j].Type))
					}
				}
				code.WriteString("}\n")
			}

			method.Code = append(method.Code, code.String())
			code.Reset()

			funcSeq++
		}

		code.WriteString("return ")
		for i := range returns {
			if i > 0 {
				code.WriteString(", ")
			}
			switch i {
			case errIdx:
				code.WriteString("nil")
			default:
				code.WriteString(makeReturn(funcSeq-1, i))
			}
		}

		method.Code = append(method.Code, code.String())
		code.Reset()

		pkg.Methods = append(pkg.Methods, method)

		value, err := generator.GenerateFile(pkg)
		if err != nil {
			return err
		}

		f, err := os.Create(fileName)
		if err != nil {
			return err
		}

		if _, err = f.Write(value); err != nil {
			return err
		}
	}

	return nil
}
