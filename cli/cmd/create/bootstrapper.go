package create

import (
	"os"
	"os/exec"
	"strings"
	"sync"

	cp "github.com/otiai10/copy"

	"github.com/natewong1313/go-react-ssr/cli/logger"
	"github.com/natewong1313/go-react-ssr/cli/utils"
	"github.com/rs/zerolog"
)

type Bootstrapper struct {
	Logger          zerolog.Logger
	TempDirPath     string
	ProjectDir      string
	GoModuleName    string
	FrontendDir     string
	WebFramework    string
	IsUsingTailwind bool
}

func (b *Bootstrapper) Start() {
	b.TempDirPath = createTempDir()
	b.cloneRepo()
	b.moveGoFiles()
	var wg sync.WaitGroup
	wg.Add(2)
	go b.setupFrontend(&wg)
	go b.setupBackend(&wg)
	wg.Wait()
	logger.L.Info().Msg("Project setup complete! 🎉")
}

func (b *Bootstrapper) cloneRepo() {
	logger.L.Info().Msg("Cloning example repository")
	cmd := exec.Command("git", "clone", "-b", "go-ssr-cli", "https://github.com/natewong1313/go-react-ssr.git")
	cmd.Dir = b.TempDirPath
	err := cmd.Run()
	if err != nil {
		utils.HandleError(err)
	}
}

func (b *Bootstrapper) moveGoFiles() {
	logger.L.Info().Msg("Setting up Go files")
	err := cp.Copy(b.TempDirPath+"/go-react-ssr/examples/"+strings.ToLower(b.WebFramework), b.ProjectDir)
	if err != nil {
		utils.HandleError(err)
	}
}

func (b *Bootstrapper) setupFrontend(wg *sync.WaitGroup) {
	b.createFrontendFolder()
	b.installNPMDependencies()
	wg.Done()
}

func (b *Bootstrapper) setupBackend(wg *sync.WaitGroup) {
	b.updateGoModules()
	b.replaceImportsInGoFile()
	// projectFolderName := filepath.Base(b.ProjectDir)
	// b.GoModuleName = "example.com/" + strings.Replace(projectFolderName, " ", "-", -1)
	// logger.L.Debug().Msgf("Creating go module %s", b.GoModuleName)

	// b.createFileInFolder("main.go", "")

	// logger.L.Info().Msg("Backend setup complete")
	wg.Done()
}

func (b *Bootstrapper) createFrontendFolder() {
	logger.L.Info().Msg("Creating /frontend folder")
	frontendFolderFromGit := b.TempDirPath + "/go-react-ssr/examples/frontend"
	if b.IsUsingTailwind {
		frontendFolderFromGit = b.TempDirPath + "/go-react-ssr/examples/frontend-tailwind"
	}
	err := cp.Copy(frontendFolderFromGit, b.ProjectDir+"/frontend")
	if err != nil {
		utils.HandleError(err)
	}
}

func (b *Bootstrapper) installNPMDependencies() {
	logger.L.Info().Msg("Installing npm dependencies")
	// args := []string{"install", "typescript", "--save-dev"}
	// if b.IsUsingTailwind {
	// 	args = append(args, "tailwindcss")
	// 	args = append(args, "--save-dev")
	// }
	// cmd := exec.Command("npm", args...)
	cmd := exec.Command("npm", "install")
	cmd.Dir = b.ProjectDir + "/frontend"
	err := cmd.Run()
	if err != nil {
		utils.HandleError(err)
	}
}

func (b *Bootstrapper) updateGoModules() {
	logger.L.Info().Msg("Installing Go modules")
	cmd := exec.Command("go", "get", "-u", "github.com/natewong1313/go-react-ssr")
	cmd.Dir = b.ProjectDir
	err := cmd.Run()
	if err != nil {
		utils.HandleError(err)
	}
}

func (b *Bootstrapper) replaceImportsInGoFile() {
	logger.L.Info().Msg("Updating imports in main.go")
	read, err := os.ReadFile(b.ProjectDir + "/main.go")
	if err != nil {
		utils.HandleError(err)
	}
	newContents := strings.Replace(string(read), "../frontend", "./frontend", -1)
	err = os.WriteFile(b.ProjectDir+"/main.go", []byte(newContents), 0644)
	if err != nil {
		utils.HandleError(err)
	}
}
