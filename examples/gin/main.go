package main

import (
	"example.com/gin/models"

	"github.com/gin-gonic/gin"
	go_ssr "github.com/natewong1313/go-react-ssr"
	"github.com/natewong1313/go-react-ssr/config"
	"github.com/natewong1313/go-react-ssr/react_renderer"
)

func main() {
	g := gin.Default()
	g.StaticFile("favicon.ico", "./frontend/public/favicon.ico")
	go_ssr.Init(config.Config{
		FrontendDir:        "./frontend/src",
		GeneratedTypesPath: "./frontend/src/generated.d.ts",
		PropsStructsPath:   "./models/props.go",
		GlobalCSSFilePath:  "./frontend/src/Main.css",
		TailwindConfigPath: "./frontend/tailwind.config.js",
	})

	g.GET("/", func(c *gin.Context) {
		c.Writer.Write(react_renderer.RenderRoute(react_renderer.Config{
			File:  "Home.tsx",
			Title: "Gin example app",
			MetaTags: map[string]string{
				"og:title":    "Gin example app",
				"description": "Hello world!",
			},
			Props: &models.IndexRouteProps{
				InitialCount: 0,
			},
		}))
	})
	g.Run()
}
