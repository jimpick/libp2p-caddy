package queryaskservice

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/node/modules/lp2p"
	"github.com/filecoin-project/lotus/node/modules/moduleapi"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/jimpick/libp2p-caddy/go-wasm/queryaskservice/node"
	"github.com/jimpick/lotus-query-ask-api-daemon/api"
	"github.com/libp2p/go-libp2p-daemon/p2pclient"
)

var jsHandler js.Value

func Start(p2pclientNode *p2pclient.Client) {
	var queryAskAPI api.QueryAskAPI

	ctx := context.Background()

	r := repo.NewMemory(nil)

	nilRouting, err := lp2p.NilRouting(ctx)
	if err != nil {
		panic(err)
	}

	var fullNodeStruct = apistruct.FullNodeStruct{}
	var nodeAPI lotusapi.FullNode = &fullNodeStruct
	var closer jsonrpc.ClientCloser
	defer func() {
		if closer != nil {
			closer()
		}
	}()

	_, err = node.New(ctx,
		node.QueryAskAPI(&queryAskAPI),
		node.Repo(r),
		node.Online(),
		node.Override(new(lp2p.BaseIpfsRouting), nilRouting),
		node.Override(new(moduleapi.ChainModuleAPI), nodeAPI),
		node.Override(new(moduleapi.StateModuleAPI), nodeAPI),
		node.Override(new(*p2pclient.Client), p2pclientNode),
	)
	if err != nil {
		panic(err)
	}

	cbOpt := jsonrpc.WithConnectCallback(func(environment js.Value) {
		requestsForLotusHandler := environment.Get("requestsForLotusHandler")

		// closer, err := jsonrpc.NewJSMergeClient(context.Background(), requestsForLotusHandler, "Filecoin", []interface{}{&nodeAPI})
		closer, err = jsonrpc.NewJSMergeClient(context.Background(), requestsForLotusHandler, "Filecoin",
			[]interface{}{
				&fullNodeStruct.CommonStruct.Internal,
				&fullNodeStruct.Internal,
			})
		if err != nil {
			fmt.Printf("connecting with lotus failed: %s\n", err)
			panic(err)
		}
	})
	rpcServer := jsonrpc.NewJSServer("connectQueryAskService", cbOpt)
	rpcServer.Register("Filecoin", queryAskAPI)
}
