package impl

import (
	"github.com/filecoin-project/lotus/node/impl/client"
	"github.com/jimpick/lotus-query-ask-api-daemon/api"
)

type QueryAskAPI struct {
	client.API
}

var _ api.QueryAskAPI = &QueryAskAPI{}
