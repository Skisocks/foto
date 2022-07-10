package upload

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
	uploadLong = command.LongDesc(`
		Import one or more resources.
		` + validResources + `
`)

	uploadExample = command.Examples(`
		# Import photos from dir
		foto import photos
	`)
)

// NewCmdUpload creates a command object for the generic "import" action, which
// retrieves one or more resources from a file.
func NewCmdUpload() *cobra.Command {
	o := &Options{}

	cmd := &cobra.Command{
		Use:     "upload TYPE [flags]",
		Short:   "Upload one or more resources",
		Long:    uploadLong,
		Example: uploadExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Cmd = cmd
			o.Args = args
			err := o.Run()
			command.CheckErr(err)
		},
		SuggestFor: []string{"up", "load"},
	}
	cmd.AddCommand(NewCmdPhoto())
	return cmd
}

// Run implements this command
func (o *Options) Run() error {
	return o.Cmd.Help()
}
