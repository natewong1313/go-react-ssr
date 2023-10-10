package update

import (
	"io"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/natewong1313/go-react-ssr/gossr-cli/utils"
)

func getLatestVersion() string {
	res, err := http.Get("https://proxy.golang.org/github.com/natewong1313/go-react-ssr/@latest")
	if err != nil {
		utils.HandleError(err)
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		utils.HandleError(err)
	}
	version, err := jsonparser.GetString(resBody, "Version")
	if err != nil {
		utils.HandleError(err)
	}
	return version
}
