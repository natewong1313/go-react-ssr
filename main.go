package main

import (
	"gossr/config"
	typeconverter "gossr/pkg/type_converter"
	"gossr/router"

	"github.com/gin-gonic/gin"
)

func main() {
    typeconverter.Scan()
    err := config.LoadConfig()
    if err != nil{
        panic(err)
    }

    g := gin.Default()
    g.Use(gin.Recovery())

    // g.Static("/public", "./public")
    g.LoadHTMLGlob(config.Config.Web.PublicDirectory+"/*")
    router.InitRoutes(g)
    
    g.Run()
}