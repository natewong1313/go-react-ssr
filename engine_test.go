package go_ssr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"testing"
	"time"
)

func TestNew_Dev(t *testing.T) {
	config := Config{
		AppEnv:              "development",
		FrontendDir:         "./examples/frontend-tailwind/src",
		GeneratedTypesPath:  "./examples/frontend-tailwind/src/generated.d.ts",
		PropsStructsPath:    "./examples/gin/models/props.go",
		HotReloadServerPort: 4001,
	}

	originalContents, err := os.ReadFile(config.GeneratedTypesPath)
	assert.Nil(t, err, "ReadFile should not return an error")
	err = os.Truncate(config.GeneratedTypesPath, 0)
	assert.Nil(t, err, "os.Truncate should not return an error, got %v", err)
	_, err = New(config)
	assert.Nil(t, err, "gossr.New should not return an error, got %v", err)

	contents, err := os.ReadFile(config.GeneratedTypesPath)
	assert.Nil(t, err, "ReadFile should not return an error")
	assert.Contains(t, string(contents), "Do not change, this code is generated from Golang structs", "Types file should have generated code in it")

	var conn net.Conn
	for i := 1; i <= 3; i++ {
		conn, _ = net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprintf("%d", config.HotReloadServerPort)), time.Second)
		if conn != nil {
			conn.Close()
			break
		}
	}
	assert.NotNil(t, conn, "Hot reload server should be running on port %d", config.HotReloadServerPort)

	err = os.WriteFile(config.GeneratedTypesPath, originalContents, 0644)
}

func TestNew_Prod(t *testing.T) {
	config := Config{
		AppEnv:             "production",
		FrontendDir:        "./examples/frontend/src",
		GeneratedTypesPath: "./examples/frontend/src/generated.d.ts",
		PropsStructsPath:   "./examples/gin/models/props.go",
	}

	originalContents, err := os.ReadFile(config.GeneratedTypesPath)
	assert.Nil(t, err, "ReadFile should not return an error")
	err = os.Truncate(config.GeneratedTypesPath, 0)
	assert.Nil(t, err, "os.Truncate should not return an error, got %v", err)
	_, err = New(config)
	assert.Nil(t, err, "gossr.New should not return an error, got %v", err)

	contents, err := os.ReadFile(config.GeneratedTypesPath)
	assert.Nil(t, err, "ReadFile should not return an error")
	assert.Equal(t, string(contents), "", "Generated types file should be empty")

	err = os.WriteFile(config.GeneratedTypesPath, originalContents, 0644)
}
