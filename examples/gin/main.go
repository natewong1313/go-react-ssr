package main

import (
	"math/rand"

	"example.com/gin/models"

	"github.com/gin-gonic/gin"
	go_ssr "github.com/natewong1313/go-react-ssr"
	"github.com/natewong1313/go-react-ssr/config"
	"github.com/natewong1313/go-react-ssr/react_renderer"
)

func main() {
	g := gin.Default()
	g.StaticFile("favicon.ico", "../frontend-tailwind/public/favicon.ico")
	go_ssr.Init(config.Config{
		FrontendDir:        "../frontend-tailwind/src",
		GeneratedTypesPath: "../frontend-tailwind/src/generated.d.ts",
		TailwindConfigPath: "../frontend-tailwind/tailwind.config.js",
		GlobalCSSFilePath:  "../frontend-tailwind/src/Main.css",
		PropsStructsPath:   "./models/props.go",
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
				InitialCount: rand.Intn(100),
			},
		}))
	})
	g.Run("127.0.0.1:8080")
}
