//go:build !wasm

package resources

import "golang.design/x/clipboard"

func init() {
	if err := clipboard.Init(); err != nil {
		panic(err)
	}
}

func ReadClipboard() string {
	return string(clipboard.Read(clipboard.FmtText))
}

func WriteClipboard(text string) {
	clipboard.Write(clipboard.FmtText, []byte(text))
}
