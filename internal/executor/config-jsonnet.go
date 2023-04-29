package executor

import (
	"github.com/google/go-jsonnet"
	"os"
	"path/filepath"
	"strings"
)

func ConfigNew(path string) {
	if err := os.MkdirAll("template/config/"+filepath.Dir(path), os.ModePerm); err != nil {
		panic(err)
	}

	f, err := os.Create("template/config/" + path)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := f.WriteString(`local config = {
name: "louis",
age: 21,
address: "world",
};
		`); err != nil {
		panic(err)
	}
}

func ClientJsonnet(path string) {
	vm := jsonnet.MakeVM()

	jsonnetPath := "template/config/" + path

	result, err := vm.EvaluateFile(jsonnetPath)
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(generated+"/config/"+filepath.Dir(path), os.ModePerm); err != nil {
		panic(err)
	}

	f, err := os.Create(generated + "/config/" + filepath.Dir(path) + "/" + strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) + ".json")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := f.WriteString(result); err != nil {
		panic(err)
	}
}
