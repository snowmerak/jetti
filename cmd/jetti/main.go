package main

import (
	"flag"
	"github.com/snowmerak/jetti/internal/executor"
)

func main() {
	initFlag := flag.String("init", "", "-init <project-name> : initialize go project")
	beanFlag := flag.Bool("bean", false, "-bean : generate bean container")
	protoFlag := flag.Bool("proto", false, "-proto : generate protobuf messages and grpc services")
	protoMakeFlag := flag.String("proto-make", "", "-proto-make <path>/<filename.proto>")
	helpFlag := flag.Bool("help", false, "show help")
	flag.Parse()

	if *helpFlag {
		flag.PrintDefaults()
	}
	if *initFlag != "" {
		executor.Init(*initFlag)
	}
	if *beanFlag {
		executor.Bean()
	}
	if *protoFlag {
		executor.Proto()
	}
	if *protoMakeFlag != "" {
		executor.ProtoMake(*protoMakeFlag)
	}
}
