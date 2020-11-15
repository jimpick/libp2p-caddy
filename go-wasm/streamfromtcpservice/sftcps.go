package streamfromtcpservice

import (
	"fmt"
	"io/ioutil"
	"syscall/js"

	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-daemon/p2pclient"
	multiaddr "github.com/multiformats/go-multiaddr"
)

type StreamFromTCPService struct {
	Node *p2pclient.Client
}

// The following functions implement a window.ping() entrypoint callable from
// JS that returns a promise.

// See: https://withblue.ink/2020/10/03/go-webassembly-http-requests-and-promises.html

func (sfwss *StreamFromTCPService) GetStream(this js.Value, param []js.Value) interface{} {
	maddr := param[0].String()
	println("Go maddr: ", maddr)

	streamFromTCPHandler := func(this js.Value, args []js.Value) interface{} {
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

			if err := sfwss.Node.Connect(peer.ID, peer.Addrs); err != nil {
				fmt.Printf("Connect error %v\n", err)
				reject.Invoke(js.ValueOf("Connect error"))
				return
			}
			_, s, err := sfwss.Node.NewStream(peer.ID, []string{"/cats"})
			if err != nil {
				fmt.Println("huh, this should have worked: ", err)
				return
			}
			out, err := ioutil.ReadAll(s)
			fmt.Println(string(out))
			resolve.Invoke(js.ValueOf(string(out)))
		}()
		return nil
	}

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(js.FuncOf(streamFromTCPHandler))
}
