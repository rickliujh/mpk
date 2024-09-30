package fileio

import "github.com/bnb-chain/tss-lib/v2/tss"

type (
	Meta struct {
		Threshold int
		Peers     []*tss.PartyID
	}
)
