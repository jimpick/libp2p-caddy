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
	"github.com/jimpick/libp2p-caddy/go-wasm/queryaskservice/api"
	"github.com/jimpick/libp2p-caddy/go-wasm/queryaskservice/node"
)

var jsHandler js.Value

func Start() {
	var queryAskAPI api.QueryAskAPI = api.QueryAskAPI{}

	ctx := context.Background()

	r := repo.NewMemory(nil)

	nilRouting, err := lp2p.NilRouting(ctx)
	if err != nil {
		panic(err)
	}

	var nodeAPI lotusapi.FullNode
	var closer jsonrpc.ClientCloser
	nodeAPI = &apistruct.FullNodeStruct{}
	defer func() {
		if closer != nil {
			closer()
		}
	}()

	_, err = node.New(ctx,
		// node.QueryAskAPI(&queryAskAPI),
		node.Repo(r),
		node.Online(),
		node.Override(new(lp2p.BaseIpfsRouting), nilRouting),
		node.Override(new(moduleapi.ChainModuleAPI), nodeAPI),
		node.Override(new(moduleapi.StateModuleAPI), nodeAPI),
	)
	if err != nil {
		panic(err)
	}

	cbOpt := jsonrpc.WithConnectCallback(func(environment js.Value) {
		fmt.Println("Jim ConnectCallback")
		requestsForLotusHandler := environment.Get("requestsForLotusHandler")
		jim := environment.Get("jim")
		fmt.Println("JimX", jim)

		var res apistruct.FullNodeStruct
		// closer, err := jsonrpc.NewJSMergeClient(context.Background(), requestsForLotusHandler, "Filecoin", []interface{}{&nodeAPI})
		closer, err = jsonrpc.NewJSMergeClient(context.Background(), requestsForLotusHandler, "Filecoin",
			[]interface{}{
				&res.CommonStruct.Internal,
				&res.Internal,
			})
		if err != nil {
			fmt.Printf("connecting with lotus failed: %s\n", err)
			panic(err)
		}
	})
	fmt.Println("Jim1")
	rpcServer := jsonrpc.NewJSServer("connectQueryAskService", cbOpt)
	fmt.Println("Jim2")
	rpcServer.Register("Filecoin", &queryAskAPI)
}
