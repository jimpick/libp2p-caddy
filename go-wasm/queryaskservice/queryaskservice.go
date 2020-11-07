package queryaskservice

import (
	"context"
	"syscall/js"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/jimpick/libp2p-caddy/go-wasm/queryaskservice/api"
	"github.com/jimpick/libp2p-caddy/go-wasm/queryaskservice/node"
)

var jsHandler js.Value

func Start() {
	var queryAskAPI api.QueryAskAPI

	ctx := context.Background()

	_, err := node.New(ctx,
		node.QueryAskAPI(&queryAskAPI),
		// node.Repo(r),
		node.Online(),
		// node.Override(new(moduleapi.ChainModuleAPI), nodeAPI),
		// node.Override(new(moduleapi.StateModuleAPI), nodeAPI),
	)
	if err != nil {
		panic(err)
	}

	rpcServer := jsonrpc.NewJSServer("connectQueryAskService")
	rpcServer.Register("Filecoin", &queryAskAPI)
}
