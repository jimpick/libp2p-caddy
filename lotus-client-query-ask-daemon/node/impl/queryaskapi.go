package impl

import (
	"github.com/filecoin-project/lotus/node/impl/client"
	"github.com/jimpick/libp2p-caddy/lotus-client-query-ask-daemon/api"
)

type QueryAskAPI struct {
	client.API
}

var _ api.QueryAskAPI = &QueryAskAPI{}
