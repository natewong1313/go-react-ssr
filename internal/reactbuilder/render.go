package reactbuilder

import (
	"errors"
	"os/exec"
	"strings"
)

func RenderReactToHTML(serverRenderJSFilePath string) (string, error) {
	cmd := exec.Command("node", serverRenderJSFilePath)
	stdOut := new(strings.Builder)
	stdErr := new(strings.Builder)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	err := cmd.Run()
	if err != nil {
		return "", errors.New(stdErr.String())
	}
	return stdOut.String(), nil
}
