package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/apistruct"
)

func main() {
	// Jimpick - nuc
	url := "wss://lotus.jimpick.com/spacerace_api/0/node/rpc/v0"
	authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.tOXKfED8C-eSNrSwSXtoIR1bGk2znmVs7eRScu1bp40" // already shared publicly
	headers := http.Header{"Authorization": []string{"Bearer " + authToken}}

	// For https://api.chain.love/
	// url := "wss://api.chain.love/rpc/v0"
	// headers := http.Header{}

	var api apistruct.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(context.Background(), url, "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	defer closer()

	// Now you can call any API you're interested in.
	tipset, err := api.ChainHead(context.Background())
	if err != nil {
		log.Fatalf("calling chain head: %s", err)
	}
	fmt.Printf("Current chain head is: %s", tipset.String())
}
