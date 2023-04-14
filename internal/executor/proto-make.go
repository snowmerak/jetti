package executor

import (
	"fmt"
	"os"
	"path/filepath"
)

func ProtoMake(path string) {
	packageName := "." + filepath.Dir(path)
	path = filepath.Join("proto", path)
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		panic("file already exists")
	}

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := f.Write([]byte(fmt.Sprintf("syntax = \"proto3\";\n\noption go_package = \"%s\";\n\npackage %s;\n\nmessage HelloRequest {\n  string name = 1;\n}\n\nmessage HelloResponse {\n  string message = 1;\n}\n\nservice HelloService {\n  rpc SayHello (HelloRequest) returns (HelloResponse);\n}\n\n", packageName, name))); err != nil {
		panic(err)
	}
}
