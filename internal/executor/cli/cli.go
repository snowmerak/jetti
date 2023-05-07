package cli

const (
	Generate = "generate"
	New      = "new <module-name>"
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
	New struct {
		ModuleName string `arg:"" help:"Module name"`
	} `cmd:"" help:"Create a new project"`
}
