package cmd

import (
	log "github.com/sirupsen/logrus"
	_import "github.com/skisocks/foto/pkg/cmd/import"
	"github.com/skisocks/foto/pkg/cmd/upload"
	"github.com/skisocks/foto/pkg/cmd/version"
	"github.com/skisocks/foto/pkg/helpers/command"
	"github.com/skisocks/foto/pkg/rootcmd"
	"github.com/spf13/cobra"
)

// Main creates the new command
func Main() *cobra.Command {
	cmd := &cobra.Command{
		Use:   rootcmd.TopLevelCommand,
		Short: "commands for creating and upgrading the MQube environment",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				log.Errorf(err.Error())
			}
		},
	}
	cmd.AddCommand(_import.NewCmdImport())
	cmd.AddCommand(upload.NewCmdUpload())
	cmd.AddCommand(command.SplitCommand(version.NewCmdVersion()))
	return cmd
}
