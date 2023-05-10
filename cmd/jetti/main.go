package main

import (
	"fmt"
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

	fmt.Println(ctx.Command())
	switch ctx.Command() {
	case cli.Generate:
		if err := executor.Generate(pwd); err != nil {
			panic(err)
		}
	case cli.New:
		if err := executor.New(pwd, param.New.ModuleName, param.New.Cmd); err != nil {
			panic(err)
		}
	case cli.Run:
	case cli.RunArgs:
	}
}
