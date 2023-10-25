package go_ssr

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestEngine_BuildLayoutCSSFile(t *testing.T) {
	type Test struct {
		name          string
		shouldContain string
		config        *Config
	}
	originalContents, err := os.ReadFile("./examples/frontend/src/Home.css")
	assert.Nil(t, err, "ReadFile should not return an error")
	tests := []Test{
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
			shouldContain: "html {",
			config: &Config{
				AppEnv:             "production",
				FrontendDir:        "./examples/frontend-tailwind/src",
				LayoutCSSFilePath:  "Main.css",
				TailwindConfigPath: "./examples/frontend-tailwind/tailwind.config.js",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.config.setFilePaths()
			engine := &Engine{
				Logger: zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger(),
				Config: test.config,
			}
			err = engine.BuildLayoutCSSFile()
			assert.Nil(t, err, "BuildLayoutCSSFile should not return an error")
			assert.NotNilf(t, engine.CachedLayoutCSSFilePath, "CachedLayoutCSSFilePath should not be nil")
			contents, err := os.ReadFile(engine.CachedLayoutCSSFilePath)
			assert.Nil(t, err, "ReadFile should not return an error")
			assert.Contains(t, string(contents), test.shouldContain)
		})
	}
}
