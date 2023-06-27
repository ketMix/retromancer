//go:build wasm

package resources

import (
	"syscall/js"
)

func ReadClipboard() string {
	wait := make(chan string)
	if js.Global().Get("navigator").Get("clipboard").Get("readText").IsUndefined() {
		js.Global().Get("alert").Invoke("Clipboard API is not supported in this browser.")
		return ""
	}
	js.Global().Get("navigator").Get("clipboard").Call("readText").Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		wait <- args[0].String()
		return nil
	}))
	return <-wait
}

func WriteClipboard(text string) {
	if js.Global().Get("navigator").Get("clipboard").Get("readText").IsUndefined() {
		js.Global().Get("alert").Invoke("Clipboard API is not supported in this browser.")
		return
	}
	wait := make(chan interface{})
	js.Global().Get("navigator").Get("clipboard").Call("writeText", text).Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		wait <- nil
		return nil
	}))
	<-wait
	return
}
