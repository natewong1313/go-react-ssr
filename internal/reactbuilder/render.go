package reactbuilder

import (
	"os/exec"
)

func RenderReactToHTML(serverRenderJSFilePath string) (string, error) {
	cmd := exec.Command("node", serverRenderJSFilePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
