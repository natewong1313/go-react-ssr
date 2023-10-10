package main

import (
	"math/rand"
	"net/http"

	"example.com/echo/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	go_ssr "github.com/natewong1313/go-react-ssr"
	"github.com/natewong1313/go-react-ssr/config"
	"github.com/natewong1313/go-react-ssr/gossr-cli/logger"
	"github.com/natewong1313/go-react-ssr/react_renderer"
)

var APP_ENV string

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Static("/assets", "../frontend/public/")

	err := go_ssr.Init(config.Config{
		AppEnv:             APP_ENV,
		AssetRoute:         "/assets",
		FrontendDir:        "../frontend/src",
		GeneratedTypesPath: "../frontend/src/generated.d.ts",
		PropsStructsPath:   "./models/props.go",
		LayoutFile:         "Layout.tsx",
	})
	if err != nil {
		logger.L.Err(err).Msg("Failed to init go-react-ssr")
		return
	}

	e.GET("/", func(c echo.Context) error {
		response := react_renderer.RenderRoute(react_renderer.Config{
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
