package commands

import (
	"github.com/spf13/cobra"
	"github.com/stumacwastaken/todo/cmd/commands/rest"
	"github.com/stumacwastaken/todo/cmd/commands/seed"
)

var (
	rootCmd = &cobra.Command{
		Use:           "todo",
		Short:         "todo all the things",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(seed.Cmd)
	rootCmd.AddCommand(rest.Cmd)

}
