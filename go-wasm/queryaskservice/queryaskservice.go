package queryaskservice

import (
	"context"
	"syscall/js"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/node/modules/lp2p"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/jimpick/libp2p-caddy/go-wasm/queryaskservice/api"
	"github.com/jimpick/libp2p-caddy/go-wasm/queryaskservice/node"
)

var jsHandler js.Value

func Start() {
	var queryAskAPI api.QueryAskAPI

	ctx := context.Background()

	r := repo.NewMemory(nil)

	nilRouting, err := lp2p.NilRouting(ctx)
	if err != nil {
		panic(err)
	}

	_, err = node.New(ctx,
		node.QueryAskAPI(&queryAskAPI),
		node.Repo(r),
		node.Online(),
		node.Override(new(lp2p.BaseIpfsRouting), nilRouting),
		// node.Override(new(moduleapi.ChainModuleAPI), nodeAPI),
		// node.Override(new(moduleapi.StateModuleAPI), nodeAPI),
	)
	if err != nil {
		panic(err)
	}

	rpcServer := jsonrpc.NewJSServer("connectQueryAskService")
	rpcServer.Register("Filecoin", &queryAskAPI)
}
