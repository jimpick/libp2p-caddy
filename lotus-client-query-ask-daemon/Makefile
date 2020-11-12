GOFLAGS+=-tags=clientQueryAskAPI

lotus-client-query-ask-daemon:
	go build .

clean:
	rm -f lotus-client-query-ask-daemon

why:
	go mod why github.com/filecoin-project/filecoin-ffi

debug:
	GOLOG_LOG_LEVEL=debug go run . daemon 2>&1 | tee ~/tmp/debug9.log
