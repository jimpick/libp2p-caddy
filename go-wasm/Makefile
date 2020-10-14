main.wasm:
	GOOS=js GOARCH=wasm go build -o main.wasm

nodemon:
	GOOS=js GOARCH=wasm nodemon --ext go --watch main.go --watch ../../go-libp2p --watch ../../go-libp2p-core --watch ../../go-libp2p-swarm --exec go build -o main.wasm --signal SIGTERM

clean:
	rm -f *.wasm
