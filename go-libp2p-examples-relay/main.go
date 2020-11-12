package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	multiaddr "github.com/multiformats/go-multiaddr"

	circuit "github.com/libp2p/go-libp2p-circuit"
)

func main() {

	var priv crypto.PrivKey
	privBytes, err := ioutil.ReadFile("privkey")
	if err != nil {
		fmt.Println("Generating new keypair")
		// Set your own keypair
		priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
		/*
			priv, _, err = crypto.GenerateKeyPair(
				crypto.Ed25519, // Select your key type. Ed25519 are nice short
				-1,             // Select key length when possible (i.e. RSA).
			)
		*/
		if err != nil {
			panic(err)
		}
		privBytes, err := crypto.MarshalPrivateKey(priv)
		if err != nil {
			panic(err)
		}
		privBase64 := base64.StdEncoding.EncodeToString(privBytes)
		err = ioutil.WriteFile("privkey", []byte(privBase64), 0644)
		if err != nil {
			panic(err)
		}
	} else {
		privBytesDecoded, err := base64.StdEncoding.DecodeString(string(privBytes))
		if err != nil {
			panic(err)
		}
		priv, err = crypto.UnmarshalPrivateKey(privBytesDecoded)
		if err != nil {
			panic(err)
		}
	}

	// Create three libp2p hosts, enable relay client capabilities on all
	// of them.

	// Tell the host to monitor for relays.
	//h1, err := libp2p.New(context.Background(), libp2p.EnableRelay(circuit.OptDiscovery))
	/*
		h1, err := libp2p.New(context.Background())
		if err != nil {
			panic(err)
		}
	*/

	// Tell the host to relay connections for other peers (The ability to *use*
	// a relay vs the ability to *be* a relay)
	h2, err := libp2p.New(context.Background(),
		libp2p.Identity(priv),
		libp2p.EnableRelay(circuit.OptHop),
		// libp2p.ChainOptions(
		// 	libp2p.Transport(ws.New),
		libp2p.DefaultTransports,
		// ),
		libp2p.ListenAddrStrings(
			"/ip4/127.0.0.1/tcp/2010/ws",
			"/ip4/127.0.0.1/tcp/2011",
		),
		libp2p.Ping(false),
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr.NewConnManager(
			100,            // Lowwater
			400,            // HighWater,
			10*time.Minute, // GracePeriod
		)),
	)
	if err != nil {
		panic(err)
	}

	// configure our own ping protocol
	h2PingService := &ping.PingService{Host: h2}
	h2.SetStreamHandler(ping.ID, h2PingService.PingHandler)

	// print the node's listening addresses
	fmt.Println("Relay listen addresses:", h2.Addrs())
	// print the node's PeerInfo in multiaddr format
	h2peerInfo := peerstore.AddrInfo{
		ID:    h2.ID(),
		Addrs: h2.Addrs(),
	}
	h2addrs, err := peerstore.AddrInfoToP2pAddrs(&h2peerInfo)
	fmt.Println("libp2p relay node address:", h2addrs[0])
	fmt.Println("libp2p relay node address:", h2addrs[0])

	// Zero out the listen addresses for the host, so it can only communicate
	// via p2p-circuit for our example
	h3, err := libp2p.New(context.Background(),
		libp2p.ListenAddrs(),
		libp2p.EnableRelay(),
		libp2p.Ping(false))
	if err != nil {
		panic(err)
	}

	// configure our own ping protocol
	h3PingService := &ping.PingService{Host: h3}
	h3.SetStreamHandler(ping.ID, h3PingService.PingHandler)

	// print the node's listening addresses
	fmt.Println("Cat server listen addresses:", h3.Addrs())

	// print the node's PeerInfo in multiaddr format
	h3peerInfo := peerstore.AddrInfo{
		ID:    h3.ID(),
		Addrs: h3.Addrs(),
	}
	h3addrs, err := peerstore.AddrInfoToP2pAddrs(&h3peerInfo)
	fmt.Println("libp2p cat server node address:", h3addrs[0])

	h2info := peer.AddrInfo{
		ID:    h2.ID(),
		Addrs: h2.Addrs(),
	}

	// Connect both h1 and h3 to h2, but not to each other
	/*
		if err := h1.Connect(context.Background(), h2info); err != nil {
			panic(err)
		}
	*/
	if err := h3.Connect(context.Background(), h2info); err != nil {
		panic(err)
	}

	// Now, to test things, let's set up a protocol handler on h3
	h3.SetStreamHandler("/cats", func(s network.Stream) {
		fmt.Println("Meow! It worked!")
		s.Close()
	})

	/*
		_, err = h1.NewStream(context.Background(), h3.ID(), "/cats")
		if err == nil {
			fmt.Println("Didnt actually expect to get a stream here. What happened?")
			return
		}
		fmt.Println("Okay, no connection from h1 to h3: ", err)
		fmt.Println("Just as we suspected")
	*/

	// Creates a relay address
	relayaddr, err := multiaddr.NewMultiaddr("/p2p-circuit/ipfs/" + h3.ID().Pretty())
	if err != nil {
		panic(err)
	}
	fmt.Println("Relay address", relayaddr.String())

	/*
		// Since we just tried and failed to dial, the dialer system will, by default
		// prevent us from redialing again so quickly. Since we know what we're doing, we
		// can use this ugly hack (it's on our TODO list to make it a little cleaner)
		// to tell the dialer "no, its okay, let's try this again"
		h1.Network().(*swarm.Swarm).Backoff().Clear(h3.ID())
	*/

	/*
		h3relayInfo := peer.AddrInfo{
			ID:    h3.ID(),
			Addrs: []multiaddr.Multiaddr{relayaddr},
		}
		if err := h1.Connect(context.Background(), h3relayInfo); err != nil {
			panic(err)
		}
	*/

	/*
		// Woohoo! we're connected!
		s, err := h1.NewStream(context.Background(), h3.ID(), "/cats")
		if err != nil {
			fmt.Println("huh, this should have worked: ", err)
			return
		}

		s.Read(make([]byte, 1)) // block until the handler closes the stream
	*/
	maddr := "/ip4/80.82.17.10/tcp/9999/p2p/12D3KooWDMpcct12Vb6jPXwjvLQHA2hoP8XKGbUZ2tpue1ydoZUm"
	fmt.Println("Dialing Lotus", maddr)
	addr, err := multiaddr.NewMultiaddr(maddr)
	if err != nil {
		panic(err)
	}
	lotusPeer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}
	connectFunc := func() {
		fmt.Printf("Connecting\n")
		if err := h2.Connect(context.Background(), *lotusPeer); err != nil {
			fmt.Printf("Error connecting %v\n", err)
		}
	}
	connectFunc()
	h2.ConnManager().Protect(lotusPeer.ID, "protected")
	fmt.Println("sending 500 ping messages to", lotusPeer.ID)
	ch := h2PingService.Ping(context.Background(), lotusPeer.ID)
	for i := 0; i < 500; i++ {
		res := <-ch
		fmt.Println("pinged", addr, "in", res.RTT, "number", i)
		if res.RTT == 0 {
			os.Exit(0)
			// connectFunc()
		}
		time.Sleep(15 * time.Second)
	}

	// wait for a SIGINT or SIGTERM signal
	pauseCh := make(chan os.Signal, 1)
	signal.Notify(pauseCh, syscall.SIGINT, syscall.SIGTERM)
	<-pauseCh
	fmt.Println("Received signal, shutting down...")
}
