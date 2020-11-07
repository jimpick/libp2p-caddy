package api

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	peer "github.com/libp2p/go-libp2p-peer"
)

// QueryAskAPI implements API passing calls to user-provided function values.
type QueryAskAPI interface {
	// ClientQueryAsk returns a signed StorageAsk from the specified miner.
	ClientQueryAsk(ctx context.Context, p peer.ID, miner address.Address) (*storagemarket.StorageAsk, error)
}

/*
type QueryAskAPI struct {
}

func (api *QueryAskAPI) ClientQueryAsk(peerID string, address string) string {
	return "Hello, " + address + " at " + peerID
}
*/
