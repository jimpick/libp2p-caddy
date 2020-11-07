package queryaskservice

import (
	"syscall/js"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/jimpick/libp2p-caddy/go-wasm/queryaskservice/api"
)

var jsHandler js.Value

func Start() {
	var queryAskAPI api.QueryAskAPI = api.QueryAskAPI{}
	rpcServer := jsonrpc.NewJSServer("connectQueryAskService")
	rpcServer.Register("Filecoin", &queryAskAPI)
}
