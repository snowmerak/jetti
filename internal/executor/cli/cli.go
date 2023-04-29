package cli

const (
	Proto  = "proto"
	Bean   = "bean"
	Cmd    = "cmd"
	New    = "new"
	Pprof  = "pprof"
	Redis  = "redis"
	Client = "client"
	Config = "config"
	Model  = "model"
)

type CLI struct {
	Proto struct {
		New      string `help:"Create a new proto file: <path>/<filename.proto>"`
		Generate bool   `help:"Generate proto all files"`
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
	Pprof struct {
		Http1 string `help:"Generate http1 pprof server: <addr>"`
		Http2 string `help:"Generate http2 pprof server: <addr>"`
		Http3 string `help:"Generate http3 pprof server: <addr>"`
	} `cmd:"" help:"Generate pprof server"`
	// Redis struct {
	//	New      string `help:"Create a new redis type: <path+name>"`
	//	Generate bool   `help:"Generate redis data types"`
	// } `cmd:"" help:"Generate redis data type"`
	Client struct {
		Rueidis string `help:"Generate redis client(rueidis): <addr>"`
		GoRedis string `help:"Generate redis client(go-redis): <addr>"`
		Nats    string `help:"Generate nats client: <addr>"`
	} `cmd:"" help:"Generate client"`
	Model struct {
		Json string `help:"Generate json model: <path+name>"`
		Yaml string `help:"Generate yaml model: <path+name>"`
		New  string `help:"Create a new model: <path+name>"`
	} `cmd:"" help:"Generate model"`
	Config struct {
		New     string `help:"Create a new jsonnet config file: <path+name>"`
		Jsonnet string `help:"Generate jsonnet config: <path+name>"`
	} `cmd:"" help:"Generate config"`
}
