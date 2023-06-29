package menu

import (
	"os/exec"
)

func OpenURL(path string) error {
	// This seems bad.
	cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}
