package router

import (
	"gossr/models"
	reactrenderer "gossr/pkg/react_renderer"

	"github.com/gin-gonic/gin"
)


func InitRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
        reactrenderer.Render(c, "Home.tsx", &models.IndexRouteProps{InitialCount: 16, Message: "Hello from Go SSR!"})
    })
}