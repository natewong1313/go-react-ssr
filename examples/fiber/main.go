package main

import (
	"log"
	"math/rand"

	"example.com/fiber/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	gossr "github.com/natewong1313/go-react-ssr"
)

var APP_ENV string

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(favicon.New(favicon.Config{
		File: "../frontend/public/favicon.ico",
		URL:  "/favicon.ico",
	}))
	app.Static("/assets", "../frontend/public/")

	engine, err := gossr.New(gossr.Config{
		AppEnv:           APP_ENV,
		AssetRoute:       "/assets",
		FrontendSrcDir:   "../frontend/src",
		PropsStructsPath: "./models/props.go",
	})
	if err != nil {
		log.Fatal("Failed to init go-react-ssr")
	}

	app.Get("/", func(c *fiber.Ctx) error {
		response := engine.RenderRoute(gossr.RenderConfig{
			File:  "Home.tsx",
			Title: "Fiber example app",
			MetaTags: map[string]string{
				"og:title":    "Fiber example app",
				"description": "Hello world!",
			},
			Props: &models.IndexRouteProps{
				InitialCount: rand.Intn(100),
			},
		})
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(string(response))
	})

	app.Listen(":8080")
}
