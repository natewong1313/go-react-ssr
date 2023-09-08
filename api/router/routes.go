package router

import (
	"gossr/api/models"
	"gossr/pkg/hot_reload"
	"gossr/pkg/react_renderer"

	"github.com/gin-gonic/gin"
)


func InitRoutes(router *gin.Engine) {
	router.GET("/ws/hotreload", hot_reload.Serve)
	router.GET("/", func(c *gin.Context) {
		react_renderer.RenderRoute(c, react_renderer.Config{
			File: "Home.tsx",
			MetaTags: map[string]string{
				"title": "My app",
				"og:title": "My app",
				"description": "Hello world!",
			},
			Props: &models.IndexRouteProps{
				InitialCount: 0,
				Message: "Hello from Go SSR!",
			},
		})
    })
}