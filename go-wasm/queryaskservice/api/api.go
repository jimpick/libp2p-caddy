package api

type QueryAskAPI struct {
}

func (api *QueryAskAPI) ClientQueryAsk(peerID string, address string) string {
	return "Hello, " + address + " at " + peerID
}
