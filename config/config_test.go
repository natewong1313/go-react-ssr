package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test no config", args{}, false},
		{"test with FrontendDir that doesnt exist", args{Config{
			FrontendDir: "test",
		}}, true},
		{"test with PropsStructsPath that doesnt exist", args{Config{
			PropsStructsPath: "test",
		}}, true},
		{"test with PropsStructsPath that doesnt exist", args{Config{
			PropsStructsPath: "test",
		}}, true},
		{"test with tailwind config and global css file", args{Config{
			FrontendDir:        "../examples/frontend-tailwind/src",
			TailwindConfigPath: "../examples/frontend-tailwind/tailwind.config.js",
			GlobalCSSFilePath:  "../examples/frontend-tailwind/src/Main.css",
		}}, false},
		{"test with tailwind config but no global css file", args{Config{
			TailwindConfigPath: "../examples/frontend-tailwind/tailwind.config.js",
		}}, true},
		{"test with tailwind config but global css file that doesnt exist", args{Config{
			FrontendDir:        "../examples/frontend-tailwind/src",
			TailwindConfigPath: "../examples/frontend-tailwind/tailwind.config.js",
			GlobalCSSFilePath:  "test",
		}}, true},
		{"test with tailwind config that doesnt exist", args{Config{
			FrontendDir:        "../examples/frontend-tailwind/src",
			TailwindConfigPath: "test",
			GlobalCSSFilePath:  "../examples/frontend-tailwind/src/Main.css",
		}}, true},
		{"test with layout file that doesnt exist", args{Config{
			FrontendDir: "../examples/frontend-tailwind/src",
			LayoutFile:  "test",
		}}, true},
		{"test with layout file that does exist", args{Config{
			FrontendDir: "../examples/frontend-tailwind/src",
			LayoutFile:  "Home.tsx",
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckPathExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args string
		want bool
	}{
		{"test existing path", "./config.go", true},
		{"test non existing path", "./config.go1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkPathExists(tt.args); got != tt.want {
				t.Errorf("checkPathExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
