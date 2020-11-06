package api

type HelloAPI struct {
}

func (api *HelloAPI) HelloName(name string) string {
	return "Hello, " + name
}
