package fileio

import (
	"crypto/ecdsa"

	"github.com/bnb-chain/tss-lib/v2/tss"
)

type (
	VaildFile interface {
		*Meta | *ecdsa.PublicKey
	}

	Meta struct {
		Threshold int
		Peers     []*tss.PartyID
	}
)
