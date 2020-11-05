package helloservice

import (
	"fmt"
	"syscall/js"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/jimpick/libp2p-caddy/go-wasm/helloservice/api"
)

var jsHandler js.Value

func Start() {
	fmt.Printf("Jim start 'listening' on Hello API here\n")
	js.Global().Set("connectHelloService", js.FuncOf(connectService))

	var helloAPI api.HelloAPI = api.HelloAPI{}
	rpcServer := jsonrpc.NewJSServer()
	rpcServer.Register("Hello", &helloAPI)
}

func connectService(this js.Value, param []js.Value) interface{} {
	jsHandler = param[0]
	return nil // Return function for submitting messages?
}
