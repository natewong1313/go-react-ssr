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
    g.SetTrustedProxies(nil)
    // g.Static("/public", "./public")
    // g.LoadHTMLGlob(config.Config.Web.PublicDirectory+"/*")
    g.LoadHTMLFiles(config.Config.Web.PublicDirectory+"/index.html")
    router.InitRoutes(g)
    
    g.Run()
}