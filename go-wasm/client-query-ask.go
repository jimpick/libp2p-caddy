package main

import (
	"syscall/js"
)

func clientQueryAsk(this js.Value, param []js.Value) interface{} {
	// jsHandler := param[0]
	println("Go ClientQueryAsk")

	clientQueryAskHandler := func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		// reject := args[1]

		go func() {
			/*
				var api FullNodeStruct
				closer, err := jsonrpc.NewJsMergeClient(context.Background(), jsHandler, "Filecoin", []interface{}{&api})
				if err != nil {
					fmt.Printf("connecting with lotus failed: %s\n", err)
					return
				}
				defer closer()

				// Now you can call any API you're interested in.
				tipset, err := api.ChainHead(context.Background())
				if err != nil {
					fmt.Printf("calling chain head: %s\n", err)
					return
				}
				fmt.Printf("Current chain head is: %s\n", tipset.String())
			*/
			resolve.Invoke(js.ValueOf(1234))
		}()
		return nil
	}

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(js.FuncOf(clientQueryAskHandler))
}
