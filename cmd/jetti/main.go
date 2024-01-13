package main

import (
	"github.com/alecthomas/kong"
	"github.com/snowmerak/jetti/v2/internal/executor"
	"github.com/snowmerak/jetti/v2/internal/executor/cli"
	"log"
	"os"
	"path/filepath"
)

func main() {
	param := &cli.CLI{}
	ctx := kong.Parse(param)

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	pwd = filepath.ToSlash(pwd)

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
	case cli.Index:
		if err := executor.Index(pwd); err != nil {
			panic(err)
		}
	case cli.Impl:
		if err := executor.ImplInteractive(pwd); err != nil {
			panic(err)
		}
	case cli.ImplTarget:
		if err := executor.ImplTargets(pwd, param.Impl.Target); err != nil {
			panic(err)
		}
	case cli.Check:
		if err := executor.Check(pwd); err != nil {
			panic(err)
		}
	case cli.Tools:
		if param.Tools.Renew {
			if err := executor.InstallRegistriesRenew(); err != nil {
				panic(err)
			}
			log.Println("tools registries renewed")
		}
		if param.Tools.Install {
			if param.Tools.Multi {
				if err := executor.InstallMultipleRegistries(); err != nil {
					panic(err)
				}
			} else {
				if err := executor.InstallRegistry(); err != nil {
					panic(err)
				}
			}
		}
	default:
		log.Println("unknown command", ctx.Command())
	}
}
