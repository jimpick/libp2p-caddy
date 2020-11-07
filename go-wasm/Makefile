main.wasm:
	GOOS=js GOARCH=wasm go build -o main.wasm

nodemon:
	GOOS=js GOARCH=wasm nodemon --ext go --watch . --watch ../../go-libp2p --watch ../../go-libp2p-core --watch ../../go-libp2p-swarm --watch ../../go-ws-transport-0xproject-feat-wss-dialing --watch ../../go-jsonrpc-wasm --watch ../../lotus-wasm --exec go build -o main.wasm --signal SIGTERM

clean:
	rm -f *.wasm
