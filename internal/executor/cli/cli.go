package cli

const (
	Generate   = "generate"
	New        = "new <module-name>"
	Run        = "run <command-name>"
	RunArgs    = "run <command-name> <args>"
	Show       = "show"
	Index      = "index"
	Impl       = "impl"
	ImplTarget = "impl <target>"
	Check      = "check"
	Tools      = "tools"
)

type CLI struct {
	Generate struct {
	} `cmd:"" help:"Generate code"`
	New struct {
		ModuleName string `arg:"" help:"Module name"`
		Cmd        bool   `cmd:"" help:"Create a new command"`
		Proto      bool   `cmd:"" help:"Create a new proto"`
	} `cmd:"" help:"Create a new project"`
	Run struct {
		CommandName string   `arg:"" help:"Command name"`
		Args        []string `arg:"" optional:"" help:"Command arguments"`
	} `cmd:"" help:"Run a command"`
	Show struct {
		Imports bool `cmd:"" help:"Show import cycle"`
	} `cmd:"" help:"Show information"`
	Index struct {
	} `cmd:"" help:"Index the project"`
	Impl struct {
		Interactive bool     `cmd:"" help:"Interactive mode"`
		Target      []string `arg:"" optional:"" help:"Implement interfaces"`
	} `cmd:"" help:"Implement interfaces"`
	Check struct {
	} `cmd:"" help:"Check the project"`
	Tools struct {
		Renew   bool `cmd:"" help:"Renew tools registry"`
		Install bool `cmd:"" help:"Install tool"`
		Multi   bool `cmd:"" help:"Install multiple tools"`
	} `cmd:"" help:"Search and install tools"`
}
