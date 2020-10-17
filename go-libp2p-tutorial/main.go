package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	ws "github.com/libp2p/go-ws-transport"
)

func main() {
	// create a background context (i.e. one that never cancels)
	ctx := context.Background()

	// start a libp2p node with default settings
	node, err := libp2p.New(ctx,
		libp2p.Transport(ws.New),
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/2000/ws"),
	)
	if err != nil {
		panic(err)
	}

	// print the node's listening addresses
	fmt.Println("Listen addresses:", node.Addrs())

	// print the node's PeerInfo in multiaddr format
	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	fmt.Println("libp2p node address:", addrs[0])

	dotEnvContents := []byte("PEER_ID=" + node.ID().Pretty() + "\n")
	err = ioutil.WriteFile("../.env", dotEnvContents, 0644)
	if err != nil {
		panic(err)
	}
	// wait for a SIGINT or SIGTERM signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")

	// shut the node down
	if err := node.Close(); err != nil {
		panic(err)
	}
}
