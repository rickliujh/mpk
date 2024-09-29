/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/google/uuid"
	"github.com/rickliujh/multi-signer/pkg/fileio"
	"github.com/spf13/cobra"
)

var (
	pcount    int
	threshold int
	timeout   int64
	keynames  []string
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
	Run: func(cmd *cobra.Command, args []string) {
		// preParams, err := keygen.GeneratePreParams(time.Duration(timeout) * time.Minute)
		// if err != nil {
		// 	fmt.Println(err, preParams)
		// 	return
		// }

		if len(keynames) == 0 {
			for i := range pcount {
				keynames = append(keynames, fmt.Sprint(i))
			}
		}

		if len(keynames) != pcount {
			fmt.Println(len(keynames), pcount)
			fmt.Println("the numbers of key names provided is not equal to party count")
			return
		}

		pids := make(tss.UnSortedPartyIDs, 0)

		for i := 0; i < pcount; i++ {
			id, err := uuid.New().MarshalBinary()
			if err != nil {
				fmt.Println(err)
			}
			pids = append(pids,
				tss.NewPartyID(fmt.Sprint(i), fmt.Sprintf("p[%d]", i), big.NewInt(0).SetBytes(id)))
		}

		sorted := tss.SortPartyIDs(pids)
		ctx := tss.NewPeerContext(sorted)

		// Select an elliptic curve
		// use ECDSA
		curve := tss.S256()
		// or use EdDSA
		// curve := tss.Edwards()

		wg := &sync.WaitGroup{}
		pks := make([]*keygen.LocalPartySaveData, len(pids))
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
									fmt.Println(err)
								}
							}()
						}
					case data := <-endCh:
						pks[i] = data
						return
					case <-time.Tick(time.Minute):
						fmt.Printf("%v runing state %v\n", party.PartyID(), party.Running())
					case <-time.After(time.Duration(timeout) * time.Minute):
						fmt.Printf("goroutine[%d]	timeout\n", i)
						return
					}
				}
			}(i)
			go func() {
				defer wg.Done()
				err := party.Start()
				if err != nil {
					fmt.Println(err)
					return
				}
			}()
		}

		wg.Wait()

		for i, pk := range pks {
			if err := fileio.Save(keynames[i], pk); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("[%d]\n%v\n", i, pk)
		}
	},
}

func init() {
	rootCmd.AddCommand(keygenCmd)

	keygenCmd.Flags().IntVarP(&pcount, "party-count", "c", 2, "defines the number of parties will be generated")
	keygenCmd.Flags().IntVarP(&threshold, "threshold", "t", 1, "defines the threshold for signature verification")
	keygenCmd.Flags().Int64VarP(&timeout, "timeout", "o", 1, "defines the minutes of timeout for preparams")
	keygenCmd.Flags().StringArrayVarP(&keynames, "key-names", "n", []string{}, "defines the name of keys")
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
