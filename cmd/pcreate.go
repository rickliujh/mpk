/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"math/big"

	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/google/uuid"
	"github.com/rickliujh/multi-signer/pkg/fileio"
	"github.com/spf13/cobra"
)

// peerCreateCmd represents the pcreate command
var peerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creating peers in the group",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) == 0 {
			return fmt.Errorf("no enough peer name specified\n")
		}

		pids := make([]*tss.PartyID, len(args))
		for i, n := range args {
			id, err := uuid.New().MarshalBinary()
			if err != nil {
				return err
			}
			party := tss.NewPartyID(fmt.Sprint(i), n, big.NewInt(0).SetBytes(id))
			fmt.Printf("peer [%s] created: %v\n", n, party)
			pids[i] = party
		}

		return fileio.SaveFile(group, "meta", &fileio.Meta{
			Threshold: threshold,
			Peers:     pids,
		})
	},
}

var (
	threshold int
)

func init() {
	peerCmd.AddCommand(peerCreateCmd)

	peerCreateCmd.Flags().StringVarP(&group, "group", "g", "default", "the peer group save the keys")
	peerCreateCmd.Flags().IntVarP(&threshold, "threshold", "t", 1, "defines the threshold for signature verification")
}
