package cmd

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/bnb-chain/tss-lib/v2/common"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/signing"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	keyname string
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
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
		for i, p := range pids {
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

		// signing
		threshold = 1
		timeout = 2
		message := big.NewInt(0).SetBytes([]byte("hello world!"))
		outCh := make(chan tss.Message, len(pids))
		endCh := make(chan *common.SignatureData, len(pids))
		signps := make([]tss.Party, len(pids))
		for i, key := range pks {
			params := tss.NewParameters(curve, ctx, pids[i], len(pids), threshold)
			party := signing.NewLocalParty(message, params, *key, outCh, endCh)
			signps[i] = party
			go func() {
				err := party.Start()
				if err != nil {
					fmt.Println(err)
				}
			}()
		}

		// go func() {
		// 	for data := range sendCh {
		// 		fmt.Println(&data)
		// 	}
		// }()

		for {
			select {
			case msg := <-outCh:
				dest := msg.GetTo()
				if dest == nil {
					go func() {
						for _, p := range signps {
							if p.PartyID().Index == msg.GetFrom().Index {
								continue
							}
							fmt.Println("no dest")
							if err := SharedPartyUpdater(p, msg); err != nil {
								fmt.Println(err)
							}
						}
					}()
				} else {
					go func() {
						fmt.Println("dest")
						if err := SharedPartyUpdater(signps[dest[0].Index], msg); err != nil {
							fmt.Println(err)
						}
					}()
				}
			case data := <-endCh:
				fmt.Println("test")
				fmt.Println(data)
				return
			case <-time.After(30 * time.Second):
				fmt.Printf("signing timeout\n")
				fmt.Printf("%d, %d, %d", len(outCh), len(endCh), 1)
				return
				// case <-time.Tick(10 * time.Second):
				// 	for _, p := range paries {
				// 		fmt.Printf("%v runing state %v\n", p.PartyID().Id, p.Running())
				// 	}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	keygenCmd.Flags().StringVarP(&keyname, "key", "k", "", "specify the name of key signing the message, default using the first local key")
}
