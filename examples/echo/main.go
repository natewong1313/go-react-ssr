package main

import (
	"log"
	"math/rand"
	"net/http"

	"example.com/echo/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	gossr "github.com/natewong1313/go-react-ssr"
)

var APP_ENV string

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Static("/assets", "../frontend-mui/public/")

	engine, err := gossr.New(gossr.Config{
		AppEnv:           APP_ENV,
		AssetRoute:       "/assets",
		FrontendSrcDir:   "../frontend/src",
		PropsStructsPath: "./models/props.go",
		LayoutFile:       "Layout.tsx",
	})
	if err != nil {
		log.Fatal("Failed to init go-react-ssr")
	}

	e.GET("/", func(c echo.Context) error {
		response := engine.RenderRoute(gossr.RenderConfig{
			File:  "Home.tsx",
			Title: "Echo example app",
			MetaTags: map[string]string{
				"og:title":    "Echo example app",
				"description": "Hello world!",
			},
			Props: &models.IndexRouteProps{
				InitialCount: rand.Intn(100),
			},
		})
		return c.HTML(http.StatusOK, string(response))
	})

	e.Start(":8080")
}
