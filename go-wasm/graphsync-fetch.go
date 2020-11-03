package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"syscall/js"

	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
	uio "github.com/ipfs/go-unixfs/io"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func graphSyncFetch(this js.Value, param []js.Value) interface{} {
	maddr := param[0].String()
	cidArg := param[1].String()
	fmt.Printf("Go maddr %v\n", maddr)
	fmt.Printf("Go cid: %v\n", cidArg)

	graphSyncFetchHandler := func(this js.Value, args []js.Value) interface{} {
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

			target, err := cid.Decode(cidArg)
			if err != nil {
				fmt.Printf("failed to decode CID '%q': %s\n", cidArg, err)
				reject.Invoke(js.ValueOf("CID Decode error"))
			}
			fmt.Printf("Go target: %v\n", target)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := node.Connect(ctx, *peer); err != nil {
				fmt.Printf("Connect error %v\n", err)
				reject.Invoke(js.ValueOf("Connect error"))
				return
			}

			bs := blockstore.NewBlockstore(dssync.MutexWrap(datastore.NewMapDatastore()))
			gs, err := newGraphsync(ctx, node, bs)
			if err != nil {
				fmt.Printf("newGraphsync error %v\n", err)
				reject.Invoke(js.ValueOf("newGraphsync error"))
				return
			}

			err = fetch(ctx, gs, peer.ID, target)
			if err != nil {
				fmt.Printf("fetch error %v\n", err)
				reject.Invoke(js.ValueOf("fetch error"))
				return
			}

			dag := merkledag.NewDAGService(blockservice.New(bs, offline.Exchange(bs)))
			root, err := dag.Get(ctx, target)
			if err != nil {
				fmt.Printf("get.Get error %v\n", err)
				reject.Invoke(js.ValueOf("get.Get error"))
				return
			}

			reader, err := uio.NewDagReader(ctx, root, dag)
			if err != nil {
				fmt.Printf("NewDagReader error %v\n", err)
				reject.Invoke(js.ValueOf("NewDagReader error"))
				return
			}

			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, reader)
			if err != nil {
				fmt.Printf("io.Copy error %v\n", err)
				reject.Invoke(js.ValueOf("io.Copy error"))
				return
			}

			resolve.Invoke(js.ValueOf(buf.String()))
		}()
		return nil
	}

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(js.FuncOf(graphSyncFetchHandler))
}
