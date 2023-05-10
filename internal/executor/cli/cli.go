package cli

const (
	Generate = "generate"
	New      = "new <module-name>"
	Run      = "run <command-name>"
	RunArgs  = "run <command-name> <args>"
)

type CLI struct {
	Generate struct {
	} `cmd:"" help:"Generate code"`
	New struct {
		ModuleName string `arg:"" help:"Module name"`
		Cmd        bool   `cmd:"" help:"Create a new command"`
	} `cmd:"" help:"Create a new project"`
	Run struct {
		CommandName string   `arg:"" help:"Command name"`
		Args        []string `arg:"" optional:"" help:"Command arguments"`
	} `cmd:"" help:"Run a command"`
}
