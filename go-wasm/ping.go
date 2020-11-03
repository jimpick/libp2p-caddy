package main

import (
	"context"
	"fmt"
	"syscall/js"

	peerstore "github.com/libp2p/go-libp2p-core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// The following functions implement a window.ping() entrypoint callable from
// JS that returns a promise.

// See: https://withblue.ink/2020/10/03/go-webassembly-http-requests-and-promises.html

func pingNode(this js.Value, param []js.Value) interface{} {
	maddr := param[0].String()
	println("Go maddr: ", maddr)

	pingHandler := func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			addr, err := multiaddr.NewMultiaddr(maddr)
			if err != nil {
				fmt.Printf("NewMultiaddr error %v\n", err)
				reject.Invoke(js.ValueOf("NewMultiaddr error"))
				return
			}

			peer, err := peerstore.AddrInfoFromP2pAddr(addr)
			if err != nil {
				fmt.Printf("AddInfoFromP2pAddr error %v\n", err)
				reject.Invoke(js.ValueOf("AddInfoFromP2pAddr error"))
				return
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := node.Connect(ctx, *peer); err != nil {
				fmt.Printf("Connect error %v\n", err)
				reject.Invoke(js.ValueOf("Connect error"))
				return
			}

			ch := pingService.Ping(ctx, peer.ID)
			res := <-ch
			fmt.Println("pinged", addr, "in", res.RTT)
			resolve.Invoke(js.ValueOf(res.RTT.Milliseconds()))
		}()
		return nil
	}

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(js.FuncOf(pingHandler))
}
