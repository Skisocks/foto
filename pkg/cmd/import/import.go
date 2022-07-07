package _import

import (
	"github.com/skisocks/foto/pkg/helpers/command"
	"github.com/spf13/cobra"
)

// Options for triggering
type Options struct {
	Apps []string
	Args []string
	Cmd  *cobra.Command
}

const (
	validResources = `Valid resource types include:
	* photos (aka 'ph')
	`
)

var (
	getLong = command.LongDesc(`
		Import one or more resources.
		` + validResources + `
`)

	getExample = command.Examples(`
		# Import photos from dir
		foto import photos
	`)
)

// NewCmdImport creates a command object for the generic "import" action, which
// retrieves one or more resources from a file.
func NewCmdImport() *cobra.Command {
	o := &Options{}

	cmd := &cobra.Command{
		Use:     "import TYPE [flags]",
		Short:   "Import one or more resources",
		Long:    getLong,
		Example: getExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Cmd = cmd
			o.Args = args
			err := o.Run()
			command.CheckErr(err)
		},
		SuggestFor: []string{"list", "ps"},
	}
	cmd.AddCommand(NewCmdPhoto())
	return cmd
}

// Run implements this command
func (o *Options) Run() error {
	return o.Cmd.Help()
}
