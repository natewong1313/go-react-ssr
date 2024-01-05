package reactbuilder

import (
	"fmt"
	"os"
	"strings"

	"github.com/buger/jsonparser"
	esbuildApi "github.com/evanw/esbuild/pkg/api"
	"github.com/natewong1313/go-react-ssr/internal/utils"
)

type BuildResult struct {
	JS           string
	CSS          string
	Dependencies []string
}

func BuildServer(buildContents, frontendDir, assetRoute string) (BuildResult, error) {
	opts := esbuildApi.BuildOptions{
		Stdin: &esbuildApi.StdinOptions{
			Contents:   buildContents,
			Loader:     esbuildApi.LoaderTSX,
			ResolveDir: frontendDir,
		},
		Platform:          esbuildApi.PlatformNode,
		Bundle:            true,
		Write:             false,
		Outdir:            "/",
		Metafile:          false,
		AssetNames:        fmt.Sprintf("%s/[name]", strings.TrimPrefix(assetRoute, "/")),
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Loader: map[string]esbuildApi.Loader{ // for loading images properly
			".png":   esbuildApi.LoaderFile,
			".svg":   esbuildApi.LoaderFile,
			".jpg":   esbuildApi.LoaderFile,
			".jpeg":  esbuildApi.LoaderFile,
			".gif":   esbuildApi.LoaderFile,
			".bmp":   esbuildApi.LoaderFile,
			".woff2": esbuildApi.LoaderFile,
			".woff":  esbuildApi.LoaderFile,
			".ttf":   esbuildApi.LoaderFile,
			".eot":   esbuildApi.LoaderFile,
		},
	}
	return build(opts, false)
}

func BuildClient(buildContents, frontendDir, assetRoute string) (BuildResult, error) {
	opts := esbuildApi.BuildOptions{
		Stdin: &esbuildApi.StdinOptions{
			Contents:   buildContents,
			Loader:     esbuildApi.LoaderTSX,
			ResolveDir: frontendDir,
		},
		Bundle:            true,
		Write:             false,
		Outdir:            "/",
		Metafile:          true,
		AssetNames:        fmt.Sprintf("%s/[name]", strings.TrimPrefix(assetRoute, "/")),
		MinifyWhitespace:  os.Getenv("APP_ENV") == "production",
		MinifyIdentifiers: os.Getenv("APP_ENV") == "production",
		MinifySyntax:      os.Getenv("APP_ENV") == "production",
		Loader: map[string]esbuildApi.Loader{ // for loading images properly
			".png":   esbuildApi.LoaderFile,
			".svg":   esbuildApi.LoaderFile,
			".jpg":   esbuildApi.LoaderFile,
			".jpeg":  esbuildApi.LoaderFile,
			".gif":   esbuildApi.LoaderFile,
			".bmp":   esbuildApi.LoaderFile,
			".woff2": esbuildApi.LoaderFile,
			".woff":  esbuildApi.LoaderFile,
			".ttf":   esbuildApi.LoaderFile,
			".eot":   esbuildApi.LoaderFile,
		},
	}
	return build(opts, true)
}

func build(buildOptions esbuildApi.BuildOptions, isClient bool) (BuildResult, error) {
	result := esbuildApi.Build(buildOptions)
	if len(result.Errors) > 0 {
		fileLocation := "unknown"
		lineNum := "unknown"
		if result.Errors[0].Location != nil {
			fileLocation = result.Errors[0].Location.File
			lineNum = result.Errors[0].Location.LineText
		}
		return BuildResult{}, fmt.Errorf("%s <br>in %s <br>at %s", result.Errors[0].Text, fileLocation, lineNum)
	}

	var br BuildResult
	for _, file := range result.OutputFiles {
		if strings.HasSuffix(file.Path, "stdin.js") {
			br.JS = string(file.Contents)
		} else if strings.HasSuffix(file.Path, "stdin.css") {
			br.CSS = string(file.Contents)
		}
	}
	if isClient {
		br.Dependencies = getDependencyPathsFromMetafile(result.Metafile)
	}
	return br, nil
}

// getDependencyPathsFromMetafile parses dependencies from esbuild metafile and returns the paths of the dependencies
func getDependencyPathsFromMetafile(metafile string) []string {
	var dependencyPaths []string
	// Parse the metafile and get the paths of the dependencies
	// Ignore dependencies in node_modules
	err := jsonparser.ObjectEach([]byte(metafile), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		if !strings.Contains(string(key), "/node_modules/") {
			dependencyPaths = append(dependencyPaths, utils.GetFullFilePath(string(key)))
		}
		return nil
	}, "inputs")
	if err != nil {
		return nil
	}
	return dependencyPaths
}
