/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/rickliujh/multi-signer/pkg/fileio"
	"github.com/spf13/cobra"
)

var (
	timeout int64
	group   string
)

// keyCmd represents the key command
var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// preParams, err := keygen.GeneratePreParams(time.Duration(timeout) * time.Minute)
		// if err != nil {
		// 	fmt.Println(err, preParams)
		// 	return
		// }

		meta, err := fileio.LoadFile[fileio.Meta](group, "meta")
		if err != nil {
			return
		}

		pids := meta.Peers
		if len(pids) == 0 {
			return fmt.Errorf("no peers in the group, create peers first")
		}

		threshold := meta.Threshold
		sorted := tss.SortPartyIDs(pids)
		ctx := tss.NewPeerContext(sorted)

		// Select an elliptic curve
		// use ECDSA
		curve := tss.S256()
		// or use EdDSA
		// curve := tss.Edwards()

		wg := &sync.WaitGroup{}
		pks := make(map[string]*keygen.LocalPartySaveData)
		paries := []tss.Party{}
		for i, p := range sorted {
			wg.Add(2)
			params := tss.NewParameters(curve, ctx, p, len(pids), threshold)
			outCh := make(chan tss.Message)
			endCh := make(chan *keygen.LocalPartySaveData)
			party := keygen.NewLocalParty(params, outCh, endCh) // Omit the last arg to compute the pre-params in round 1
			paries = append(paries, party)
			go func(i int) {
				defer wg.Done()
				for {
					select {
					case msg := <-outCh:
						for _, p := range paries {
							if p.PartyID().Index == msg.GetFrom().Index {
								continue
							}
							go func() {
								if err := SharedPartyUpdater(p, msg); err != nil {
									fmt.Fprintln(os.Stderr, err)
								}
							}()
						}
					case data := <-endCh:
						pks[sorted[i].GetMoniker()] = data
						return
					case <-time.After(time.Duration(timeout) * time.Minute):
						fmt.Printf("goroutine[%d]	keygen timeout\n", i)
						return
					}
				}
			}(i)
			go func() {
				defer wg.Done()
				err := party.Start()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}
			}()
		}

		wg.Wait()

		fmt.Println("====public key====")
		tmp := pks[sorted[0].GetMoniker()]
		fmt.Println(tmp.ECDSAPub.ToECDSAPubKey())
		err = fileio.SaveFile(group, "public", tmp.ECDSAPub.ToECDSAPubKey())
		fmt.Println("====private keys====")
		for k, pk := range pks {
			fmt.Printf("[%s]\n%v\n", k, pk)
			if err = fileio.SavePK(group, k, pk); err != nil {
				return err
			}
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(keygenCmd)

	keygenCmd.Flags().Int64VarP(&timeout, "timeout", "o", 1, "defines the minutes of timeout for preparams")
	keygenCmd.Flags().StringVarP(&group, "group", "g", "default", "the peer group save the keys")
}

func SharedPartyUpdater(party tss.Party, msg tss.Message) *tss.Error {
	s, err := json.Marshal(msg)
	fmt.Printf("%s, %v, %s\n", s, err, party.PartyID().GetMoniker())

	// do not send a message from this party back to itself
	if party.PartyID() == msg.GetFrom() {
		fmt.Println("skip")
		return nil
	}
	bz, _, err := msg.WireBytes()
	if err != nil {
		return party.WrapError(err)
	}
	pMsg, err := tss.ParseWireMessage(bz, msg.GetFrom(), msg.IsBroadcast())
	if err != nil {
		return party.WrapError(err)
	}
	if _, err := party.Update(pMsg); err != nil {
		return err
	}
	return nil
}
