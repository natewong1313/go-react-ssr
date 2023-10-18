package main

import (
	"example.com/echo/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	go_ssr "github.com/natewong1313/go-react-ssr"
	"github.com/natewong1313/go-react-ssr/react"
	"log"
	"math/rand"
	"net/http"
)

var APP_ENV string

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Static("/assets", "../frontend/public/")

	engine, err := go_ssr.New(&go_ssr.Config{
		AppEnv:             APP_ENV,
		AssetRoute:         "/assets",
		FrontendDir:        "../frontend/src",
		GeneratedTypesPath: "../frontend/src/generated.d.ts",
		PropsStructsPath:   "./models/props.go",
		LayoutFile:         "Layout.tsx",
	})
	if err != nil {
		log.Fatal("Failed to init go-react_old-ssr")
	}

	e.GET("/", func(c echo.Context) error {
		response := engine.RenderRoute(react.RenderConfig{
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
