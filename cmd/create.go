/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/rickliujh/multi-signer/pkg/fileio"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := "default"
		if len(args) != 0 {
			name = args[0]
		}

		if err := fileio.CreateGroup(name); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	groupCmd.AddCommand(createCmd)
}
