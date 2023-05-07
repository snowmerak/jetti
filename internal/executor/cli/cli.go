package cli

const (
	Generate = "generate"
)

/*
TODO: add proto
TODO: add fbs
TODO: add bean
TODO: add cmd
TODO: add pprof
TODO: add server
TODO: add json/yaml
*/
type CLI struct {
	Generate struct {
	} `cmd:"" help:"Generate code"`
}
