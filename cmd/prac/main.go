package main

import (
	"fmt"
	"github.com/snowmerak/jetti/v2/lib/parser"
)

func main() {
	pkg, err := parser.ParseFile("./lib/generator/gen.go")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", pkg)
}
