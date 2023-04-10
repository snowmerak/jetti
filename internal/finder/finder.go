package finder

import (
	"bufio"
	"github.com/snowmerak/go-bean/internal/model"
	"io"
	"strings"
)

func FindStruct(r io.Reader, direction string) model.Structs {
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
