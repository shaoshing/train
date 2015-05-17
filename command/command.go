package trainCommand

import (
	"bytes"
	"os/exec"
)

func bash(bash string) (out string, err error) {
	cmd := exec.Command("sh", "-c", bash)
	var buf bytes.Buffer
	cmd.Stderr = &buf
	cmd.Stdout = &buf
	err = cmd.Run()
	out = buf.String()
	return
}
