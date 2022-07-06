package cmd

import (
	log "github.com/sirupsen/logrus"
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
	return cmd
}
