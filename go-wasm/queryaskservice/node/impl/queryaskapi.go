package impl

import (
	"github.com/filecoin-project/lotus/node/impl/client"
	"github.com/jimpick/libp2p-caddy/go-wasm/queryaskservice/api"
)

type QueryAskAPI struct {
	client.API
}

var _ api.QueryAskAPI = &QueryAskAPI{}
