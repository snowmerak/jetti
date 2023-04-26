package finder

import (
	"bufio"
	"github.com/snowmerak/jetti/internal/model"
	"io"
	"os"
	"strings"
)

func FindStructName(r io.Reader, direction string) model.Structs {
	prefix := "//go:" + direction
	packageName := ""
	structNames := []string(nil)
	reader := bufio.NewReader(r)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		value := string(line)
		switch {
		case packageName == "" && strings.HasPrefix(value, "package "):
			packageName = strings.TrimPrefix(value, "package ")
		case strings.HasPrefix(value, prefix):
			for {
				line, _, err = reader.ReadLine()
				if err != nil {
					break
				}
				value = string(line)
				if strings.HasPrefix(value, "type ") {
					sp := strings.SplitN(value, " ", 3)
					structNames = append(structNames, sp[1])
					break
				}
			}
		}
	}
	return model.Structs{
		PackageName: packageName,
		StructNames: structNames,
	}
}

func FindStructType(r io.Reader, direction string) model.Struct {
	prefix := "//go:" + direction
	packageName := ""
	structName := ""
	fields := []model.Field(nil)
	reader := bufio.NewReader(r)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		value := string(line)
		switch {
		case packageName == "" && strings.HasPrefix(value, "package "):
			packageName = strings.TrimPrefix(value, "package ")
		case strings.HasPrefix(value, prefix):
			for {
				line, _, err = reader.ReadLine()
				if err != nil {
					break
				}
				value = string(line)
				if strings.HasPrefix(value, "type ") {
					sp := strings.SplitN(value, " ", 3)
					structName = sp[1]
					break
				}
			}
		case strings.HasPrefix(value, "type "):
			sp := strings.SplitN(value, " ", 3)
			structName = sp[1]
		case strings.HasPrefix(value, "type "+structName+" struct"):
			for {
				line, _, err = reader.ReadLine()
				if err != nil {
					break
				}
				value = string(line)
				if strings.HasPrefix(value, "}") {
					break
				}
				if strings.HasPrefix(value, "//") {
					continue
				}
				sp := strings.SplitN(value, " ", 2)
				if len(sp) != 2 {
					continue
				}
				name := sp[0]
				sp = strings.SplitN(sp[1], " ", 2)
				if len(sp) != 2 {
					continue
				}
				typ := sp[0]
				tags := make(map[string]string)
				if strings.HasPrefix(sp[1], "`") {
					tag := strings.Trim(sp[1], "`")
					sp = strings.Split(tag, " ")
					for _, v := range sp {
						sp2 := strings.SplitN(v, ":", 2)
						if len(sp2) != 2 {
							continue
						}
						tags[sp2[0]] = sp2[1]
					}
				}
				fields = append(fields, model.Field{
					Name: name,
					Type: typ,
					Tags: tags,
				})
			}
		}
	}
	return model.Struct{
		PackageName: packageName,
		StructName:  structName,
		Fields:      fields,
	}
}

func FindDirections(r io.Reader, direction string) []string {
	prefix := "//go:" + direction
	reader := bufio.NewReader(r)
	result := []string(nil)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		value := string(line)
		if strings.HasPrefix(value, prefix) {
			sp := strings.SplitN(value, " ", 2)
			if len(sp) != 2 {
				continue
			}
			result = append(result, sp[1])
		}
	}

	return result
}

func FindModuleName() (string, error) {
	f, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	reader := bufio.NewReader(f)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return "", err
		}

		if strings.HasPrefix(string(line), "module ") {
			return strings.TrimPrefix(string(line), "module "), nil
		}
	}
}

type Field struct {
	Name string
	Type string
}

func FindProtobufFields(file string, message string) []string {
	f, err := os.Open(file)
	if err != nil {
		return nil
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	reader := bufio.NewReader(f)
	result := []string(nil)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		value := string(line)
		if strings.HasPrefix(value, "type "+message) {
			for {
				line, _, err = reader.ReadLine()
				if err != nil {
					break
				}
				value = string(line)
				if strings.HasPrefix(value, "}") {
					break
				}
				if strings.HasPrefix(value, "//") {
					continue
				}
				sp := strings.SplitN(value, " ", 2)
				if len(sp) != 2 {
					continue
				}
				result = append(result, sp[0])
			}
		}
	}

	return result
}
