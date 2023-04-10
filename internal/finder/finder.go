package finder

import (
	"bufio"
	"github.com/snowmerak/go-bean/internal/model"
	"io"
	"strings"
)

func Find(r io.Reader) model.Model {
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
		case strings.HasPrefix(value, "//go:bean"):
			line, _, err = reader.ReadLine()
			if err != nil {
				break
			}
			value = string(line)
			if strings.HasPrefix(value, "type ") {
				sp := strings.SplitN(value, " ", 3)
				structNames = append(structNames, sp[1])
			}
		}
	}
	return model.Model{
		PackageName: packageName,
		StructNames: structNames,
	}
}
