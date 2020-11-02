#! /bin/bash

export DATA='{"jsonrpc":"2.0","id":2,"method":"Filecoin.ClientQueryAsk","params":["12D3KooWDMpcct12Vb6jPXwjvLQHA2hoP8XKGbUZ2tpue1ydoZUm","f02620"]}'
echo $DATA | jq .

curl -X POST -H "Content-Type: application/json" \
       	--data "$DATA" \
       	'http://10.0.1.52:1234/rpc/v0' 

