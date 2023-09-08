package main

import (
	"gossr/api/router"
	"gossr/config"
	"gossr/pkg/hot_reload"
	"gossr/pkg/type_converter"

	"github.com/gin-gonic/gin"
)

func main() {
    err := config.LoadConfig()
    if err != nil{
        panic(err)
    }

    err = type_converter.Init()
    if err != nil{
        panic(err)
    }

    hot_reload.StartWatching()

    g := gin.Default()
    g.Use(gin.Recovery())

    // g.Static("/public", "./public")
    g.LoadHTMLGlob(config.Config.Web.PublicDirectory+"/*")
    router.InitRoutes(g)
    
    g.Run()
}