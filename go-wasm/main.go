package main

import (
	"syscall/js"
)

func add(this js.Value, i []js.Value) interface{} {
	result := js.ValueOf(i[0].Int() + i[1].Int())
	println(result.String())
	return js.ValueOf(result)
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("add", js.FuncOf(add))

	println("WASM Go Initialized 2")
	<-c // wait forever
}
