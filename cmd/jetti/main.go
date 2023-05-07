package main

import (
	"github.com/alecthomas/kong"
	"github.com/snowmerak/jetti/internal/executor"
	"github.com/snowmerak/jetti/internal/executor/cli"
	"os"
)

func main() {
	param := &cli.CLI{}
	ctx := kong.Parse(param)

	switch ctx.Command() {
	case cli.Generate:
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		if err := executor.Generate(pwd); err != nil {
			panic(err)
		}
	}
}
