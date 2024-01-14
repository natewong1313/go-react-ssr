package go_ssr

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestEngine_BuildLayoutCSSFile(t *testing.T) {
	type test struct {
		name          string
		shouldContain string
		config        *Config
	}
	originalContents, err := os.ReadFile("./examples/frontend/src/Home.css")
	assert.Nil(t, err, "ReadFile should not return an error")
	tests := []test{
		{
			name:          "should clone layout css file",
			shouldContain: string(originalContents),
			config: &Config{
				AppEnv:            "production",
				FrontendDir:       "./examples/frontend/src",
				LayoutCSSFilePath: "Home.css",
			},
		},
		{
			name:          "should build layout css file with tailwind",
			shouldContain: "tailwindcss",
			config: &Config{
				AppEnv:             "production",
				FrontendDir:        "./examples/frontend-tailwind/src",
				LayoutCSSFilePath:  "Main.css",
				TailwindConfigPath: "./examples/frontend-tailwind/tailwind.config.js",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.setFilePaths()
			engine := &Engine{
				Logger: zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger(),
				Config: tt.config,
			}
			err = engine.BuildLayoutCSSFile()
			assert.Nil(t, err, "BuildLayoutCSSFile should not return an error, got %v", err)
			assert.NotNilf(t, engine.CachedLayoutCSSFilePath, "CachedLayoutCSSFilePath should not be nil")
			contents, err := os.ReadFile(engine.CachedLayoutCSSFilePath)
			assert.Nil(t, err, "ReadFile should not return an error, got %v", err)
			assert.Contains(t, string(contents), tt.shouldContain)
		})
	}
}
