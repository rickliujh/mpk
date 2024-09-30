package cmd

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/bnb-chain/tss-lib/v2/common"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/signing"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/rickliujh/mpk/pkg/fileio"
	"github.com/spf13/cobra"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Signing a message as peer group",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) == 0 {
			return fmt.Errorf("no message provided for signing")
		}

		message := args[0]
		meta, err := fileio.LoadFile[*fileio.Meta](group, "meta")
		if err != nil {
			return
		}

		pks, err := fileio.LoadPK(group)
		if err != nil {
			return
		}

		threshold := meta.Threshold
		pids := meta.Peers
		sorted := tss.SortPartyIDs(pids)
		ctx := tss.NewPeerContext(sorted)

		// Select an elliptic curve
		// use ECDSA
		curve := tss.S256()
		// or use EdDSA
		// curve := tss.Edwards()

		ended := 0
		msg := big.NewInt(0).SetBytes([]byte(message))
		outCh := make(chan tss.Message, len(pids))
		endCh := make(chan *common.SignatureData, len(pids))

		parties := make(map[string]tss.Party)
		var signatureData *common.SignatureData

		for i, id := range pids {
			pk := pks[id.GetMoniker()]
			params := tss.NewParameters(curve, ctx, pids[i], len(pids), threshold)
			party := signing.NewLocalParty(msg, params, *pk, outCh, endCh)
			parties[id.GetMoniker()] = party
			go func() {
				err := party.Start()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}()
		}

		// go func() {
		// 	for data := range sendCh {
		// 		fmt.Println(&data)
		// 	}
		// }()
	signing:
		for {
			select {
			case msg := <-outCh:
				go func() {
					dest := msg.GetTo()
					if dest == nil {
						for k, p := range parties {
							if k == msg.GetFrom().GetMoniker() {
								continue
							}
							fmt.Println("no dest")
							if err := SharedPartyUpdater(p, msg); err != nil {
								fmt.Println(err)
							}
						}
					} else {
						fmt.Println("dest")
						if err := SharedPartyUpdater(parties[dest[0].GetMoniker()], msg); err != nil {
							fmt.Println(err)
						}
					}
				}()
			case data := <-endCh:
				ended++
				if ended >= len(parties) {
					fmt.Println("test")
					fmt.Println(data)

					signatureData = data
					break signing
				}
			case <-time.After(time.Duration(timeout) * time.Minute):
				fmt.Printf("%d, %d, %d", len(outCh), len(endCh), 1)
				err = fmt.Errorf("signing timeout\n")
				break signing
				// case <-time.Tick(10 * time.Second):
				// 	for _, p := range paries {
				// 		fmt.Printf("%v runing state %v\n", p.PartyID().Id, p.Running())
				// 	}
			}
		}

		if output != "" {
			err = fileio.SaveSig(output, signatureData)
			if err != nil {
				return
			}
		}

		fmt.Printf("====signature====\n%v\n", signatureData)
		return
	},
}

var (
	keyname string
	output  string
)

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringVarP(&group, "group", "g", "default", "the peer group save the keys")
	signCmd.Flags().Int64VarP(&timeout, "timeout", "o", 1, "defines the minutes of timeout for preparams")
	signCmd.Flags().StringVarP(&output, "output", "f", "signature.json", "file path for output of signature")
}
