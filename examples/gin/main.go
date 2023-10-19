package main

import (
	"log"
	"math/rand"

	"example.com/gin/models"

	"github.com/gin-gonic/gin"
	gossr "github.com/natewong1313/go-react-ssr"
)

var APP_ENV string

func main() {
	g := gin.Default()
	g.StaticFile("favicon.ico", "../frontend-tailwind/public/favicon.ico")
	g.Static("/assets", "../frontend-tailwind/public")
	engine, err := gossr.New(gossr.Config{
		AppEnv:             APP_ENV,
		AssetRoute:         "/assets",
		FrontendDir:        "../frontend-tailwind/src",
		GeneratedTypesPath: "../frontend-tailwind/src/generated.d.ts",
		TailwindConfigPath: "../frontend-tailwind/tailwind.config.js",
		LayoutCSSFilePath:  "Main.css",
		PropsStructsPath:   "./models/props.go",
	})
	if err != nil {
		log.Fatal("Failed to init go-react-ssr")
	}

	g.GET("/", func(c *gin.Context) {
		c.Writer.Write(engine.RenderRoute(gossr.RenderConfig{
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
	g.Run()
}
