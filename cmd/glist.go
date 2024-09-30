/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/rickliujh/mpk/pkg/fileio"
	"github.com/spf13/cobra"
)

// groupListCmd represents the list command
var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List existing peer groups",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		gps, err := fileio.LoadGroup()
		for _, gp := range gps {
			fmt.Println(gp)
		}
		return err
	},
}

func init() {
	groupCmd.AddCommand(groupListCmd)

}
