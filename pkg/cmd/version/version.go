package version

import (
	"fmt"
	"github.com/skisocks/foto/pkg/helpers/command"
	"github.com/skisocks/foto/pkg/rootcmd"
	"github.com/spf13/cobra"
)

// Options for triggering
type Options struct {
	Apps []string
	Args []string
	Cmd  *cobra.Command
}

const (
	// TestVersion used in test cases for the current version if no
	// version can be found - such as if the version property is not properly
	// included in the go test flags
	TestVersion = "1.0.0-SNAPSHOT"
)

var (
	createLong = command.LongDesc(`
		Shows the version of mqa
`)

	createExample = command.Examples(`
		version
	`)

	Version string
)

// NewCmdTrigger
func NewCmdVersion() (*cobra.Command, *Options) {

	o := &Options{}

	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Shows the version of foto",
		Long:    createLong,
		Example: fmt.Sprintf(createExample, rootcmd.BinaryName),
		Run: func(cmd *cobra.Command, args []string) {
			o.Cmd = cmd
			o.Args = args
			err := o.Run()
			command.CheckErr(err)
		},
	}
	o.Cmd = cmd

	return cmd, o

}

func (o *Options) Run() error {
	fmt.Println(GetVersion())

	return nil
}

func GetVersion() string {
	if Version != "" {
		return Version
	}
	return TestVersion
}
