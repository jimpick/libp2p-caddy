package queryaskservice

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/node/modules/lp2p"
	"github.com/filecoin-project/lotus/node/modules/moduleapi"
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

	nodeAPI := &apistruct.FullNodeStruct{}

	_, err = node.New(ctx,
		node.QueryAskAPI(&queryAskAPI),
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
		fmt.Println("Jim ConnectCallback", environment)
	})
	fmt.Println("Jim1")
	rpcServer := jsonrpc.NewJSServer("connectQueryAskService", cbOpt)
	fmt.Println("Jim2")
	rpcServer.Register("Filecoin", queryAskAPI)
}
