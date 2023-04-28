package main

import (
	"github.com/alecthomas/kong"
	"github.com/snowmerak/jetti/internal/executor"
	"github.com/snowmerak/jetti/internal/executor/cli"
)

func main() {
	param := &cli.CLI{}
	ctx := kong.Parse(param)

	switch ctx.Command() {
	case cli.New:
		switch {
		case param.New.Init != "":
			executor.Init(param.New.Init)
		}
	case cli.Proto:
		switch {
		case param.Proto.New != "":
			executor.ProtoNew(param.Proto.New)
		case param.Proto.Generate:
			executor.ProtoBuild()
		}
	case cli.Bean:
		if param.Bean.Generate {
			executor.Bean()
		}
	case cli.Cmd:
		switch {
		case param.Cmd.New != "":
			executor.CmdNew(param.Cmd.New)
		case param.Cmd.Build != nil:
			executor.CmdBuild(param.Cmd.Build[0], param.Cmd.Build[1:]...)
		case param.Cmd.Run != nil:
			executor.CmdRun(param.Cmd.Run[0], param.Cmd.Run[1:]...)
		}
	case cli.Pprof:
		if param.Pprof.Http1 != "" {
			executor.PprofHttp1(param.Pprof.Http1)
		}
		if param.Pprof.Http2 != "" {
			executor.PprofHttp2(param.Pprof.Http2)
		}
		if param.Pprof.Http3 != "" {
			executor.PprofHttp3(param.Pprof.Http3)
		}
	// case cli.Redis:
	//	if param.Redis.New != "" {
	//		executor.RedisNew(param.Redis.New)
	//	}
	//	if param.Redis.Generate {
	//		executor.RedisGenerate()
	//	}
	case cli.Client:
		if param.Client.Rueidis != "" {
			executor.ClientRueidis(param.Client.Rueidis)
		}
		if param.Client.GoRedis != "" {
			executor.ClientGoRedis(param.Client.GoRedis)
		}
		if param.Client.Nats != "" {
			executor.ClientNats(param.Client.Nats)
		}
	}
}
