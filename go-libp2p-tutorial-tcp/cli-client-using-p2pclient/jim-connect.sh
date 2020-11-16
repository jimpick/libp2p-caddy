#! /bin/bash

. ../../../libp2p-caddy/.env.tcp
go run . /ip4/127.0.0.1/tcp/2020/p2p/$PEER_ID_TCP

