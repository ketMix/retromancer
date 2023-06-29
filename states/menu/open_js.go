//go:build wasm

package menu

import (
	"syscall/js"
)

func OpenURL(path string) error {
	js.Global().Get("window").Call("open", path, "_blank")
	return nil
}
