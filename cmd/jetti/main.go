package main

import (
	"github.com/alecthomas/kong"
	"github.com/snowmerak/jetti/v2/internal/executor"
	"github.com/snowmerak/jetti/v2/internal/executor/cli"
	"os"
)

func main() {
	param := &cli.CLI{}
	ctx := kong.Parse(param)

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	switch ctx.Command() {
	case cli.Generate:
		if err := executor.Generate(pwd); err != nil {
			panic(err)
		}
	case cli.New:
		if err := executor.New(pwd, param.New.ModuleName); err != nil {
			panic(err)
		}
	}
}
