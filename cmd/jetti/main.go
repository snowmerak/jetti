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
		kind := executor.NewKindModule
		if param.New.Cmd {
			kind = executor.NewKindCmd
		}
		if param.New.Proto {
			kind = executor.NewKindProto
		}
		if err := executor.New(pwd, param.New.ModuleName, kind); err != nil {
			panic(err)
		}
	case cli.Run:
		fallthrough
	case cli.RunArgs:
		if err := executor.Run(pwd, param.Run.CommandName, param.Run.Args...); err != nil {
			panic(err)
		}
	case cli.Show:
		if param.Show.Imports {
			if err := executor.ShowImports(pwd); err != nil {
				panic(err)
			}
		}
	}
}
