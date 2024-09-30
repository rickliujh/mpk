/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/rickliujh/multi-signer/pkg/fileio"
	"github.com/spf13/cobra"
)

// peerListCmd represents the plist command
var peerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List existing peers in group",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		meta, err := fileio.LoadFile[*fileio.Meta](group, "meta")
		if err != nil {
			return err
		}

		for i, pn := range meta.Peers {
			fmt.Printf("%d. %s\n", i, pn.GetMoniker())
		}

		return
	},
}

func init() {
	peerCmd.AddCommand(peerListCmd)

	peerListCmd.Flags().StringVarP(&group, "group", "g", "default", "the peer group save the keys")
}
