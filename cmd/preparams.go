/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/spf13/cobra"
)

// dirty code
const preparamsName = "preparams.json"

var (
	timeout int64
	cfgdir  string
)

// preparamsCmd represents the preparams command
var preparamsCmd = &cobra.Command{
	Use:   "preparams",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// should allow update or start interactive action to confirm if updating
		f, err := preparefile(cfgdir + preparamsName)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()

		fmt.Printf("generating preparams, timeout: %d\n", timeout)

		preParams, err := keygen.GeneratePreParams(time.Duration(timeout) * time.Minute)
		if err != nil {
			fmt.Println(err, preParams)
			return
		}

		bz, err := json.Marshal(preParams)
		if err != nil {
			fmt.Println(err, preParams)
			return
		}

		f.Write(bz)
	},
}

func init() {
	genCmd.AddCommand(preparamsCmd)

	// dirty code
	usercfg, _ := os.UserConfigDir()
	cfgdir = usercfg + "/multi-signer/"

	preparamsCmd.Flags().Int64VarP(&timeout, "timeout", "o", 1, "defines the minutes of timeout for preparams")
}

// preparedir checks the file, returns file and err if exist
// creates the dir and file otherwise
func preparefile(path string) (*os.File, error) {
	finfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(cfgdir, 0700); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if finfo != nil {
		return nil, fmt.Errorf("preparams already existed")
	}

	f, err := os.Create(cfgdir + preparamsName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return f, nil
}
