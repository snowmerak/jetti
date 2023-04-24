package cli

const (
	Proto = "proto"
	Bean  = "bean"
	Cmd   = "cmd"
	New   = "new"
)

type CLI struct {
	Proto struct {
		New   string `help:"Create a new proto file: <path>/<filename.proto>"`
		Build bool   `help:"Build proto all files"`
	} `cmd:"" help:"Generate protobuf messages and grpc services"`
	Bean struct {
		Generate bool `help:"Generate bean container"`
	} `cmd:"" help:"Generate bean container"`
	Cmd struct {
		New   string   `help:"Create a new cmd package: <cmd-name>"`
		Build []string `help:"Build a cmd package: <cmd-name>,[<options>...]"`
		Run   []string `help:"Run a cmd package: <cmd-name>,[<args>...]"`
	} `cmd:"" help:"managing cmd package"`
	New struct {
		Init string `help:"Initialize a new project"`
	} `cmd:"" help:"Initialize a new project"`
}
