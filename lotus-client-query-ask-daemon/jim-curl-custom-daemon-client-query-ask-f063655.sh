#! /bin/bash

export DATA='{"jsonrpc":"2.0","id":2,"method":"Filecoin.ClientQueryAsk","params":["12D3KooWEUS7VnaRrHF24GTWVGYtcEsmr3jsnNLcsEwPU7rDgjf5","f063655"]}'
echo $DATA | jq .

curl -X POST -H "Content-Type: application/json" \
       	--data "$DATA" \
       	'http://127.0.0.1:1238/rpc/v0' 

