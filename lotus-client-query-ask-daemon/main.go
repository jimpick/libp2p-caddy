package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api"
	lcli "github.com/filecoin-project/lotus/cli/cmd"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/jimpick/libp2p-caddy/lotus-client-query-ask-daemon/api"
	"github.com/jimpick/libp2p-caddy/lotus-client-query-ask-daemon/node"
)

const flagQueryAskRepo = "query-ask-repo"

const listenAddr = "127.0.0.1:1238"

var daemonCmd = &cli.Command{
	Name:  "daemon",
	Usage: "run client query ask api daemon",
	Action: func(cctx *cli.Context) error {
		var queryAskAPI api.FullNodeStruct

		nodeAPI, ncloser, err := lcli.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer ncloser()
		ctx := lcli.DaemonContext(cctx)

		r, err := repo.NewFS(cctx.String(flagQueryAskRepo))
		if err != nil {
			return xerrors.Errorf("opening fs repo: %w", err)
		}

		// Re-use repo.Worker type as it has no config defaults
		if err := r.Init(repo.Worker); err != nil && err != repo.ErrRepoExists {
			return xerrors.Errorf("repo init error: %w", err)
		}

		_, err = node.New(ctx,
			node.QueryAskAPI(&queryAskAPI),
			node.Repo(r),
			node.Online(),
			node.Override(new(api.FullNode), nodeAPI),
		)
		if err != nil {
			return xerrors.Errorf("initializing node: %w", err)
		}
		rpcServer := jsonrpc.NewServer()
		rpcServer.Register("Filecoin", queryAskAPI)

		http.Handle("/rpc/v0", rpcServer)

		fmt.Printf("Listening on http://%s\n", listenAddr)
		return http.ListenAndServe(listenAddr, nil)
	},
}

func main() {
	app := &cli.App{
		Name: "lotus-client-query-ask-api-daemon",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "repo",
				EnvVars: []string{"LOTUS_PATH"},
				Hidden:  true,
				Value:   "~/.lotus", // TODO: Consider XDG_DATA_HOME
			},
			&cli.StringFlag{
				Name:    flagQueryAskRepo,
				EnvVars: []string{"LOTUS_QUERY_ASK_PATH"},
				Value:   "~/.lotus-client-query-ask", // TODO: Consider XDG_DATA_HOME
			},
		},
		Commands: []*cli.Command{
			daemonCmd,
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
