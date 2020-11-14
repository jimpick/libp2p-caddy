package main

import (
	"fmt"
	"io/ioutil"
	"os"

	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-daemon/p2pclient"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func main() {
	// controlMaddr, _ := multiaddr.NewMultiaddr("/unix/tmp/p2pd.sock")
	controlMaddr, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/2030/ws")
	listenMaddr, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/0")

	node, err := p2pclient.NewClient(controlMaddr, listenMaddr)
	// start a libp2p node that listens on a random local TCP port,
	// but without running the built-in ping protocol
	/*
		node, err := libp2p.New(ctx,
			libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
			libp2p.Ping(false),
		)
	*/
	if err != nil {
		panic(err)
	}

	// configure our own ping protocol
	/*
		pingService := &ping.PingService{Host: node}
		node.SetStreamHandler(ping.ID, pingService.PingHandler)
	*/

	nodeID, nodeAddrs, err := node.Identify()

	// print the node's PeerInfo in multiaddr format
	peerInfo := peerstore.AddrInfo{
		ID:    nodeID,
		Addrs: nodeAddrs,
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("libp2p node address:", addrs[0])

	addr, err := multiaddr.NewMultiaddr(os.Args[1])
	if err != nil {
		panic(err)
	}
	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Dialing", addr.String())
	if err := node.Connect(peer.ID, peer.Addrs); err != nil {
		panic(err)
	}

	_, s, err := node.NewStream(peer.ID, []string{"/cats"})
	if err != nil {
		fmt.Println("huh, this should have worked: ", err)
		return
	}

	out, err := ioutil.ReadAll(s)
	fmt.Println(string(out))

	// shut the node down
	if err := node.Close(); err != nil {
		panic(err)
	}
}
