package main

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

type FullNodeStruct struct {
	ChainHead func(context.Context) (*types.TipSet, error) `perm:"read"`
}

type FileRef struct {
	Path  string
	IsCAR bool
}

type RetrievalOrder struct {
	// TODO: make this less unixfs specific
	Root  cid.Cid
	Piece *cid.Cid
	Size  uint64
	// TODO: support offset
	Total                   types.BigInt
	UnsealPrice             types.BigInt
	PaymentInterval         uint64
	PaymentIntervalIncrease uint64
	Client                  address.Address
	Miner                   address.Address
	MinerPeer               retrievalmarket.RetrievalPeer
}
