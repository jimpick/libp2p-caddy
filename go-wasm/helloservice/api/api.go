package api

type HelloAPI struct {
}

func (api *HelloAPI) Hello(name string) string {
	return "Hello, " + name
}
