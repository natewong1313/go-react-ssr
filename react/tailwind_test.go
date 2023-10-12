package react

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/natewong1313/go-react-ssr/config"
)

func Test_compileTailwindCssFile(t *testing.T) {
	type args struct {
		inputFilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test compiling tailwind with global css file", args{
			inputFilePath: "../examples/frontend-tailwind/src/Main.css"},
			"! tailwindcss v3.3.3 | MIT License | https://tailwindcss.com",
			false,
		},
	}

	// Create a temporary folder in the local cache directory to store the temporary CSS file
	folderPath, err := createTempCSSFolder()
	if err != nil {
		t.Errorf("createTempCSSFolder() error = %v", err)
		return
	}
	tempFilePath := filepath.ToSlash(filepath.Join(folderPath, "test.css"))
	defer os.Remove(tempFilePath)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Load(config.Config{
				GlobalCSSFilePath:  tt.args.inputFilePath,
				TailwindConfigPath: "../examples/frontend-tailwind/tailwind.config.js",
			})
			got, err := compileTailwindCssFile(tempFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("compileTailwindCssFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.Contains(got, tt.want) {
				t.Errorf("compileTailwindCssFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
