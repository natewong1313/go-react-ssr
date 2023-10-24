package create

import (
	"os"
	"os/exec"
	"strings"
	"sync"

	cp "github.com/otiai10/copy"

	"github.com/natewong1313/go-react-ssr/gossr-cli/logger"
	"github.com/natewong1313/go-react-ssr/gossr-cli/utils"
)

type Bootstrapper struct {
	TempDirPath   string
	ProjectDir    string
	GoModuleName  string
	FrontendDir   string
	WebFramework  string
	StylingPlugin string
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
	cmd := exec.Command("git", "clone", "https://github.com/natewong1313/go-react-ssr.git")
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
	b.updateDockerFile()
	wg.Done()
}

func (b *Bootstrapper) createFrontendFolder() {
	logger.L.Info().Msg("Creating /frontend folder")
	frontendFolderFromGit := b.TempDirPath + "/go-react-ssr/examples/frontend"
	if b.StylingPlugin == "Tailwind" {
		frontendFolderFromGit = b.TempDirPath + "/go-react-ssr/examples/frontend-tailwind"
	} else if b.StylingPlugin == "Material UI" {
		frontendFolderFromGit = b.TempDirPath + "/go-react-ssr/examples/frontend-mui"
	}
	err := cp.Copy(frontendFolderFromGit, b.ProjectDir+"/frontend")
	if err != nil {
		utils.HandleError(err)
	}
}

func (b *Bootstrapper) installNPMDependencies() {
	logger.L.Info().Msg("Installing npm dependencies")
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
	newContents = strings.Replace(newContents, "-tailwind/", "/", -1)
	newContents = strings.Replace(newContents, "-mui/", "/", -1)
	err = os.WriteFile(b.ProjectDir+"/main.go", []byte(newContents), 0644)
	if err != nil {
		utils.HandleError(err)
	}
}

func (b *Bootstrapper) updateDockerFile() {
	logger.L.Info().Msg("Updating Dockerfile")
	read, err := os.ReadFile(b.ProjectDir + "/Dockerfile")
	if err != nil {
		utils.HandleError(err)
	}
	var contents string
	contents = strings.Replace(string(read), "frontend-tailwind", "frontend", -1)
	contents = strings.Replace(string(read), "frontend-mui", "frontend", -1)
	err = os.WriteFile(b.ProjectDir+"/Dockerfile", []byte(contents), 0644)
	if err != nil {
		utils.HandleError(err)
	}
}
