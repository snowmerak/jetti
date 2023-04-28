package executor

import (
	"os"
	"strings"
)

func ClientNewJson(name string) {
	dep := getDependency(name)

	folder := makeSubPath(clientFolder, dep.Import)

	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		panic(err)
	}

	f, err := os.Create(folder + strings.ToLower(dep.Type) + ".json")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := f.WriteString(`{}`); err != nil {
		panic(err)
	}
}
