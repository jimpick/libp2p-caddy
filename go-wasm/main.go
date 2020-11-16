package main

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/jimpick/libp2p-caddy/go-wasm/helloservice"
	"github.com/jimpick/libp2p-caddy/go-wasm/pingservice"
	"github.com/jimpick/libp2p-caddy/go-wasm/queryaskservice"
	"github.com/jimpick/libp2p-caddy/go-wasm/streamfromtcpservice"
	"github.com/jimpick/libp2p-caddy/go-wasm/streamfromwsservice"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-daemon/p2pclient"
	multiaddr "github.com/multiformats/go-multiaddr"

	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	ws "github.com/libp2p/go-ws-transport"
)

var (
	node        host.Host
	pingService *ping.PingService
)

func main() {
	ctx := context.Background()

	// local libp2p node for wss
	var err error
	node, err = libp2p.New(ctx,
		libp2p.Transport(ws.New),
		libp2p.Ping(false),
	)
	if err != nil {
		panic(err)
	}

	// remote libp2p node for non-wss
	controlMaddr, _ := multiaddr.NewMultiaddr("/dns4/libp2p-caddy-p2pd.localhost/tcp/9059/wss")
	listenMaddr, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/0")
	p2pclientNode, err := p2pclient.NewClient(controlMaddr, listenMaddr)
	fmt.Printf("Jim p2pclientNode %v\n", p2pclientNode)
	nodeID, nodeAddrs, err := p2pclientNode.Identify()
	peerInfo := peerstore.AddrInfo{
		ID:    nodeID,
		Addrs: nodeAddrs,
	}
	fmt.Printf("Jim peerInfo %v\n", peerInfo)
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("p2pclient->p2pd node address:", addrs[0])

	// configure our own ping protocol
	pingService = &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	// APIs
	helloservice.Start()
	queryaskservice.Start()

	ping := pingservice.PingService{
		Node:        &node,
		PingService: pingService,
	}

	sfwss := streamfromwsservice.StreamFromWSService{
		Node: &node,
	}

	sftcps := streamfromtcpservice.StreamFromTCPService{
		Node: p2pclientNode,
	}

	js.Global().Set("ping", js.FuncOf(ping.PingNode))
	js.Global().Set("graphSyncFetch", js.FuncOf(graphSyncFetch))
	js.Global().Set("chainHead", js.FuncOf(chainHead))
	js.Global().Set("streamFromWs", js.FuncOf(sfwss.GetStream))
	js.Global().Set("streamFromTcp", js.FuncOf(sftcps.GetStream))

	println("WASM Go Initialized 2")

	c := make(chan struct{}, 0)
	<-c // wait forever
}
