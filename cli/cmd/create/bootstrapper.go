package create

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/natewong1313/go-react-ssr/go-ssr-cli/logger"
	"github.com/natewong1313/go-react-ssr/go-ssr-cli/utils"
	"github.com/rs/zerolog"
)

type Bootstrapper struct {
	Logger          zerolog.Logger
	ProjectDir      string
	GoModuleName    string
	FrontendDir     string
	WebFramework    string
	IsUsingTailwind bool
}

func (b *Bootstrapper) Start() {
	var wg sync.WaitGroup
	wg.Add(2)
	go b.setupFrontend(&wg)
	go b.setupBackend(&wg)
	wg.Wait()

}

func (b *Bootstrapper) setupFrontend(wg *sync.WaitGroup) {
	logger.L.Info().Msg("Setting up frontend")
	b.createFrontendFolder()
	b.createSrcFolder()
	b.createFileInFrontendFolder("package.json", PACKAGE_JSON)
	b.installNPMDependencies()
	b.createFilesInFrontendFolder()
	logger.L.Info().Msg("Frontend setup complete")
	wg.Done()
}

func (b *Bootstrapper) createFrontendFolder() {
	logger.L.Info().Msg("Creating /frontend folder")
	if err := os.MkdirAll(b.ProjectDir+"/frontend", 0777); err != nil {
		utils.HandleError(err)
	}
	b.FrontendDir = b.ProjectDir + "/frontend"
}

func (b *Bootstrapper) createSrcFolder() {
	logger.L.Info().Msg("Creating /frontend/src folder")
	if err := os.MkdirAll(b.FrontendDir+"/src", 0777); err != nil {
		utils.HandleError(err)
	}
}

func (b *Bootstrapper) installNPMDependencies() {
	logger.L.Info().Msg("Installing npm dependencies")
	args := []string{"install", "typescript", "--save-dev"}
	if b.IsUsingTailwind {
		args = append(args, "tailwindcss")
		args = append(args, "--save-dev")
	}
	cmd := exec.Command("npm", args...)
	cmd.Dir = b.FrontendDir
	err := cmd.Run()
	if err != nil {
		utils.HandleError(err)
	}
}

func (b *Bootstrapper) createFilesInFrontendFolder() {
	var wg sync.WaitGroup
	var filePaths = map[string]string{
		"src/index.tsx": REACT_FILE,
		"src/index.css": INDEX_CSS,
		"tsconfig.json": TSCONFIG,
	}
	if b.IsUsingTailwind {
		filePaths["tailwind.config.js"] = TAILWIND_CONFIG
		filePaths["src/main.css"] = TAILWIND_CSS
	}
	for fileName, contents := range filePaths {
		wg.Add(1)
		go func(fileName, contents string) {
			defer wg.Done()
			b.createFileInFrontendFolder(fileName, contents)
		}(fileName, contents)
	}
	wg.Wait()
}

func (b *Bootstrapper) createFileInFrontendFolder(fileName, contents string) {
	logger.L.Debug().Msgf("Creating %s file", fileName)
	f, err := os.Create(b.FrontendDir + "/" + fileName)
	if err != nil {
		utils.HandleError(err)
	}
	defer f.Close()

	_, err = f.WriteString(contents)
	if err != nil {
		utils.HandleError(err)
	}
	f.Sync()
}

func (b *Bootstrapper) setupBackend(wg *sync.WaitGroup) {
	projectFolderName := filepath.Base(b.ProjectDir)
	b.GoModuleName = "example.com/" + strings.Replace(projectFolderName, " ", "-", -1)
	logger.L.Debug().Msgf("Creating go module %s", b.GoModuleName)

	b.createFileInFolder("main.go", "")

	logger.L.Info().Msg("Backend setup complete")
	wg.Done()
}

func (b *Bootstrapper) createFileInFolder(fileName, contents string) {
	logger.L.Debug().Msgf("Creating %s file", fileName)
	f, err := os.Create(b.ProjectDir + "/" + fileName)
	if err != nil {
		utils.HandleError(err)
	}
	defer f.Close()

	_, err = f.WriteString(contents)
	if err != nil {
		utils.HandleError(err)
	}
	f.Sync()
}
