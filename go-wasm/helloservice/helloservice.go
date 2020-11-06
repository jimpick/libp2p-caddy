package helloservice

import (
	"syscall/js"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/jimpick/libp2p-caddy/go-wasm/helloservice/api"
)

var jsHandler js.Value

func Start() {
	var helloAPI api.HelloAPI = api.HelloAPI{}
	rpcServer := jsonrpc.NewJSServer("connectHelloService")
	rpcServer.Register("Filecoin", &helloAPI)
}
