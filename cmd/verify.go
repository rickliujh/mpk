/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/rickliujh/mpk/pkg/fileio"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verifying the signature",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if sigpath == "" {
			return fmt.Errorf("provided file path invalid")
		}

		if len(args) == 0 {
			return fmt.Errorf("no message provided for verifying")
		}

		bmsg := []byte(args[0])

		sig, err := fileio.LoadSig(sigpath)
		if err != nil {
			return
		}

		pks, err := fileio.LoadPK(group)
		if err != nil {
			return
		}

		for k, pk := range pks {
			pub := pk.ECDSAPub.ToECDSAPubKey()
			if ok := ecdsa.Verify(
				pub,
				bmsg, big.NewInt(0).SetBytes(sig.R),
				big.NewInt(0).SetBytes(sig.S),
			); !ok {
				return fmt.Errorf("pk[%s] failed to verify", k)
			}
		}

		fmt.Println("verified")
		return
	},
}

var (
	sigpath string
)

func init() {
	rootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVarP(&sigpath, "sig-file", "f", "signature.json", "file path for signature file to be verified")
	verifyCmd.Flags().StringVarP(&group, "group", "g", "default", "the peer group save the keys")
}
