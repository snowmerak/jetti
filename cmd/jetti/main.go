package main

import (
	"github.com/alecthomas/kong"
	"github.com/snowmerak/jetti/internal/executor"
	"github.com/snowmerak/jetti/internal/executor/cli"
	"strings"
)

func main() {
	param := &cli.CLI{}
	ctx := kong.Parse(param)

	if err := ctx.PrintUsage(true); err != nil {
		return
	}

	switch ctx.Command() {
	case cli.Proto:
		switch {
		case param.Proto.New != "":
			executor.ProtoMake(param.Proto.New)
		case param.Proto.Build:
			executor.Proto()
		}
	case cli.Bean:
		if param.Bean.Generate {
			executor.Bean()
		}
	case cli.Cmd:
		switch {
		case param.Cmd.New != "":
			executor.CmdMake(param.Cmd.New)
		case param.Cmd.Build != "":
			executor.CmdBuild(param.Cmd.Build)
		case param.Cmd.Run != "":
			args := strings.Split(param.Cmd.Run, " ")
			executor.Cmd(args[0], args[1:]...)
		}
	}
}
