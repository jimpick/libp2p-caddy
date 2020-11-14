#! /bin/bash

#go run . -hostAddrs /ip4/127.0.0.1/tcp/2071
go run . -listen /ip4/0.0.0.0/tcp/2030/ws -hostAddrs /ip4/127.0.0.1/tcp/2040
